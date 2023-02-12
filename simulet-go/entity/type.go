package entity

import (
	"git.fiblab.net/sim/simulet-go/entity/person/schedule"
	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
	"git.fiblab.net/sim/simulet-go/utils/container"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"git.fiblab.net/sim/simulet-go/utils/randengine"
)

// 人的状态
type PersonStatus int32

const (
	PersonStatus_SLEEP               PersonStatus = iota // 未到出发时间——Aoi.sleep
	PersonStatus_CROWD                                   // 室内行走——Aoi.crowd
	PersonStatus_WAITING_FOR_LEAVING                     // 等待Aoi放行——Aoi.gate
	PersonStatus_WALKING                                 // 行人——Pedestrian
	PersonStatus_DRIVING                                 // 司机——Vehicle
)

// 方位常量
const (
	LEFT   = 0 // 左侧
	RIGHT  = 1 // 右侧
	BEFORE = 0 // 后方，等价于prev
	AFTER  = 1 // 前方，等价于next
)

type BaseRuntime struct {
	Position  geometry.Point // 位置
	Speed     float64        // 速度
	Direction float64        // 方向
}

type BaseRuntimeOnRoad struct {
	Lane ILane   // 所在车道id
	S    float64 // 车道上的位置
}

type BaseRuntimeInAoi struct {
	Aoi IAoi // 所在Aoi
}

type IPerson interface {
	// 自身属性
	ID() int32
	Attr() *agentv2.AgentAttribute
	VehicleAttr() *agentv2.VehicleAttribute
	BikeAttr() *agentv2.BikeAttribute
	IsOnBike() bool
	LaneS() float64
	Position() geometry.Point
	Speed() float64
	Direction() float64
	InJunction() bool
	Status() PersonStatus

	// 供aoi调用更新人的位置与状态

	SetRuntimeStatusByAoi(PersonStatus)
	SetSnapshotByAoi(BaseRuntime, BaseRuntimeInAoi)
	InitPedestrianByAoi()
	InitVehicleByAoi()
	SetCrowdByAoi()   // crowds.MarkAsAdded
	UnsetCrowdByAoi() // crowds.MarkAsRemoved
	SetRuntimeActivityByAoi(activity string)

	// 日程、导航相关
	Schedule() *schedule.Schedule
	CheckDeparture() bool
	ResetScheduleIfNeed()
	NextTrip() bool
	RequestRouteFromAoi(IAoi)
	RouteSuccessful() bool
	CurrentVehicleLaneID() int32
	CurrentPedestrianLaneID() int32
	NextPedestrianRouteJourney(ILane) bool

	// 随机数发生器
	Generator() *randengine.Engine
}

type IVehicle interface {
	ID() int32
	Length() float64
	Attr() *agentv2.AgentAttribute
	VehicleAttr() *agentv2.VehicleAttribute
	Speed() float64
	Lane() ILane
	S() float64
	Direction() float64
	Position() geometry.Point
	ShadowS() float64
	ShadowLane() ILane
	IsLaneChanging() bool
}

type IPedestrian interface {
	Lane() ILane
	S() float64
	Speed() float64
	Direction() float64
	Position() geometry.Point
	IsForward() bool
}

type Connection struct {
	Lane ILane
	Type mapv2.LaneConnectionType
}

type VehicleSideLink struct {
	// [LEFT/RIGHT][BACK/FRONT]
	Links [2][2]*container.ListNode[IVehicle, VehicleSideLink]
}

func (l *VehicleSideLink) Clear() {
	l.Links[0][0] = nil
	l.Links[0][1] = nil
	l.Links[1][0] = nil
	l.Links[1][1] = nil
}

type ILane interface {
	ID() int32
	Length() float64
	Width() float64

	ProjectFromLane(l ILane, s float64) float64

	Predecessors() map[int32]Connection
	Successors() map[int32]Connection
	Aois() map[int32]IAoi
	FirstLeftLane() ILane
	FirstRightLane() ILane
	GetPositionByS(s float64) geometry.Point
	GetDirectionByS(s float64) float64
	InJunction() bool

	// 获取特定位置车辆

	GetFirstVehicle() IVehicle
	Vehicles() *container.List[IVehicle, VehicleSideLink]
	Pedestrians() *container.List[IPedestrian, struct{}]

	// 道路状态

	MaxSpeed() float64
	LightState() lightv2.LightState
	LightStateRemainingTime() float64

	// Lane链表
	ReportVehicleAdded(node *container.ListNode[IVehicle, VehicleSideLink])
	ReportVehicleRemoved(node *container.ListNode[IVehicle, VehicleSideLink])
	ReportPedestrianAdded(node *container.ListNode[IPedestrian, struct{}])
	ReportPedestrianRemoved(node *container.ListNode[IPedestrian, struct{}])
}

type IJunction interface {
	ID() int32
}

// 车道的信控部分接口
type ILaneTrafficLightSetter interface {
	SetLightState(state lightv2.LightState)
	SetLightRemainingTime(time float64)
}

type AoiMoveType int32

const (
	AoiMoveType_UNSPECIFIED AoiMoveType = iota
	AoiMoveType_GATE
	AoiMoveType_SLEEP
	AoiMoveType_CROWD
	AoiMoveType_LANE
	AoiMoveType_INIT
)

type IAoi interface {
	// 自身属性
	ID() int32
	Positions() []geometry.Point
	Centroid() geometry.Point

	// 道路连接关系
	DrivingLanes() map[int32]ILane
	DrivingS(laneID int32) float64
	WalkingLanes() map[int32]ILane
	WalkingS(laneID int32) float64
	LaneSs() map[int32]float64

	Add(person IPerson, from, to AoiMoveType, laneID int32)
	ReportRemoved(person IPerson, from, to AoiMoveType, laneID int32)
	// 模块间流转
	MoveBetweenSubmodules(person IPerson, from, to AoiMoveType, laneID int32)
}
