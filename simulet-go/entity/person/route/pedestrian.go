package route

import (
	"context"

	"git.fiblab.net/sim/simulet-go/entity"
	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"
)

type PedestrianSegment struct {
	Lane      entity.ILane
	Direction routingv2.MovingDirection
}

func (s PedestrianSegment) IsForward() bool {
	return s.Direction == routingv2.MovingDirection_MOVING_DIRECTION_FORWARD
}

type PedestrianRoute struct {
	Start        RouteStartPosition          // 导航起点
	base         *routingv2.GetRouteResponse // 导航请求结果
	ok           bool                        // 导航请求是否成功
	indexJourney int                         // 当前journey的索引
	indexRoute   int                         // 当前行驶车道对应的索引编号，即route[route_index_]==parent_lane
	route        []PedestrianSegment         // 转换为指针形式后的路径
	end          *geov2.Position             // 记录路径规划终点
	EndS         float64                     // 终点s
	EndAoi       entity.IAoi                 // 终点Aoi
}

func NewPedestrianRoute() *PedestrianRoute {
	return &PedestrianRoute{
		route: make([]PedestrianSegment, 0),
	}
}

func (r *PedestrianRoute) Ok() bool {
	return r.ok
}

func (r *PedestrianRoute) AtLast() bool {
	return r.indexRoute+1 >= len(r.route)
}

func (r *PedestrianRoute) Current() PedestrianSegment {
	return r.route[r.indexRoute]
}

func (r *PedestrianRoute) Next() PedestrianSegment {
	return r.route[r.indexRoute+1]
}

func (r *PedestrianRoute) Last() PedestrianSegment {
	return r.route[len(r.route)-1]
}

// 向前增加index，返回是否正常（true: 正常, false：越界）
func (r *PedestrianRoute) Step() bool {
	r.indexRoute++
	if r.indexRoute >= len(r.route) {
		r.indexRoute = len(r.route) - 1
		return false
	}
	return true
}

func (r *PedestrianRoute) ProduceRouting(target *geov2.Position, startPosition RouteStartPosition, routeType routingv2.RouteType) {
	r.Start = startPosition
	r.ok = false
	// 仅支持从aoi出发
	req := routingv2.GetRouteRequest{
		Type: routeType,
		Start: &geov2.Position{
			AoiPosition: &geov2.AoiPosition{AoiId: startPosition.Aoi.ID()},
		},
		End: target,
	}
	// 记录路径规划终点
	r.end = target
	// 发送路径规划请求
	RoutingClient.GetRoute(context.Background(), &req, r.ProcessRouting)
}

func (r *PedestrianRoute) ProcessRouting(res *routingv2.GetRouteResponse) {
	var laneStart entity.ILane
	if len(res.Journeys) == 0 {
		r.route = make([]PedestrianSegment, 0)
		r.indexRoute = 0
		r.ok = false
		return
	}
	firstJourney := res.Journeys[0]
	if firstJourney.Type == routingv2.JourneyType_JOURNEY_TYPE_WALKING {
		if route := firstJourney.Walking.Route; len(route) != 0 {
			laneStart = r.Start.Aoi.WalkingLanes()[route[0].LaneId]
		}
	} else {
		for _, l := range r.Start.Aoi.WalkingLanes() {
			laneStart = l
		}
	}
	if laneStart == nil {
		r.route = make([]PedestrianSegment, 0)
		r.indexRoute = 0
		r.ok = false
		return
	}
	r.base = res
	r.indexJourney = -1
	r.NextJourney(laneStart)
	r.ok = true
}

func (r *PedestrianRoute) NextJourney(lane entity.ILane) bool {
	if r.indexJourney+1 >= len(r.base.Journeys) {
		return false
	}
	r.indexJourney++
	r.route = make([]PedestrianSegment, 0)
	r.indexRoute = 0
	pb := r.base.Journeys[r.indexJourney]
	switch pb.Type {
	case routingv2.JourneyType_JOURNEY_TYPE_WALKING:
		pbRoute := pb.Walking.Route
		if lane.ID() != pbRoute[0].LaneId {
			log.Panic("PedestrianRoute: wrong start lane when processing")
		}
		direction := pbRoute[0].MovingDirection
		r.route = append(r.route, PedestrianSegment{lane, direction})
		for _, segment := range pbRoute[1:] {
			laneID := segment.LaneId
			switch direction {
			case routingv2.MovingDirection_MOVING_DIRECTION_FORWARD:
				lt := lane.Successors()[laneID]
				lane = lt.Lane
				if lt.Type == mapv2.LaneConnectionType_LANE_CONNECTION_TYPE_HEAD {
					direction = routingv2.MovingDirection_MOVING_DIRECTION_FORWARD
				} else {
					direction = routingv2.MovingDirection_MOVING_DIRECTION_BACKWARD
				}
			case routingv2.MovingDirection_MOVING_DIRECTION_BACKWARD:
				lt := lane.Predecessors()[laneID]
				lane = lt.Lane
				if lt.Type == mapv2.LaneConnectionType_LANE_CONNECTION_TYPE_HEAD {
					direction = routingv2.MovingDirection_MOVING_DIRECTION_FORWARD
				} else {
					direction = routingv2.MovingDirection_MOVING_DIRECTION_BACKWARD
				}
			default:
				log.Panic("PedestrianRoute: unknown nextLaneType")
			}
			if lane == nil {
				log.Panicf("PedestrianRoute: lane %v in res does not exist", laneID)
			}
			if lane.ID() != laneID {
				log.Panicf("PedestrianRoute: lane %v in route response and %v do not match",
					laneID, lane.ID())
			}
			r.route = append(r.route, PedestrianSegment{lane, direction})
		}
		if lanePosition := r.end.LanePosition; lanePosition != nil {
			r.EndS = lanePosition.S
			r.EndAoi = nil
		} else {
			r.EndAoi = lane.Aois()[r.end.AoiPosition.AoiId]
			r.EndS = r.EndAoi.WalkingS(pb.Walking.Route[len(pb.Walking.Route)-1].LaneId)
		}
	default:
		log.Panic("PedestrianRoute: unsupported journeyType")
	}
	return true
}
