package vehicle

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/person/route"
	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	"git.fiblab.net/sim/simulet-go/utils/container"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
)

// 计算本时刻的速度与移动距离
// v(t)=v(t-1)+acc*dt, ds=v(t-1)*dt+acc*dt*dt/2
func computeSpeedAndDistance(speed, acc, dt float64) (float64, float64) {
	dv := acc * dt
	if speed+dv < 0 {
		// 刹车到停止
		return 0, speed * speed / 2 / -acc
	}
	return speed + dv, (speed + dv/2) * dt
}

// 车辆状态
type Status int32

const (
	Status_RUNNING           Status = iota // 行进
	Status_PAUSE                           // 停在路上
	Status_WAITING_FOR_ROUTE               // 到达出发时间等待导航请求
)

type Vehicle struct {
	attr        *agentv2.AgentAttribute
	vehicleAttr *agentv2.VehicleAttribute
	length      float64 // 车辆长度属性，避免直接从base中读取带来的多重指针访问

	runtime  Runtime // 运行时数据
	snapshot Runtime // 上一时刻快照

	driver entity.IPerson // 司机

	route *route.VehicleRoute // 路经规划

	skipToEnd bool // 异常情况，车辆需要瞬移到终点
	isEnd     bool // 车辆生命周期是否结束

	// Lane链表
	node, shadowNode *container.ListNode[entity.IVehicle, entity.VehicleSideLink]
}

func NewVehicle(driver entity.IPerson, route *route.VehicleRoute) *Vehicle {
	lane := route.Current().Lane
	// NewVehicle一定由从AOI里面出发的行程创建
	s := route.Start.Aoi.DrivingS(lane.ID())
	v := &Vehicle{
		attr:        driver.Attr(),
		vehicleAttr: driver.VehicleAttr(),
		length:      driver.Attr().Length,
		runtime: Runtime{
			BaseRuntime: entity.BaseRuntime{
				Position:  lane.GetPositionByS(s),
				Speed:     0,
				Direction: lane.GetDirectionByS(s),
			},
			BaseRuntimeOnRoad: entity.BaseRuntimeOnRoad{
				Lane: lane,
				S:    s,
			},
			VehicleStatus: Status_RUNNING,
			DistanceToEnd: route.GetDistanceToEnd(s),
		},
		driver: driver,
		route:  route,
	}
	v.node = newNode(v.runtime.S, v)
	v.shadowNode = newNode(v.runtime.S, v)
	lane.ReportVehicleAdded(v.node)
	return v
}

func (v *Vehicle) Prepare() {
	v.snapshot = v.runtime
}

func (v *Vehicle) Update(stepInterval float64) {
	switch v.runtime.VehicleStatus {
	case Status_RUNNING:
		v.runtime.Action = newController(v).Update(stepInterval)
		v.refreshRuntimeByAction(v.runtime.Action, stepInterval)
		reachTarget := v.checkCloseToEndAndRefreshRuntime()

		// 完成计算，清空支链
		v.node.Extra.Clear()
		if v.snapshot.IsLaneChanging {
			v.shadowNode.Extra.Clear()
		}

		if reachTarget {
			// 增量更新车道索引（不再维护数据）
			v.snapshot.Lane.ReportVehicleRemoved(v.node)
			if v.snapshot.IsLaneChanging {
				v.snapshot.ShadowLane.ReportVehicleRemoved(v.shadowNode)
			}
			endAoi := v.route.GetEndAoi()
			if endAoi != nil {
				// 目标点为Aoi，私家车的车辆生命周期结束
				v.isEnd = true
			} else {
				// 目标为lane上位置，进入PAUSE状态，等待司机的下一trip开始
				v.runtime.VehicleStatus = Status_PAUSE
				v.driver.NextTrip()
			}
			return
		}

		// 车道链表更新
		v.node.Key = v.runtime.S
		if v.runtime.IsLaneChanging {
			v.shadowNode.Key = v.runtime.ShadowS
		}
		// 增量更新车道索引（维护数据）
		if v.snapshot.Lane != v.runtime.Lane {
			v.snapshot.Lane.ReportVehicleRemoved(v.node)
			// 换一个新的node来避免remove操作和add操作处理同一个对象需要保证先后顺序
			v.node = newNode(v.runtime.S, v)
			v.runtime.Lane.ReportVehicleAdded(v.node)
		}
		if !v.snapshot.IsLaneChanging && !v.runtime.IsLaneChanging {
			// do nothing
		} else if v.snapshot.IsLaneChanging && !v.runtime.IsLaneChanging {
			v.snapshot.ShadowLane.ReportVehicleRemoved(v.shadowNode)
		} else if !v.snapshot.IsLaneChanging && v.runtime.IsLaneChanging {
			v.runtime.ShadowLane.ReportVehicleAdded(v.shadowNode)
		} else {
			if v.snapshot.ShadowLane != v.runtime.ShadowLane {
				v.snapshot.ShadowLane.ReportVehicleRemoved(v.shadowNode)
				v.shadowNode = newNode(v.runtime.ShadowS, v)
				v.runtime.ShadowLane.ReportVehicleAdded(v.shadowNode)
			}
		}
	case Status_PAUSE:
		if v.driver.CheckDeparture() {
			// 到达出发时间，发送导航请求
			v.runtime.VehicleStatus = Status_WAITING_FOR_ROUTE
			v.route.RerouteFlag = true
		}
		// 在路上暂停的车辆不将信息记录到lane中
	case Status_WAITING_FOR_ROUTE:
		if v.route.Ok() {
			// 导航成功，出发
			v.runtime.VehicleStatus = Status_RUNNING
			// 维护Linked list
			v.node.Key = v.runtime.S
			v.runtime.Lane.ReportVehicleAdded(v.node)
		} else {
			// 导航失败，等待下一trip
			v.driver.NextTrip()
			v.runtime.VehicleStatus = Status_PAUSE
		}
		// 在路上暂停的车辆不将信息记录到lane中
	default:
		log.Panicf("unknown vehicle %d status %v", v.ID(), v.runtime.VehicleStatus)
	}

	if v.route.RerouteFlag {
		// 当该标志被设置时，无条件地进行路径规划
		v.route.ProduceRouting(
			v.driver.Schedule().GetTrip().End,
			route.RouteStartPosition{Lane: v.runtime.Lane, S: v.runtime.S},
		)
		v.route.RerouteFlag = false
	}
}

func (v *Vehicle) refreshRuntimeByAction(ac Action, stepInterval float64) {
	speed, ds := computeSpeedAndDistance(v.Speed(), ac.Acc, stepInterval)
	v.runtime.Speed = speed

	// 更新位置

	if ac.EnableLaneChange {
		// 发起变道
		targetLane := v.route.Next().Lane
		//  --------------------------------------------
		//   [2] → → (lane_change_length / ds) → → [3]
		//  --↑-----------------------------------------
		//   [1]     (ignore the width)
		//  --------------------------------------------
		// 1: motion.lane + motion.s
		// 2: target_lane + neighbor_s
		// 3: target_lane + target_s
		neighborS := targetLane.ProjectFromLane(v.runtime.Lane, v.runtime.S)
		// 变道必须在当前道路内完成
		targetS := math.Min(neighborS+ac.LaneChangeLength, targetLane.Length())
		if neighborS+ds >= targetS {
			// 跳过变道
			v.route.Step()
			v.driveStraightAndRefreshLocation(targetLane, neighborS, ds)
		} else {
			//  --------------------------------------------
			//   [ns] → → → → [ns+ds] → → → → [ts]
			//  --------------------------------------------
			//   [1]            [s]
			//  --------------------------------------------
			// ns: neighbor_s
			// ds: ds
			// ts: target_s
			// s: motion.s
			v.runtime.IsLaneChanging = true
			v.runtime.ShadowLane = targetLane
			v.runtime.ShadowS = neighborS + ds
			v.runtime.S = v.runtime.Lane.ProjectFromLane(v.runtime.ShadowLane, neighborS+ds)
			v.runtime.LaneChangeTotalLength = targetS - neighborS
			v.runtime.LaneChangeCompletedLength = ds
		}
	} else {
		if v.runtime.IsLaneChanging {
			// 正在变道
			if v.runtime.LaneChangeCompletedLength+ds >= v.runtime.LaneChangeTotalLength {
				// 变道完成
				v.route.Step()
				// 更新位置为变道目标
				v.driveStraightAndRefreshLocation(v.runtime.ShadowLane, v.runtime.ShadowS, ds)
				v.runtime.clearLaneChange()
			} else {
				v.runtime.LaneChangeCompletedLength += ds
				v.runtime.ShadowS += ds
				v.runtime.S = v.runtime.Lane.ProjectFromLane(v.runtime.ShadowLane, v.runtime.ShadowS)
			}
		} else {
			// 直行
			v.driveStraightAndRefreshLocation(v.runtime.Lane, v.runtime.S, ds)
		}
	}
	// 更新xy坐标
	v.runtime.Position = v.runtime.Lane.GetPositionByS(v.runtime.S)
	// 更新到终点的距离
	v.runtime.DistanceToEnd = v.route.GetDistanceToEnd(v.runtime.S)
	// 更新车辆方向角
	v.runtime.computeDirection()
}

func (v *Vehicle) driveStraightAndRefreshLocation(lane entity.ILane, s, ds float64) {
	s += ds
	if s > lane.Length() {
		if v.runtime.ShadowLane != nil {
			log.Warnf("vehicle: vehicle %v skipped the change to lane %v (status=%v)",
				v.ID(), v.runtime.ShadowLane.ID(),
				v.runtime.IsLaneChanging)
		}
		v.runtime.clearLaneChange()
		for s > lane.Length() {
			isLaneChange := v.route.Current().IsLaneChange()
			if ok := v.route.Step(); !ok {
				v.skipToEnd = true
				return
			}
			teleport := false
			// 步进到第一个NextLaneType是直行的路段的下一路段（进行一次直行）
			for isLaneChange {
				isLaneChange = v.route.Current().IsLaneChange()
				if ok := v.route.Step(); !ok {
					v.skipToEnd = true
					return
				}
				teleport = true
			}
			nextLane := v.route.Current().Lane
			s -= lane.Length()
			if teleport {
				log.Warnf("vehicle: teleport vehicle %v to lane %v, which may be caused by too short lane %v (length=%v) or lane-change failure",
					v.ID(), nextLane.ID(), lane.ID(), lane.Length())
			}
			lane = nextLane
		}
	}
	v.runtime.Lane = lane
	v.runtime.S = s
}

// 检查车辆是否到达目标地点，是则返回true
func (v *Vehicle) checkCloseToEndAndRefreshRuntime() bool {
	if v.skipToEnd || v.runtime.DistanceToEnd <= CLOSE_TO_END {
		// 到达目的地，设置motion为目的地的路面位置（供人进入aoi时选择gate）
		v.runtime.Lane = v.route.Last().Lane
		v.runtime.S = v.route.GetEndS()
		v.runtime.Speed = 0
		v.runtime.clearLaneChange()
		v.runtime.computeDirection()
		if v.skipToEnd {
			log.Warnf("skipToEnd: vehicle %v from %+v to %+v",
				v.ID(), v.snapshot, v.runtime,
			)
		}
		return true
	} else {
		return false
	}
}

// base

func (v *Vehicle) ID() int32 {
	return v.driver.ID()
}

func (v *Vehicle) Attr() *agentv2.AgentAttribute {
	return v.attr
}

func (v *Vehicle) VehicleAttr() *agentv2.VehicleAttribute {
	return v.vehicleAttr
}

func (v *Vehicle) Length() float64 {
	return v.length
}

func (v *Vehicle) FetchBaseSnapshotForPerson() (entity.BaseRuntime, entity.BaseRuntimeOnRoad) {
	return v.snapshot.BaseRuntime, v.snapshot.BaseRuntimeOnRoad
}

// getter

func (v *Vehicle) Snapshot() Runtime {
	return v.snapshot
}

func (v *Vehicle) Speed() float64 {
	return v.snapshot.Speed
}

func (v *Vehicle) Lane() entity.ILane {
	return v.snapshot.Lane
}

func (v *Vehicle) S() float64 {
	return v.snapshot.S
}

func (v *Vehicle) Direction() float64 {
	return v.snapshot.Direction
}

func (v *Vehicle) Position() geometry.Point {
	return v.snapshot.Position
}

func (v *Vehicle) ShadowS() float64 {
	return v.snapshot.ShadowS
}

func (v *Vehicle) ShadowLane() entity.ILane {
	return v.snapshot.ShadowLane
}

func (v *Vehicle) IsLaneChanging() bool {
	return v.snapshot.IsLaneChanging
}

func (v *Vehicle) GetEndByPerson() (entity.IAoi, entity.ILane, bool) {
	if v.isEnd {
		return v.route.GetEndAoi(), v.runtime.Lane, true
	} else {
		return nil, nil, false
	}
}
