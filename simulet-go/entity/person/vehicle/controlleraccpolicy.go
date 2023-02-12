package vehicle

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	traffic_lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"git.fiblab.net/sim/simulet-go/utils"
	"github.com/samber/lo"
)

// 策略模板1：前方有其他车辆，跟车
func (l *controller) carFollow(ahead entity.IVehicle, distance float64) (ac Action) {
	ac.Acc = l.computeCarFollowAcc(
		ahead.Speed(),
		ahead.Attr().MaxBrakingAcceleration,
		distance,
	)
	return
}

// 策略模板2：前一车道
func (l *controller) laneAhead(ahead entity.ILane, distance float64) (ac Action) {
	ac.Acc = utils.INF
	if l.speed > ahead.MaxSpeed() {
		// 超速需要减速
		acc := computeAcc(
			ahead.MaxSpeed(),
			l.speed,
			distance-l.speed*l.stepInterval,
		)
		if acc < l.usualBrakingAcc {
			ac.Acc = acc
		}
	}
	// 一般刹车不能让车辆停在路口了，得开始减速
	// vt+at^2/2 (a=一般刹车) > distance to lane end
	cannotStop := l.speed*l.stepInterval+computeBrakingDistance(l.speed, l.usualBrakingAcc) >= distance
	if ahead.InJunction() && cannotStop {
		// 需要开始判断路口信控情况
		switch ahead.LightState() {
		case traffic_lightv2.LightState_LIGHT_STATE_RED:
			// 红灯减速停车
			ac.Acc = lo.Min([]float64{ac.Acc, -0.1, computeAcc(0, l.speed, distance)})
		case traffic_lightv2.LightState_LIGHT_STATE_YELLOW:
			// 黄灯，倒计时结束前不可过线，减速停车
			if ahead.LightStateRemainingTime()*l.speed <= distance {
				ac.Acc = lo.Min([]float64{ac.Acc, -0.1, computeAcc(0, l.speed, distance)})
			}
		default:
			// 绿灯或没灯，跳过
		}
	}
	return
}

// 策略1：开到限速
func (l *controller) policyToLimit() (ac Action) {
	maxSpeed := math.Min(l.v.attr.MaxSpeed, l.curLane.MaxSpeed())
	ac.Acc = l.computeUsualAccToTargetSpeed(maxSpeed)
	return
}

// 策略2：前车
func (l *controller) policyCarFollow() (ac Action) {
	ac.Acc = utils.INF
	// 感知前车
	aheadV := l.node.Next().ValueOrDefault(nil)
	aheadS := 0.0
	if aheadV != nil {
		aheadS = aheadV.S()
	} else if l.route.Current().IsForward() {
		// 检查前方车道
		forwardLane := l.route.Next().Lane
		aheadV = forwardLane.GetFirstVehicle()
		if aheadV != nil {
			aheadS = aheadV.S() + l.curLane.Length()
		}
	}
	// 精确重叠的，当作不存在
	if aheadV != nil && math.Abs(aheadS-l.s) > 1e-6 {
		return l.carFollow(aheadV, aheadS-l.s-aheadV.Length())
	} else {
		return
	}
}

// 策略2：影子的前车
func (l *controller) policyShadowCarFollow() (ac Action) {
	ac.Acc = utils.INF
	// 感知前车
	aheadV := l.v.shadowNode.Next().ValueOrDefault(nil)
	aheadS := 0.0
	if aheadV != nil {
		aheadS = aheadV.S()
	} else if l.route.ShadowCurrent().IsForward() {
		// 检查前方车道
		forwardLane := l.route.ShadowNext().Lane
		aheadV = forwardLane.GetFirstVehicle()
		if aheadV != nil {
			aheadS = aheadV.S() + l.v.ShadowLane().Length()
		}
	}
	shadowS := l.v.ShadowS()
	// 精确重叠的，当作不存在
	if aheadV != nil && math.Abs(aheadS-shadowS) > 1e-6 {
		return l.carFollow(aheadV, aheadS-l.v.ShadowS()-aheadV.Length())
	} else {
		return
	}
}

// 策略3：下一车道
func (l *controller) policyLaneAhead() (ac Action) {
	ac.Acc = utils.INF
	if l.route.Current().IsForward() {
		return l.laneAhead(l.route.Next().Lane, l.reverseS)
	}
	return
}

// 策略3：下一车道
func (l *controller) policyShadowLaneAhead() (ac Action) {
	ac.Acc = utils.INF
	if l.route.ShadowCurrent().IsForward() {
		return l.laneAhead(l.route.ShadowNext().Lane, l.reverseS)
	}
	return
}

// 策略4：终点
func (l *controller) policyToEnd() (ac Action) {
	distanceToEnd := l.v.runtime.DistanceToEnd
	ac.Acc = utils.INF
	if l.speed*l.stepInterval+computeBrakingDistance(l.speed, l.usualBrakingAcc) >= distanceToEnd {
		ac.Acc = computeAcc(0, l.speed, distanceToEnd)
	}
	return
}
