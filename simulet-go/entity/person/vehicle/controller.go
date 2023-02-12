package vehicle

import (
	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/person/route"
	"git.fiblab.net/sim/simulet-go/utils"
	"git.fiblab.net/sim/simulet-go/utils/container"

	"github.com/samber/lo"
)

const (
	// 变道意愿阈值
	LEFT_MOTIVATION_THRESHOLD  = 4
	RIGHT_MOTIVATION_THRESHOLD = 12

	DEADLOCK_SPEED                                      = .1    // 死锁判决速度
	COLLISION_AVOIDANCE_RATIO                           = 1     // 为了避免碰撞所留出的空间相对于半车身宽度的比例
	OVERLAP_CONGESTION_OCCUPANCY_THRESHOLD              = .6    // 判断交叉车道拥堵的占有率阈值
	NEXT_LANE_CONGESTION_STATIC_VEHICLE_COUNT_THRESHOLD = 2     // 判断下一条车道拥堵的静止车辆数阈值
	VIEW_DISTANCE                                       = 200.0 // 车辆最大视野
	CLOSE_TO_END                                        = 5     // 车辆到达终点的判定范围
)

type controller struct {
	v                               *Vehicle // 模块所在车辆
	route                           *route.VehicleRoute
	node                            *container.ListNode[entity.IVehicle, entity.VehicleSideLink]
	leftMotivation, rightMotivation float64 // 变道意愿
	usualBrakingAcc                 float64
	maxBrakingAcc                   float64
	usualAcc                        float64
	length                          float64
	minGap                          float64
	laneChangeLength                float64
	speed                           float64
	s                               float64
	reverseS                        float64
	curLane                         entity.ILane
	stepInterval                    float64 // TODO: 考虑直接使用clock里面的值，不再传递
}

func newController(v *Vehicle) *controller {
	// 数据预读
	attr := v.attr
	c := &controller{
		v:                v,
		route:            v.route,
		node:             v.node,
		usualBrakingAcc:  attr.UsualBrakingAcceleration,
		maxBrakingAcc:    attr.MaxBrakingAcceleration,
		usualAcc:         attr.UsualAcceleration,
		length:           attr.Length,
		minGap:           v.vehicleAttr.MinGap,
		laneChangeLength: v.vehicleAttr.LaneChangeLength,
		speed:            v.Speed(),
		s:                v.S(),
		curLane:          v.Lane(),
	}
	c.reverseS = c.curLane.Length() - c.s
	return c
}

func (l *controller) Update(stepInterval float64) Action {
	// 更新参数
	l.stepInterval = stepInterval

	// 执行策略
	ac := Action{}
	ac.Acc = utils.INF

	// 执行变道决策
	if !l.v.IsLaneChanging() && !l.curLane.InJunction() {
		lcEnv := l.getLaneChangeEnv()
		l.planNecessaryLaneChange(&lcEnv, &ac)
	}

	// 执行加速度决策
	ac.UpdateByMinAcc(
		l.policyToLimit(),
		l.policyCarFollow(),
		l.policyLaneAhead(),
		l.policyToEnd(),
	)
	if l.v.IsLaneChanging() {
		ac.UpdateByMinAcc(
			l.policyShadowCarFollow(),
			l.policyShadowLaneAhead(),
		)
	}

	// 后处理
	ac.Acc = lo.Clamp(ac.Acc, l.v.attr.MaxBrakingAcceleration, l.v.attr.MaxAcceleration)
	if ac.EnableLaneChange {
		ac.LaneChangeLength = lo.Clamp(
			ac.LaneChangeLength,
			l.length,
			l.laneChangeLength,
		)
	}
	return ac
}
