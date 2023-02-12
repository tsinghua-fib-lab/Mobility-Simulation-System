 #include "wolong/map/v2/map.pb.h"
 #include "wolong/agent//v2//agent.pb.h"
 #include "wolong/geo/v2/geo.pb.h"
 #include "wolong/routing/v2/routing.pb.h"
 #include "wolong/routing/v2/routing_service.pb.h"
 #include "wolong/traffic_light/v2/traffic_light.pb.h"
 #include "wolong/trip/v2/trip.pb.h"
#include "wolong/routing/v2/routing_service.grpc.pb.h"

namespace simulet {
    //地图
    using PAoi = wolong::map::v2::Aoi;
    using PMap = wolong::map::v2::Map;
    using PLane = wolong::map::v2::Lane;
    using PRoad = wolong::map::v2::Road;
    using Pjunction = wolong::map::v2::Junction;
    using Pagents = wolong::agent::v2::Agents;
    using Agent = wolong::agent::v2::Agent;


    //导航
    using PbGetRouteRequest = wolong::routing::v2::GetRouteRequest;
    using PbGetRouteResponse = wolong::routing::v2::GetRouteResponse;
    using PbNextLaneType = wolong::routing::v2::NextLaneType;
    using PbRouteType = wolong::routing::v2::RouteType;
    using PbJourneyType = ::wolong::routing::v2::JourneyType;
    //using AgentPbGetRouteRequest =wolong::routing::v2::getr

    using PbPosition = ::wolong::geo::v2::Position;
}