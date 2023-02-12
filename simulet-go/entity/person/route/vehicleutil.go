package route

import routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"

func (r *VehicleRoute) canShadow() {
	if !r.route[0].IsLaneChange() {
		log.Panic("shadow lane is out of route")
	}
}

func (r *VehicleRoute) initDistanceToEnd() {
	// 计算到终点的距离
	r.distanceToEnd = make([]float64, r.size)
	endPos := r.end
	if lanePos := endPos.LanePosition; lanePos != nil {
		r.distanceToEnd[len(r.distanceToEnd)-1] = lanePos.S
	} else if aoiPos := endPos.AoiPosition; aoiPos != nil {
		endLane := r.route[len(r.route)-1].Lane
		if aoi := endLane.Aois()[aoiPos.AoiId]; aoi != nil {
			r.distanceToEnd[len(r.distanceToEnd)-1] = aoi.DrivingS(endLane.ID())
		} else {
			log.Panicf("VehicleRoute: no aoi %v on lane %v", aoiPos.AoiId, endLane.ID())
		}
	}
	for i := r.size - 2; i >= 0; i-- {
		segment := r.route[i]
		switch segment.NextLaneType {
		case routingv2.NextLaneType_NEXT_LANE_TYPE_FORWARD:
			// 直行，累加本车道长度
			r.distanceToEnd[i] = r.distanceToEnd[i+1] + segment.Lane.Length()
		case routingv2.NextLaneType_NEXT_LANE_TYPE_LEFT,
			routingv2.NextLaneType_NEXT_LANE_TYPE_RIGHT:
			r.distanceToEnd[i] = r.distanceToEnd[i+1]
		default:
			log.Panic("VehicleRoute: wrong NextLaneType")
		}
	}
}
