package person

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/person/pedestrian"
	"git.fiblab.net/sim/simulet-go/entity/person/route"
	"git.fiblab.net/sim/simulet-go/entity/person/schedule"
	"git.fiblab.net/sim/simulet-go/entity/person/vehicle"
	agentv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/agent/v2"
	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
	routingv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/routing/v2"
	tripv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/trip/v2"
	"git.fiblab.net/sim/simulet-go/utils/clock"
	"git.fiblab.net/sim/simulet-go/utils/container"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"git.fiblab.net/sim/simulet-go/utils/randengine"
	"github.com/samber/lo"
)

const (
	MAX_NOISE_ON_VEHICLE_SPEED       = 5  // 车辆速度随机扰动最大值
	MAX_NOISE_ON_VEHICLE_ACC         = .5 // 车辆加速度随机扰动最大值s
	MAX_NOISE_ON_PEDESTRIAN_POSITION = 2  // 行人位置输出随机扰动最大值
)

type Snapshot struct {
	// 上一时刻状态，室外由runtime提供，室内由SetSnapshotByAoi触发从runtime复制
	Status entity.PersonStatus

	// 供输出或外部接口调用的人的数据快照，与status对应
	entity.BaseRuntime
	entity.BaseRuntimeOnRoad
	entity.BaseRuntimeInAoi

	// 活动输出
	Activity string
}

type Runtime struct {
	// 当前时刻状态，通过SetRuntimeStatusByAoi更新（室外不主动更新）
	Status entity.PersonStatus
	// 活动输出，由Aoi通过SetRuntimeActivityByAoi更新
	Activity string

	// 接口对于室外的人为车辆、行人、乘客其一,且与status对应
	// 对于室内的人为nil
	submodule submodule
}

type Person struct {
	container.ActiveElement

	id          int32
	attr        *agentv2.AgentAttribute   // 人的属性
	vehicleAttr *agentv2.VehicleAttribute // 车的属性
	bikeAttr    *agentv2.BikeAttribute    // 自行车的属性
	home        *geov2.Position           // 人的家

	runtime  Runtime  // 运行时数据
	snapshot Snapshot // 快照

	schedule          *schedule.Schedule // 时刻表
	newSchedule       []*tripv2.Schedule // schedule修改buffer
	scheduleResetFlag bool               // 时刻表是否被修改

	vehicleRoute    *route.VehicleRoute    // 车辆导航
	pedestrianRoute *route.PedestrianRoute // 行人导航

	generator              *randengine.Engine // 随机数生成器，以ID为seed
	pedestrianOutputOffset geometry.Point     // 为行人输出添加的随机扰动
}

func NewPerson(base *agentv2.Agent, m *PersonManager) *Person {
	p := &Person{
		id:              base.Id,
		attr:            base.Attribute,
		vehicleAttr:     base.VehicleAttribute,
		bikeAttr:        base.BikeAttribute,
		home:            base.Home,
		schedule:        schedule.NewSchedule(),
		newSchedule:     make([]*tripv2.Schedule, 0),
		vehicleRoute:    route.NewVehicleRoute(),
		generator:       randengine.New(uint64(base.Id)),
		pedestrianRoute: route.NewPedestrianRoute(),
	}
	p.SetSchedules(base.GetSchedules())
	// 行人输出添加随机扰动
	p.pedestrianOutputOffset = geometry.Point{
		X: lo.Clamp(p.generator.NormFloat64(), -MAX_NOISE_ON_PEDESTRIAN_POSITION,
			MAX_NOISE_ON_PEDESTRIAN_POSITION),
		Y: lo.Clamp(p.generator.NormFloat64(), -MAX_NOISE_ON_PEDESTRIAN_POSITION,
			MAX_NOISE_ON_PEDESTRIAN_POSITION),
	}
	// 为车辆属性添加随机扰动
	// 最大速度
	p.attr.MaxSpeed = math.Max(p.attr.MaxSpeed+
		lo.Clamp(.5*MAX_NOISE_ON_VEHICLE_SPEED*p.generator.NormFloat64(),
			-MAX_NOISE_ON_VEHICLE_SPEED,
			MAX_NOISE_ON_VEHICLE_SPEED),
		.1)
	// 最大加速度
	p.attr.MaxAcceleration = math.Max(p.attr.MaxAcceleration+
		lo.Clamp(.5*MAX_NOISE_ON_VEHICLE_ACC*p.generator.NormFloat64(),
			-MAX_NOISE_ON_VEHICLE_ACC,
			MAX_NOISE_ON_VEHICLE_ACC),
		.1)
	// 最大刹车加速度
	p.attr.MaxBrakingAcceleration = math.Min(p.attr.MaxBrakingAcceleration+
		lo.Clamp(.5*MAX_NOISE_ON_VEHICLE_ACC*p.generator.NormFloat64(),
			-MAX_NOISE_ON_VEHICLE_ACC,
			MAX_NOISE_ON_VEHICLE_ACC),
		-.1)
	// 通常加速度
	p.attr.UsualAcceleration = math.Max(p.attr.UsualAcceleration+
		lo.Clamp(.5*MAX_NOISE_ON_VEHICLE_ACC*p.generator.NormFloat64(),
			-MAX_NOISE_ON_VEHICLE_ACC,
			MAX_NOISE_ON_VEHICLE_ACC),
		.1)
	// 通常刹车加速度
	p.attr.UsualBrakingAcceleration = math.Min(p.attr.UsualBrakingAcceleration+
		lo.Clamp(.5*MAX_NOISE_ON_VEHICLE_ACC*p.generator.NormFloat64(),
			-MAX_NOISE_ON_VEHICLE_ACC,
			MAX_NOISE_ON_VEHICLE_ACC,
		),
		-.1)
	return p
}

// 根据各种子模块更新快照（室外）
func (p *Person) Prepare() {
	switch p.runtime.Status {
	case entity.PersonStatus_WALKING,
		entity.PersonStatus_DRIVING:
		sub := p.runtime.submodule
		sub.Prepare()
		rt, rtOnRoad := sub.FetchBaseSnapshotForPerson()
		p.snapshot = Snapshot{
			Status:            p.runtime.Status,
			BaseRuntime:       rt,
			BaseRuntimeOnRoad: rtOnRoad,
			Activity:          p.runtime.Activity,
		}
	default:
		log.Panicf("bad person %d status %v when prepare outdoor", p.ID(), p.runtime.Status)
	}
	// 优先执行新的schedule
	p.ResetScheduleIfNeed()
}

// 更新（室外）
func (p *Person) Update(stepInterval float64) {
	switch p.runtime.Status {
	case entity.PersonStatus_WALKING:
		p.runtime.submodule.Update(stepInterval)
		if aoi, endLane, ok := p.runtime.submodule.GetEndByPerson(); ok {
			// 行人结束路面行为（生命周期结束）的后处理
			// 本行程走完，进入sleep
			p.schedule.NextTrip(clock.GlobalTime)
			aoi.Add(
				p,
				entity.AoiMoveType_LANE,
				entity.AoiMoveType_SLEEP,
				endLane.ID(),
			)
			// 进入AOI，将人标记为移除出行人列表
			Manager.pedestrians.MarkAsRemoved(p)
			p.runtime.submodule = nil
		}
	case entity.PersonStatus_DRIVING:
		p.runtime.submodule.Update(stepInterval)
		if aoi, endLane, ok := p.runtime.submodule.GetEndByPerson(); ok {
			// 车辆结束路面行为（生命周期结束）的后处理
			// 对于驾车出行，一个trip对应一个journey，因此结束后开始下一trip
			p.schedule.NextTrip(clock.GlobalTime)
			aoi.Add(
				p,
				entity.AoiMoveType_LANE,
				entity.AoiMoveType_SLEEP,
				endLane.ID(),
			)
			Manager.vehicles.MarkAsRemoved(p)
			p.runtime.submodule = nil
		}
	default:
		log.Panicf("unknown person %d status %v when update", p.ID(), p.runtime.Status)
	}
}

func (p *Person) ID() int32 {
	return p.id
}

func (p *Person) Attr() *agentv2.AgentAttribute {
	return p.attr
}

func (p *Person) VehicleAttr() *agentv2.VehicleAttribute {
	return p.vehicleAttr
}

func (p *Person) BikeAttr() *agentv2.BikeAttribute {
	return p.bikeAttr
}

func (p *Person) Position() geometry.Point {
	return p.snapshot.Position
}

func (p *Person) Speed() float64 {
	return p.snapshot.Speed
}

func (p *Person) Direction() float64 {
	return p.snapshot.Direction
}

func (p *Person) InJunction() bool {
	return (p.snapshot.Status == entity.PersonStatus_WALKING ||
		p.snapshot.Status == entity.PersonStatus_DRIVING) && p.snapshot.Lane.InJunction()
}

func (p *Person) LaneS() float64 {
	return p.snapshot.S
}

func (p *Person) Status() entity.PersonStatus {
	return p.snapshot.Status
}

func (p *Person) Home() *geov2.Position {
	return p.home
}

// 供子模块上报人的当前status用
func (p *Person) SetRuntimeStatusByAoi(status entity.PersonStatus) {
	p.runtime.Status = status
}

func (p *Person) SetSnapshotByAoi(baseRt entity.BaseRuntime, baseRta entity.BaseRuntimeInAoi) {
	p.snapshot.Status = p.runtime.Status
	p.snapshot.BaseRuntime = baseRt
	p.snapshot.BaseRuntimeInAoi = baseRta
}

func (p *Person) SetRuntimeActivityByAoi(activity string) {
	p.runtime.Activity = activity
}

func (p *Person) Schedule() *schedule.Schedule {
	return p.schedule
}

func (p *Person) SetSchedules(schedules []*tripv2.Schedule) {
	p.newSchedule = schedules
	p.scheduleResetFlag = true
}

func (p *Person) ResetScheduleIfNeed() {
	if p.scheduleResetFlag {
		p.schedule.Set(p.newSchedule, clock.GlobalTime)
		p.scheduleResetFlag = false
	}
}

func (p *Person) NextTrip() bool {
	return p.schedule.NextTrip(clock.GlobalTime)
}

func (p *Person) IsOnBike() bool {
	return false
}

func (p *Person) CurrentVehicleLaneID() int32 {
	return p.vehicleRoute.Current().Lane.ID()
}

func (p *Person) CurrentPedestrianLaneID() int32 {
	return p.pedestrianRoute.Current().Lane.ID()
}

func (p *Person) NextPedestrianRouteJourney(walkingLane entity.ILane) bool {
	return p.pedestrianRoute.NextJourney(walkingLane)
}

// 检查是否到达出发时间
func (p *Person) CheckDeparture() bool {
	return clock.GlobalTime >= p.schedule.GetDepartureTime()
}

// 通知person初始化行人模块
// 函数调用方保证行人导航可用
func (p *Person) InitPedestrianByAoi() {
	p.runtime.submodule = pedestrian.NewPedestrian(
		p,
		p.pedestrianRoute,
	)
	Manager.pedestrians.MarkAsAdded(p)
}

// 通知person初始化车辆模块
// 函数调用方保证车辆导航可用
func (p *Person) InitVehicleByAoi() {
	p.runtime.submodule = vehicle.NewVehicle(p, p.vehicleRoute)
	Manager.vehicles.MarkAsAdded(p)
}

func (p *Person) SetCrowdByAoi() {
	Manager.crowds.MarkAsAdded(p)
}

func (p *Person) UnsetCrowdByAoi() {
	Manager.crowds.MarkAsRemoved(p)
}

// 以aoi为起点发出导航请求
func (p *Person) RequestRouteFromAoi(aoi entity.IAoi) {
	trip := p.schedule.GetTrip()
	if trip.Mode == tripv2.TripMode_TRIP_MODE_DRIVE_ONLY {
		p.vehicleRoute.ProduceRouting(trip.End, route.RouteStartPosition{Aoi: aoi})
	} else {
		p.pedestrianRoute.ProduceRouting(
			trip.End,
			route.RouteStartPosition{Aoi: aoi},
			routingv2.RouteType_ROUTE_TYPE_WALKING,
		)
	}
}

// 导航请求是否成功,成功则返回true，否则转到下一trip并返回false
func (p *Person) RouteSuccessful() bool {
	trip := p.schedule.GetTrip()
	if trip.Mode == tripv2.TripMode_TRIP_MODE_DRIVE_ONLY {
		if p.vehicleRoute.Ok() {
			return true
		}
	} else {
		if p.pedestrianRoute.Ok() {
			return true
		}
	}
	p.schedule.NextTrip(clock.GlobalTime)
	return false
}

func (p *Person) Generator() *randengine.Engine {
	return p.generator
}
