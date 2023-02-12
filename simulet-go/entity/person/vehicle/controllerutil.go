package vehicle

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	"github.com/samber/lo"
)

// 计算在给定距离内改变速度所需的加速度
func computeAcc(vEnd, vNow, distance float64) float64 {
	return (vEnd*vEnd - vNow*vNow) / 2 / math.Max(distance, 1e-6)
}

// 计算在distance距离以加速度a减速为v所需的初速度
func computeSpeed(v, a, distance float64) float64 {
	return math.Sqrt(v*v - 2*a*distance)
}

// 计算采用给定加速度刹停所需的距离
func computeBrakingDistance(v, a float64) float64 {
	if a >= 0 {
		log.Panic("bad braking acc")
	}
	return v * v / 2 / -a
}

// 检查在lane的position位置是否会与行人冲突
// 冲突范围是"行人5米内到达冲突点或离开冲突点不足3米"
func checkPersonInLane(checkLane entity.ILane, position float64) bool {
	// 假定有check_before>check_after
	const CHECK_BEFORE, CHECK_AFTER = 5, 3

	list := checkLane.Pedestrians()
	if length := list.Length(); length > 10 {
		log.Errorf("try to checkPersonInLane with len=%d at lane %d", length, checkLane.ID())
	}
	for node := list.First(); node != nil; node = node.Next() {
		p := node.Value
		if p.IsForward() {
			if p.S() < position+CHECK_AFTER {
				return true
			}
		} else {
			if p.S() > position-CHECK_BEFORE {
				return true
			}
		}
	}
	return false
}

// 采用一般刹车和加速clamp，计算调整到targetSpeed的加速度
func (l *controller) computeUsualAccToTargetSpeed(target float64) float64 {
	return lo.Clamp(
		(target-l.speed)/l.stepInterval,
		l.usualBrakingAcc,
		l.usualAcc,
	)
}

// Krauss模型：
// 记本车为B，前车为A，本时刻速度为v_B和v_A，那么本车下一时刻的速度u需要满足
// 下一时刻本车普通刹车的刹车距离 + 本时刻距离 <= 前车急刹的距离+现有距离
// u^2/(2a_B)+(v_B+u)t/2 ≤ (v_A^2)/(2a_A)+d
// distance: 本车车头到前车车尾的距离
func (l *controller) computeCarFollowAcc(aheadSpeed, aheadMaxBrakingAcc, distance float64) float64 {
	if distance <= 0 {
		return l.maxBrakingAcc
	}
	a, b := .5/-l.usualBrakingAcc, .5*l.stepInterval
	c := .5*l.speed*l.stepInterval -
		aheadSpeed*aheadSpeed*.5/-aheadMaxBrakingAcc +
		l.minGap - distance
	det := b*b - 4*a*c
	if det < 0 {
		// 紧急刹车
		return l.maxBrakingAcc
	} else {
		targetSpeed := math.Max(0, (math.Sqrt(det)-b)/a/2)
		return l.computeUsualAccToTargetSpeed(targetSpeed)
	}
}
