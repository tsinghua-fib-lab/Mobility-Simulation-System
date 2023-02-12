package vehicle

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
)

type Runtime struct {
	entity.BaseRuntime
	entity.BaseRuntimeOnRoad

	// 车辆额外信息

	VehicleStatus Status
	Action        Action // 车辆行为
	DistanceToEnd float64

	// 以下成员在变道时使用，仅当shadow_lane不为空时有意义

	ShadowLane                entity.ILane // 影子所在车道
	ShadowS                   float64      // 影子在车道上的位置
	LaneChangeTotalLength     float64      // 变道总长
	LaneChangeCompletedLength float64      // 变道已完成的长度
	IsLaneChanging            bool         // 变道状态
}

func (rt *Runtime) computeDirection() {
	const TAU = 6.2831853
	if rt.ShadowLane != nil {
		direction := rt.ShadowLane.GetDirectionByS(rt.ShadowS)
		laneChangeBias := math.Atan2((rt.Lane.Width()+rt.ShadowLane.Width())/2,
			rt.LaneChangeTotalLength) *
			(1 - math.Abs(rt.LaneChangeCompletedLength/rt.LaneChangeTotalLength*2-1))
		if rt.ShadowLane == rt.Lane.FirstLeftLane() {
			direction += laneChangeBias
			if direction >= TAU {
				rt.Direction = direction - TAU
			} else {
				rt.Direction = direction
			}
		} else if rt.ShadowLane == rt.Lane.FirstRightLane() {
			direction -= laneChangeBias
			if direction >= TAU {
				rt.Direction = direction - TAU
			} else {
				rt.Direction = direction
			}
		} else {
			log.Panicf("vehicle: shadow lane %d is neither the left one of lane %d nor the right one", rt.ShadowLane.ID(), rt.Lane.ID())
		}
	} else {
		rt.Direction = rt.Lane.GetDirectionByS(rt.S)
	}
}

func (rt *Runtime) clearLaneChange() {
	rt.ShadowLane = nil
	rt.ShadowS = 0
	rt.LaneChangeTotalLength = 0
	rt.LaneChangeCompletedLength = 0
	rt.IsLaneChanging = false
}
