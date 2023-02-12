package vehicle

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"
)

const (
	LC_FORBIDDEN_DISTANCE    = 5        // 车道前后不能变道的长度
	LC_CHECK_LANE_MIN_LENGTH = 20       // 检查(前后不能变道的)车道长度的最小值
	LC_MIN_P                 = 0.2      // 随机发起变道概率下限
	LC_LENGTH_FACTOR         = 3        // 变道长度与车速的关系
	LC_MIN_SPEED             = 10 / 3.6 // 减速等待变道能下到的最小速度
	LC_MIN_INTERVAL          = 20       // 自主变道后期望最少在目标车道停留的长度
	LC_MIN_CHANGE_SPEED_GAIN = 5 / 3.6  // 触发变道决策的最小速度增益（m/s）
)

type laneChangeEnv struct {
	count    int     // 变道数量
	distance float64 // 预期车辆变道距离

	sideLane      entity.ILane // 变道目标（通过route获得的）
	sideS         float64
	sideLaneAhead entity.ILane    // 变道目标的下一车道
	beforeV       entity.IVehicle // 变道目标对应位置的后车
	afterV        entity.IVehicle // 变道目标对应位置的前车
}

func (l *controller) getLaneChangeEnv() (e laneChangeEnv) {
	curSegment := l.route.Current()
	if curSegment.IsLaneChange() {
		// 数变道次数
		e.count = l.route.CountContinuousLaneChanges()
	}
	e.distance = computeBrakingDistance(l.speed, l.usualBrakingAcc) +
		l.speed*l.stepInterval +
		float64(e.count)*l.laneChangeLength

	if curSegment.IsLaneChange() {
		shadowCurSegment := l.route.ShadowCurrent()
		e.sideLane = shadowCurSegment.Lane
		e.sideS = e.sideLane.ProjectFromLane(l.curLane, l.s)
		if shadowCurSegment.IsForward() {
			e.sideLaneAhead = l.route.ShadowNext().Lane
		}
		// 找出变道位置前后车辆
		var side int
		if curSegment.NextLaneType == routingv2.NextLaneType_NEXT_LANE_TYPE_LEFT {
			side = entity.LEFT
		} else {
			side = entity.RIGHT
		}
		links := l.node.Extra.Links[side]
		if beforeV := links[entity.BEFORE].ValueOrDefault(nil); beforeV != nil {
			e.beforeV = beforeV
		}
		if afterV := links[entity.AFTER].ValueOrDefault(nil); afterV != nil {
			e.afterV = afterV
		}
	}
	return
}

func (l *controller) setLaneChange(lc *laneChangeEnv, ac *Action, force bool) {
	var vehicleAheadDistance float64
	if force {
		goto LANE_CHANGE
	}
	if lc.afterV != nil {
		distance := lc.afterV.S() - lc.sideS - lc.afterV.Length()
		if distance <= 0 || computeSpeed(
			lc.afterV.Speed(),
			l.usualBrakingAcc,
			distance,
		) < l.speed {
			// 判断现在直接变道会撞前车，暂时无法变道，减速(直到LC_MIN_SPEED)等待路况合适
			if acc := l.computeUsualAccToTargetSpeed(LC_MIN_SPEED); acc < 0 {
				ac.UpdateByMinAcc(Action{
					Acc: acc,
				})
			}
			return
		}
		vehicleAheadDistance = distance
	}
	if lc.beforeV != nil {
		distance := lc.sideS - lc.beforeV.S()
		if distance <= 0 || computeSpeed(
			l.speed,
			lc.beforeV.Attr().UsualBrakingAcceleration,
			distance,
		) < lc.beforeV.Speed() {
			// 判断现在直接变道会让后车追尾，暂时无法变道，并减速(直到LC_MIN_SPEED)等待路况合适
			if acc := l.computeUsualAccToTargetSpeed(LC_MIN_SPEED); acc < 0 {
				ac.UpdateByMinAcc(Action{
					Acc: acc,
				})
			}
		}
		return
	}
LANE_CHANGE:
	ac.StartLaneChange(l.speed * LC_LENGTH_FACTOR)
	// 执行跟车策略
	if lc.afterV != nil {
		ac.UpdateByMinAcc(l.carFollow(lc.afterV, vehicleAheadDistance))
	}
}

func (l *controller) planNecessaryLaneChange(e *laneChangeEnv, ac *Action) {
	if e.count > 0 {
		// 距离不足，强制变道
		if l.reverseS <= e.distance {
			l.setLaneChange(e, ac, true)
			return
		}
		// 车道前后5m不变道（仅限车道长度超过20m的）
		if l.curLane.Length() > LC_CHECK_LANE_MIN_LENGTH &&
			(l.s < LC_FORBIDDEN_DISTANCE || l.reverseS < LC_FORBIDDEN_DISTANCE) {
			return
		}
		// 根据到车道尽头的距离调整触发变道的概率
		if l.v.driver.Generator().PTrue(math.Max(
			l.reverseS/l.curLane.Length(),
			LC_MIN_P,
		)) {
			l.setLaneChange(e, ac, false)
			return
		}
	}
}
