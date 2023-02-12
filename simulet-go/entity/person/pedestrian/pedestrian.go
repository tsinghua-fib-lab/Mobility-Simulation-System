package pedestrian

import (
	"math"

	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/person/route"
	"git.fiblab.net/sim/simulet-go/utils/container"
	"git.fiblab.net/sim/simulet-go/utils/geometry"
	"github.com/samber/lo"
)

const (
	BIKE_SPEED                    = 4
	DEFAULT_DESIRED_SPEED_ON_LANE = 1.34
	MAX_NOISE_ON_PEDESTRIAN_SPEED = .5
)

type Runtime struct {
	entity.BaseRuntime
	entity.BaseRuntimeOnRoad
	IsForward bool // 是否正向行走
}
type Pedestrian struct {
	person entity.IPerson // 对应人的指针

	snapshot Runtime // 快照
	runtime  Runtime // 运行时数据

	route *route.PedestrianRoute // 路径规划

	isEnd        bool    // 是否到达终点
	walkingSpeed float64 // 行走速度

	// Lane链表
	node *container.ListNode[entity.IPedestrian, struct{}]
}

func NewPedestrian(person entity.IPerson, route *route.PedestrianRoute) *Pedestrian {
	// 行走速度添加正态分布随机扰动
	speed := math.Max(.1,
		DEFAULT_DESIRED_SPEED_ON_LANE+lo.Clamp(
			.5*MAX_NOISE_ON_PEDESTRIAN_SPEED*person.Generator().NormFloat64(),
			-MAX_NOISE_ON_PEDESTRIAN_SPEED,
			MAX_NOISE_ON_PEDESTRIAN_SPEED))

	firstSeg := route.Current()
	lane := firstSeg.Lane
	s := route.Start.Aoi.WalkingS(lane.ID())
	p := &Pedestrian{
		person: person,
		runtime: Runtime{
			BaseRuntime: entity.BaseRuntime{
				Position: lane.GetPositionByS(s),
			},
			BaseRuntimeOnRoad: entity.BaseRuntimeOnRoad{
				Lane: lane,
				S:    s,
			},
			IsForward: firstSeg.IsForward(),
		},
		route:        route,
		walkingSpeed: speed,
	}
	p.node = newNode(p.runtime.S, p)
	lane.ReportPedestrianAdded(p.node)
	return p
}

func (p *Pedestrian) Prepare() {
	p.snapshot = p.runtime
}

func (p *Pedestrian) Update(stepInterval float64) {
	s := p.S()
	if s > p.route.Current().Lane.Length() || s < 0 {
		log.Panicf("Pedestrian: s %v out of lane range {%v,%v}",
			s, 0, p.route.Current().Lane.Length())
	}
	seg := p.route.Current()
	nextLane := false
	redLight := false
	var ds float64
	if p.person.IsOnBike() {
		ds = BIKE_SPEED * stepInterval
	} else {
		ds = p.walkingSpeed * stepInterval
	}
	if p.IsForward() {
		s += ds
		length := seg.Lane.Length()
		if s > length {
			if !p.route.AtLast() {
				redLight = true
				s = length
			} else {
				nextLane = true
				for {
					s -= length
					if ok := p.route.Step(); !ok {
						p.isEnd = true
						break
					}
					length = p.route.Current().Lane.Length()
					if s <= length {
						break
					}
				}
			}
		}
	} else {
		s -= ds
		if s <= 0 {
			if !p.route.AtLast() {
				redLight = true
				s = 0
			} else {
				nextLane = true
				length := 0.0
				for s+length < 0 {
					s += length
					if ok := p.route.Step(); !ok {
						p.isEnd = true
						break
					}
					length = p.route.Current().Lane.Length()
				}
				s = -s
			}
		}
	}
	endS := p.route.EndS
	if !p.isEnd {
		sFlag := false
		if seg.IsForward() {
			sFlag = s >= endS
		} else {
			sFlag = s <= endS
		}
		if p.route.AtLast() && sFlag {
			p.isEnd = true
		}
	}
	if p.isEnd {
		p.runtime.Lane = p.route.Last().Lane
		p.runtime.S = endS
		// 增量更新车道索引（不再维护数据）
		p.snapshot.Lane.ReportPedestrianRemoved(p.node)
		return
	}
	if !redLight {
		p.runtime.Speed = p.walkingSpeed
		if nextLane {
			seg = p.route.Current()
			p.runtime.IsForward = seg.IsForward()
			if !p.runtime.IsForward {
				s = seg.Lane.Length() - s
			}
			p.runtime.Lane = seg.Lane
		}
	} else {
		p.runtime.Speed = 0
	}
	p.runtime.S = s
	p.runtime.Position = p.runtime.Lane.GetPositionByS(p.runtime.S)
	p.runtime.Direction = p.runtime.Lane.GetDirectionByS(s)
	if !p.runtime.IsForward {
		p.runtime.Direction += math.Pi
	}
	// 车道链表更新
	p.node.Key = p.runtime.S
	// 增量更新车道索引（维护数据）
	if p.snapshot.Lane != p.runtime.Lane {
		p.snapshot.Lane.ReportPedestrianRemoved(p.node)
		// 换一个新的node来避免remove操作和add操作处理同一个对象需要保证先后顺序
		p.node = newNode(p.runtime.S, p)
		p.runtime.Lane.ReportPedestrianAdded(p.node)
	}
}

// getter

func (p *Pedestrian) FetchBaseSnapshotForPerson() (entity.BaseRuntime, entity.BaseRuntimeOnRoad) {
	return p.snapshot.BaseRuntime, p.snapshot.BaseRuntimeOnRoad
}

func (p *Pedestrian) Lane() entity.ILane {
	return p.snapshot.Lane
}

func (p *Pedestrian) S() float64 {
	return p.snapshot.S
}

func (p *Pedestrian) Speed() float64 {
	return p.snapshot.Speed
}

func (p *Pedestrian) Direction() float64 {
	return p.snapshot.Direction
}

func (p *Pedestrian) Position() geometry.Point {
	return p.snapshot.Position
}

func (p *Pedestrian) Snapshot() Runtime {
	return p.snapshot
}

func (p *Pedestrian) IsForward() bool {
	return p.snapshot.IsForward
}

func (p *Pedestrian) GetEndByPerson() (entity.IAoi, entity.ILane, bool) {
	if p.isEnd {
		return p.route.EndAoi, p.runtime.Lane, true
	} else {
		return nil, nil, false
	}
}
