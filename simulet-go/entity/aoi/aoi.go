package aoi

import (
	"sync/atomic"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/aoi/crowd"
	"git.fiblab.net/sim/simulet-go/entity/aoi/gate"
	"git.fiblab.net/sim/simulet-go/entity/aoi/sleep"
	"git.fiblab.net/sim/simulet-go/entity/lane"
	geov2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/geo/v2"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"github.com/samber/lo"
)

type Runtime struct {
	AddCount    atomic.Int32 // 上个step进入的人数
	RemoveCount atomic.Int32 // 上个step离开的人数
	HeadCount   int32        // 记录aoi内总人数
}

type Aoi struct {
	id        int32
	positions []geometry.Point
	area      *float64
	centroid  geometry.Point

	// 初始化临时值
	initDrivingPositions []*geov2.LanePosition
	initWalkingPositions []*geov2.LanePosition

	laneSs       map[int32]float64      // aoi连接的车道id到对应道路上位置的映射
	drivingLanes map[int32]entity.ILane // 对应的行车路网车道指针
	walkingLanes map[int32]entity.ILane // 对应的步行路网车道指针

	// 子模块
	sleep *sleep.Sleep // 休眠人群管理
	gate  *gate.Gate   // 人车放行
	crowd *crowd.Crowd // 室内人流模拟

	// 统计
	runtime  Runtime
	snapshot Runtime
}

func NewAoi(base *mapv2.Aoi, m *AoiManager) *Aoi {
	aoi := &Aoi{
		id: base.Id,
		positions: lo.Map(base.Positions, func(p *geov2.XYPosition, _ int) geometry.Point {
			return geometry.NewPointFromPb(p)
		}),
		area:                 base.Area,
		initDrivingPositions: base.DrivingPositions,
		initWalkingPositions: base.WalkingPositions,
		laneSs:               make(map[int32]float64),
		drivingLanes:         make(map[int32]entity.ILane),
		walkingLanes:         make(map[int32]entity.ILane),
	}
	aoi.centroid = geometry.GetPolygonCentroid(aoi.positions)

	// 初始化子模块

	// 每个aoi都有
	aoi.sleep = sleep.New(aoi)
	aoi.gate = gate.New(aoi, base)
	// 根据微观范围决定是否激活室内模拟模块
	if base.Area != nil {
		aoi.crowd = crowd.New(base, aoi)
	}
	return aoi
}

func (a *Aoi) InitLanes(laneManager *lane.LaneManager) {
	for _, position := range a.initDrivingPositions {
		if lane, err := laneManager.Get(position.LaneId); err != nil {
			log.Panic(err)
		} else {
			a.drivingLanes[lane.ID()] = lane
			a.laneSs[lane.ID()] = position.S
			lane.AddAoiWhenInit(a)
		}
	}
	for _, position := range a.initWalkingPositions {
		if lane, err := laneManager.Get(position.LaneId); err != nil {
			log.Panic(err)
		} else {
			a.walkingLanes[lane.ID()] = lane
			a.laneSs[lane.ID()] = position.S
			lane.AddAoiWhenInit(a)
		}
	}
	a.initDrivingPositions = nil
	a.initWalkingPositions = nil
}

func (a *Aoi) Prepare() {
	a.gate.Prepare()
	a.sleep.Prepare()
	a.crowd.Prepare()

	a.runtime.HeadCount = a.runtime.HeadCount + a.runtime.AddCount.Load() - a.runtime.RemoveCount.Load()
	a.snapshot = a.runtime
	a.runtime.AddCount.Store(0)
	a.runtime.RemoveCount.Store(0)
}

func (a *Aoi) Update(stepInterval float64) {
	a.gate.Update()
	a.sleep.Update()
	a.crowd.Update(stepInterval)
}

func (a *Aoi) ID() int32 {
	return a.id
}

func (a *Aoi) Positions() []geometry.Point {
	return a.positions
}

func (a *Aoi) Centroid() geometry.Point {
	return a.centroid
}

func (a *Aoi) DrivingLanes() map[int32]entity.ILane {
	return a.drivingLanes
}

func (a *Aoi) WalkingLanes() map[int32]entity.ILane {
	return a.walkingLanes
}

func (a *Aoi) LaneSs() map[int32]float64 {
	return a.laneSs
}

func (a *Aoi) DrivingS(laneID int32) float64 {
	return a.laneSs[laneID]
}

func (a *Aoi) WalkingS(laneID int32) float64 {
	return a.laneSs[laneID]
}

func (a *Aoi) EnableIndoor() bool {
	return a.crowd != nil
}

func (a *Aoi) HeadCount() int32 {
	return a.snapshot.HeadCount
}

func (a *Aoi) Add(p entity.IPerson, from, to entity.AoiMoveType, laneID int32) {
	a.runtime.AddCount.Add(1)
	switch from {
	case entity.AoiMoveType_INIT:
		switch to {
		case entity.AoiMoveType_SLEEP:
			a.sleep.Add(p)
			return
		}
	case entity.AoiMoveType_LANE:
		switch to {
		case entity.AoiMoveType_SLEEP:
			if a.EnableIndoor() {
				a.crowd.AddFromLaneToSleep(p, laneID)
				return
			} else {
				a.sleep.Add(p)
				return
			}
		}
	}
	log.Panicf("fail to handle ADD from %v to %v for person %d", from, to, p.ID())
}

func (a *Aoi) ReportRemoved(p entity.IPerson, from, to entity.AoiMoveType, laneID int32) {
	a.runtime.RemoveCount.Add(1)
}

func (a *Aoi) MoveBetweenSubmodules(p entity.IPerson, from, to entity.AoiMoveType, laneID int32) {
	switch from {
	case entity.AoiMoveType_CROWD:
		switch to {
		case entity.AoiMoveType_GATE:
			a.gate.Add(p)
			return
		case entity.AoiMoveType_SLEEP:
			a.sleep.Add(p)
			return
		}
	case entity.AoiMoveType_SLEEP:
		if a.EnableIndoor() {
			a.crowd.AddFromSleep(p, to, laneID)
			return
		} else {
			switch to {
			case entity.AoiMoveType_GATE:
				a.gate.Add(p)
				return
			}
		}
	}
	log.Panicf("fail to handle MOVE from %v to %v for person %d", from, to, p.ID())
}
