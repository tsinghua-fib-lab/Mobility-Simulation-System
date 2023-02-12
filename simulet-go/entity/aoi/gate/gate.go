package gate

import (
	"sync"

	"git.fiblab.net/sim/simulet-go/entity"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	tripv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/trip/v2"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
)

type Gate struct {
	aoi           entity.IAoi              // 所在Aoi
	walkingLanes  map[int32]entity.ILane   // 人行道
	drivingLanes  map[int32]entity.ILane   // 车道
	laneSs        map[int32]float64        // Aoi与道路连接处位置
	gatePositions map[int32]geometry.Point // gate点的坐标

	personInserted       []entity.IPerson // 新加入的人
	personInsertMutex    sync.Mutex
	pedestrianForLeaving []entity.IPerson // 待放行的行人
	vehicleForLeaving    []entity.IPerson // 待放行的车辆
}

func New(aoi entity.IAoi, base *mapv2.Aoi) *Gate {
	g := &Gate{
		aoi:                  aoi,
		walkingLanes:         aoi.WalkingLanes(),
		drivingLanes:         aoi.DrivingLanes(),
		laneSs:               aoi.LaneSs(),
		gatePositions:        make(map[int32]geometry.Point),
		personInserted:       make([]entity.IPerson, 0),
		personInsertMutex:    sync.Mutex{},
		pedestrianForLeaving: make([]entity.IPerson, 0),
		vehicleForLeaving:    make([]entity.IPerson, 0),
	}
	for i, lanePos := range base.DrivingPositions {
		g.gatePositions[lanePos.LaneId] = geometry.NewPointFromPb(base.DrivingGates[i])
	}
	for i, lanePos := range base.WalkingPositions {
		g.gatePositions[lanePos.LaneId] = geometry.NewPointFromPb(base.WalkingGates[i])
	}
	return g
}

func (g *Gate) Prepare() {
	for _, p := range g.personInserted {
		// 检查人的出行方式
		if p.Schedule().GetTrip().Mode == tripv2.TripMode_TRIP_MODE_DRIVE_ONLY {
			p.SetSnapshotByAoi(entity.BaseRuntime{
				Position: g.gatePositions[p.CurrentVehicleLaneID()],
			}, entity.BaseRuntimeInAoi{
				Aoi: g.aoi,
			})
			g.vehicleForLeaving = append(g.vehicleForLeaving, p)
		} else {
			p.SetSnapshotByAoi(entity.BaseRuntime{
				Position: g.gatePositions[p.CurrentPedestrianLaneID()],
			}, entity.BaseRuntimeInAoi{
				Aoi: g.aoi,
			})
			g.pedestrianForLeaving = append(g.pedestrianForLeaving, p)
		}
	}
	g.personInserted = []entity.IPerson{}
}

func (g *Gate) Update() {
	// 行人直接放行，不限制数量
	for _, p := range g.pedestrianForLeaving {
		p.SetRuntimeStatusByAoi(entity.PersonStatus_WALKING)
		p.InitPedestrianByAoi()
		g.aoi.ReportRemoved(
			p,
			entity.AoiMoveType_GATE,
			entity.AoiMoveType_LANE,
			p.CurrentPedestrianLaneID(),
		)
	}
	g.pedestrianForLeaving = []entity.IPerson{}
	if len(g.vehicleForLeaving) > 0 {
		// 车辆检查路口情况
		p := g.vehicleForLeaving[0]
		g.vehicleForLeaving = g.vehicleForLeaving[1:]
		p.SetRuntimeStatusByAoi(entity.PersonStatus_DRIVING)
		p.InitVehicleByAoi()
		g.aoi.ReportRemoved(
			p,
			entity.AoiMoveType_GATE,
			entity.AoiMoveType_LANE,
			p.CurrentVehicleLaneID(),
		)
	}
}

func (g *Gate) Add(p entity.IPerson) {
	p.SetRuntimeStatusByAoi(entity.PersonStatus_WAITING_FOR_LEAVING)
	g.personInsertMutex.Lock()
	defer g.personInsertMutex.Unlock()
	g.personInserted = append(g.personInserted, p)
}

func (g *Gate) HeadCount() int32 {
	return int32(len(g.vehicleForLeaving) + len(g.pedestrianForLeaving))
}
