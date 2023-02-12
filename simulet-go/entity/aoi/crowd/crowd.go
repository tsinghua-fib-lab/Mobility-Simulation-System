package crowd

import (
	"math"
	"sync"

	"git.fiblab.net/sim/simulet-go/entity"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"git.fiblab.net/sim/simulet-go/utils/parallel"
	"git.fiblab.net/sim/simulet-go/utils/randengine"
)

const (
	DEFAULT_DESIRED_SPEED = 1.34 // 默认期望速度
	GATE_SIZE             = 3    // 出入口宽度
	TRY_GET_RANDOM_POINT  = 30   // 随机取点最大尝试次数
	STEP_INTERVAL_IN_DOOR = .1   // 室内行人模拟的更新间隔
)

// 数据快照
type Snapshot struct {
	Status   Status         // 状态
	Position geometry.Point // 位置
	Velocity geometry.Point // 速度
}

type Crowd struct {
	aoi entity.IAoi // 所在Aoi

	gates          map[int32]geometry.Point // 出入口列表（包括driving_gates和walking_gates）
	sleepPoints    []geometry.Point         // 楼电梯间位置（行人走到此处后进入sleep模块，不再有计算任务）
	boundaryPoints []geometry.Point         // Aoi 边界点列表。各点顺序给出，注意第一点与最后一点相同
	triangleAreas  []float64                // 将 Aoi 进行三角划分后的 n-2 个三角形的面积，用于在 Aoi 中随机取点

	generator *randengine.Engine // 随机数生成器

	personInserted    []*Runtime // 新加入的人
	personInsertMutex sync.Mutex
	persons           []*Runtime // 待更新的室内行人
}

func New(base *mapv2.Aoi, aoi entity.IAoi) *Crowd {
	crowd := &Crowd{
		aoi:               aoi,
		gates:             make(map[int32]geometry.Point),
		sleepPoints:       make([]geometry.Point, 0),
		boundaryPoints:    make([]geometry.Point, 0),
		triangleAreas:     make([]float64, 0),
		generator:         randengine.New(uint64(base.Id)),
		personInserted:    make([]*Runtime, 0),
		personInsertMutex: sync.Mutex{},
		persons:           make([]*Runtime, 0),
	}
	for i, p := range base.WalkingGates {
		crowd.gates[base.WalkingPositions[i].LaneId] = geometry.NewPointFromPb(p)
	}
	for i, p := range base.DrivingGates {
		crowd.gates[base.DrivingPositions[i].LaneId] = geometry.NewPointFromPb(p)
	}
	for _, p := range base.Positions {
		crowd.boundaryPoints = append(crowd.boundaryPoints, geometry.NewPointFromPb(p))
	}
	// 计算以 p0 为顶点对 Aoi 区域进行三角划分后的 n-2 个三角形的面积
	p0 := crowd.boundaryPoints[0]
	for i, p1 := range crowd.boundaryPoints[1 : len(crowd.boundaryPoints)-2] {
		p2 := crowd.boundaryPoints[i+2]
		crowd.triangleAreas = append(crowd.triangleAreas,
			.5*math.Abs(geometry.Cross(p1.Sub(p0), p2.Sub(p0))))
	}
	// 每个aoi内随机选择若干点作为楼电梯间位置
	if points := crowd.boundaryPoints; len(points) > 1 {
		for i, j := 0, len(points)-2; i < len(points)-1; i++ {
			p1, p2, p3 := points[j], points[i], points[i+1]
			p12, p23 := geometry.Distance(p1, p2), geometry.Distance(p2, p3)
			p12, p23 = math.Max(p12, 1e-5), math.Max(p23, 1e-5)
			sin := geometry.Cross(p2.Sub(p1), p3.Sub(p2)) / p12 / p23
			if sin < -.98 {
				p := p2.Add(p1.Sub(p2).Scale(5 / p12)).Add(p3.Sub(p2).Scale(5 / p23))
				if p.InPolygon(points) {
					crowd.sleepPoints = append(crowd.sleepPoints, p)
				}
			}
			j = i
		}
		for i := len(crowd.sleepPoints); i < 3; i++ {
			crowd.sleepPoints = append(crowd.sleepPoints, crowd.getRandomPosition())
		}
	} else {
		crowd.sleepPoints = append(crowd.sleepPoints, points[0])
	}
	return crowd
}

func (c *Crowd) Prepare() {
	if c == nil {
		return
	}
	// 新插入的人加入待更新人群
	c.persons = append(c.persons, c.personInserted...)
	c.personInserted = []*Runtime{}
	// 更新人的位置
	for _, runtime := range c.persons {
		// 维护人的快照信息
		runtime.Person.SetSnapshotByAoi(entity.BaseRuntime{
			Position:  runtime.Position,
			Speed:     runtime.Velocity.Length(),
			Direction: runtime.Velocity.Angle(),
		}, entity.BaseRuntimeInAoi{
			Aoi: c.aoi,
		})
	}
}

func (c *Crowd) Update(stepInterval float64) {
	if c == nil {
		return
	}
	for steps := int(stepInterval / STEP_INTERVAL_IN_DOOR); steps >= 0; steps-- {
		// 快照准备
		personPositionMap := make(map[int32]Snapshot, len(c.persons))
		for _, runtime := range c.persons {
			personPositionMap[runtime.Person.ID()] = Snapshot{
				Status:   runtime.Status,
				Position: runtime.Position,
				Velocity: runtime.Velocity,
			}
		}
		// 更新一步
		parallel.GoFor(c.persons, 1, func(runtime *Runtime) {
			runtime.Update(personPositionMap, c)
		})
	}
	// 检查人是否到终点并进行人的流转
	for i := 0; i < len(c.persons); {
		rt := c.persons[i]
		if rt.Status == Status_REACH_TARGET {
			var to entity.AoiMoveType
			switch rt.Interest {
			case InterestType_EXIT:
				to = entity.AoiMoveType_GATE
			case InterestType_SLEEP:
				to = entity.AoiMoveType_SLEEP
			default:
				log.Panicf("unknown crowd runtime interest %v for person %d", rt.Interest, rt.Person.ID())
			}
			c.aoi.MoveBetweenSubmodules(rt.Person, entity.AoiMoveType_CROWD, to, -1)
			// 采用交换法删除一个元素
			c.persons[i] = c.persons[len(c.persons)-1]
			c.persons = c.persons[:len(c.persons)-1]
			rt.Person.UnsetCrowdByAoi()
		} else {
			// 没有删除，步进
			i++
		}
	}
}

func (c *Crowd) AddFromSleep(p entity.IPerson, to entity.AoiMoveType, laneID int32) {
	start := c.getSleepPointByRandom()
	switch to {
	case entity.AoiMoveType_GATE:
		c.add(p, start, c.getGatePointByLane(laneID), InterestType_EXIT)
	default:
		log.Panicf("unknown to %v from sleep for person %d", to, p.ID())
	}
}

func (c *Crowd) AddFromLaneToSleep(p entity.IPerson, laneID int32) {
	c.add(p, c.getGatePointByLane(laneID), c.getSleepPointByRandom(), InterestType_SLEEP)
}

// 新增室内行人
func (c *Crowd) add(p entity.IPerson, start, destination geometry.Point, interest InterestType) {
	p.SetRuntimeStatusByAoi(entity.PersonStatus_CROWD)
	r := &Runtime{
		Status:       Status_IDLE,
		Position:     start,
		Velocity:     geometry.Point{},
		Destination:  destination,
		DesiresSpeed: DEFAULT_DESIRED_SPEED,
		Interest:     interest,
		Person:       p,
	}
	p.SetCrowdByAoi()
	c.personInsertMutex.Lock()
	defer c.personInsertMutex.Unlock()
	c.personInserted = append(c.personInserted, r)
}

func (c *Crowd) getSleepPointByRandom() geometry.Point {
	// 随机从sleep points里选一个当起点
	return c.sleepPoints[c.generator.IntnSafe(len(c.sleepPoints))]
}

func (c *Crowd) getGatePointByLane(laneID int32) geometry.Point {
	// 出去的门
	p, ok := c.gates[laneID]
	if !ok {
		log.Panic("Aoi Crowd: add from wrong lane")
	}
	return p
}

// 判断给定点是否在某个门口处
func (c *Crowd) atAnyGate(p geometry.Point) bool {
	for _, g := range c.gates {
		if geometry.Distance(p, g) < GATE_SIZE/2 {
			return true
		}
	}
	return false
}

// 获取Aoi范围内随机点
// 注意此函数未保证多线程安全
func (c *Crowd) getRandomPosition() geometry.Point {
	if len(c.boundaryPoints) == 1 {
		return c.boundaryPoints[0]
	}
	if len(c.boundaryPoints) == 2 {
		return geometry.Blend(c.boundaryPoints[0], c.boundaryPoints[1], c.generator.Float64())
	}
	p0, p := c.boundaryPoints[0], geometry.Point{}
	inAoi := false
	for cnt := TRY_GET_RANDOM_POINT; !inAoi && cnt >= 0; cnt-- {
		// 以各三角形面积为概率选择一个三角形，再在其中随机取点
		index := c.generator.DiscreteDistribution(c.triangleAreas) + 1
		a, b := c.boundaryPoints[index], c.boundaryPoints[index+1]
		x, y := c.generator.Float64(), c.generator.Float64()
		if x+y >= 1 {
			x, y = 1-x, 1-y
		}
		// 使取点离边界有 5% 以上的距离
		p = geometry.Blend(p0, a, x*.9+.05).Add(geometry.Blend(p0, b, y*.9+.05)).Sub(p0)
		inAoi = p.InPolygon(c.boundaryPoints)
	}
	if !inAoi {
		log.Warnf("generate a point on the boundary of aoi %d. Consider choose a larger kTryGetRandomPoint(%d now)",
			c.aoi.ID(), TRY_GET_RANDOM_POINT)
		p = p0
	}
	return p
}

func (c *Crowd) HeadCount() int32 {
	if c == nil {
		return 0
	}
	return int32(len(c.persons))
}
