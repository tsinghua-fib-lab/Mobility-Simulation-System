package route

import (
	"context"
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
	routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"
)

type VehicleSegment struct {
	Lane         entity.ILane           // 车道指针
	NextLaneType routingv2.NextLaneType // 当前车道前往下一车道的方式
}

func (s VehicleSegment) IsForward() bool {
	return s.NextLaneType == routingv2.NextLaneType_NEXT_LANE_TYPE_FORWARD
}

func (s VehicleSegment) IsLaneChange() bool {
	return s.NextLaneType == routingv2.NextLaneType_NEXT_LANE_TYPE_LEFT ||
		s.NextLaneType == routingv2.NextLaneType_NEXT_LANE_TYPE_RIGHT
}

// 路径规划结果指针化，主要处理车辆
type VehicleRoute struct {
	Start         RouteStartPosition // 导航起点
	RerouteFlag   bool               // 路径更新请求，为true时，车辆需要在update最后发出路径规划请求
	ok            bool               // 导航请求是否成功
	size          int                // 总长度
	route         []VehicleSegment   // 转换为指针形式后的路径，满足route[0]==parent_lane
	distanceToEnd []float64          // 路段起点距终点的距离
	end           *geov2.Position    // 记录路径规划终点
}

func NewVehicleRoute() *VehicleRoute {
	return &VehicleRoute{
		route:         make([]VehicleSegment, 0),
		distanceToEnd: make([]float64, 0),
	}
}

func (r *VehicleRoute) Ok() bool {
	return r.ok
}

func (r *VehicleRoute) Current() VehicleSegment {
	return r.route[0]
}

func (r *VehicleRoute) Next() VehicleSegment {
	return r.route[1]
}

func (r *VehicleRoute) Last() VehicleSegment {
	return r.route[r.size-1]
}

func (r *VehicleRoute) ShadowCurrent() VehicleSegment {
	r.canShadow()
	return r.route[1]
}

func (r *VehicleRoute) ShadowNext() VehicleSegment {
	r.canShadow()
	return r.route[2]
}

// 是否在最后一个车道
func (r *VehicleRoute) AtLast() bool {
	return r.size <= 1
}

// 向前增加index，返回是否正常（true: 正常, false：越界）
func (r *VehicleRoute) Step() bool {
	if r.size == 1 {
		return false
	}
	r.route = r.route[1:]
	r.size--
	return true
}

func (r *VehicleRoute) GetDistanceToEnd(s float64) float64 {
	return math.Max(r.distanceToEnd[0]-s, 0)
}

func (r *VehicleRoute) CountContinuousLaneChanges() int {
	count := 0
	for i := 0; i < r.size; i++ {
		if r.route[i].IsLaneChange() {
			count++
		} else {
			break
		}
	}
	return count
}

func (r *VehicleRoute) Refresh(newRoute []VehicleSegment) {
	r.route = newRoute
	r.size = len(r.route)
	r.initDistanceToEnd()
}

// 计算终点AOI，如果终点不是AOI则返回nil
func (r *VehicleRoute) GetEndAoi() entity.IAoi {
	if aoiPos := r.end.AoiPosition; aoiPos != nil {
		endAoi, ok := r.route[len(r.route)-1].Lane.Aois()[aoiPos.AoiId]
		if !ok {
			log.Panicf(
				"VehicleRoute: lane %v does not contain aoi %v",
				r.route[len(r.route)-1].Lane.ID(),
				endAoi.ID(),
			)
		}
		return endAoi
	}
	return nil
}

func (r *VehicleRoute) GetEndS() float64 {
	return r.distanceToEnd[len(r.distanceToEnd)-1]
}

func (r *VehicleRoute) ProduceRouting(target *geov2.Position, startPosition RouteStartPosition) {
	r.Start = startPosition
	r.ok = false
	req := &routingv2.GetRouteRequest{
		Type: routingv2.RouteType_ROUTE_TYPE_DRIVING,
	}
	if r.Start.Lane != nil {
		req.Start = &geov2.Position{
			LanePosition: &geov2.LanePosition{
				LaneId: r.Start.Lane.ID(),
				S:      r.Start.S,
			},
		}
	} else {
		if r.Start.Aoi == nil {
			log.Panic("VehicleRoute: start position should contain aoi or lane")
		}
		req.Start = &geov2.Position{
			AoiPosition: &geov2.AoiPosition{
				AoiId: r.Start.Aoi.ID(),
			},
		}
	}
	req.End, r.end = target, target
	// 发送请求
	RoutingClient.GetRoute(context.Background(), req, r.ProcessRouting)
}

func (r *VehicleRoute) ProcessRouting(res *routingv2.GetRouteResponse) {
	if len(res.Journeys) == 0 {
		r.ok = false
		return
	}
	var laneStart entity.ILane
	if r.Start.Lane != nil {
		laneStart = r.Start.Lane
	} else {
		laneID := res.Journeys[0].Driving.Route[0].LaneId
		if drivingLane, ok := r.Start.Aoi.DrivingLanes()[laneID]; !ok {
			log.Panicf("VehicleRoute: no driving lane %v in Aoi %v", laneID, r.Start.Aoi.ID())
		} else {
			laneStart = drivingLane
		}
	}
	r.route = make([]VehicleSegment, 0)
	lane := laneStart
	// res check
	if !(len(res.Journeys) == 1 &&
		res.Journeys[0].Type == *routingv2.JourneyType_JOURNEY_TYPE_DRIVING.Enum() &&
		res.Journeys[0].Driving != nil &&
		len(res.Journeys[0].Driving.Route) > 0) {
		log.Panic("VehicleRoute: wrong res")
	}
	pb := res.Journeys[0].Driving.Route
	r.size = len(pb)
	// 将导航响应转换为指针
	if pb[0].LaneId != lane.ID() {
		log.Panic("VehicleRoute: wrong start lane when processing")
	}
	nextType := pb[0].NextLaneType
	r.route = append(r.route, VehicleSegment{lane, nextType})
	for i, segment := range pb[1:] {
		laneID := segment.LaneId
		switch nextType {
		case routingv2.NextLaneType_NEXT_LANE_TYPE_FORWARD:
			lane = lane.Successors()[laneID].Lane
		case routingv2.NextLaneType_NEXT_LANE_TYPE_LEFT:
			lane = lane.FirstLeftLane()
		case routingv2.NextLaneType_NEXT_LANE_TYPE_RIGHT:
			lane = lane.FirstRightLane()
		case routingv2.NextLaneType_NEXT_LANE_TYPE_LAST:
			if i+1 != r.size {
				log.Panic("VehicleRoute: detect NEXT_LANE_TYPE_UNKNOWN in route response")
			}
		default:
			log.Panic("VehicleRoute: unknown PbNextLaneType")
		}
		if lane == nil {
			log.Panicf("VehicleRoute: lane %v in res does not exist", laneID)
		}
		if lane.ID() != laneID {
			log.Panicf("VehicleRoute: lane %v in route response and %v do not match",
				laneID, lane.ID())
		}
		nextType = segment.NextLaneType
		r.route = append(r.route, VehicleSegment{lane, nextType})
	}
	if r.route[len(r.route)-1].NextLaneType != routingv2.NextLaneType_NEXT_LANE_TYPE_LAST {
		log.Panic("VehicleRoute: the last lane's NextLaneType should be NEXT_LANE_TYPE_LAST")
	}
	// 计算到终点的距离
	r.initDistanceToEnd()
	r.ok = true
}
