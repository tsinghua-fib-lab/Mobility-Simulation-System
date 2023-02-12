package sleep

import (
	"container/heap"
	"sync"

	"git.fiblab.net/sim/simulet-go/entity"
	tripv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/trip/v2"
	"git.fiblab.net/sim/simulet-go/utils/container"
)

type Sleep struct {
	aoi               entity.IAoi      // 对应aoi
	personInserted    []entity.IPerson // 新加入的人
	personInsertMutex sync.Mutex
	sleepingPeople    container.PriorityQueue[entity.IPerson] // 未到活动时间的人，无计算任务
	awakePeople       []entity.IPerson                        // 被唤醒的人
}

func New(aoi entity.IAoi) *Sleep {
	return &Sleep{
		aoi:               aoi,
		personInserted:    make([]entity.IPerson, 0),
		personInsertMutex: sync.Mutex{},
		sleepingPeople:    make(container.PriorityQueue[entity.IPerson], 0),
		awakePeople:       make([]entity.IPerson, 0),
	}
}

func (s *Sleep) Prepare() {
	pos := s.aoi.Centroid()
	for _, p := range s.personInserted {
		p.ResetScheduleIfNeed()
		p.SetSnapshotByAoi(entity.BaseRuntime{
			Position: pos,
		}, entity.BaseRuntimeInAoi{
			Aoi: s.aoi,
		})
		s.sleepingPeople.Push(&container.Item[entity.IPerson]{
			Value:    p,
			Priority: p.Schedule().GetDepartureTime(),
		})
	}
	heap.Init(&s.sleepingPeople)
	s.personInserted = []entity.IPerson{}
}

func (s *Sleep) Update() {
	// 检查awake person导航请求是否成功
	for _, p := range s.awakePeople {
		if !p.RouteSuccessful() {
			// 导航不成功，重新插入以再次请求
			s.personInsertMutex.Lock()
			s.personInserted = append(s.personInserted, p)
			s.personInsertMutex.Unlock()
			continue
		}
		// 出门，修改事件为当前t对应值
		t := p.Schedule().GetTrip()
		p.SetRuntimeActivityByAoi(t.GetActivity())
		// 检查人的出行方式
		if t.Mode == tripv2.TripMode_TRIP_MODE_DRIVE_ONLY {
			s.aoi.MoveBetweenSubmodules(
				p,
				entity.AoiMoveType_SLEEP,
				entity.AoiMoveType_GATE,
				p.CurrentVehicleLaneID(),
			)
		} else {
			s.aoi.MoveBetweenSubmodules(
				p,
				entity.AoiMoveType_SLEEP,
				entity.AoiMoveType_GATE,
				p.CurrentPedestrianLaneID(),
			)
		}
	}
	s.awakePeople = []entity.IPerson{}
	// 检查sleep person是否到达出发时间
	if len(s.sleepingPeople) > 0 {
		for p := s.sleepingPeople[0].Value; p.CheckDeparture(); p = s.sleepingPeople[0].Value {
			// 发出导航请求
			p.RequestRouteFromAoi(s.aoi)
			// 放入醒来的人
			s.awakePeople = append(s.awakePeople, p)
			heap.Pop(&s.sleepingPeople)
			if len(s.sleepingPeople) == 0 {
				// 没人了
				break
			}
		}
	}
}

func (s *Sleep) Add(p entity.IPerson) {
	p.SetRuntimeStatusByAoi(entity.PersonStatus_SLEEP)
	s.personInsertMutex.Lock()
	defer s.personInsertMutex.Unlock()
	s.personInserted = append(s.personInserted, p)
}

func (s *Sleep) HeadCount() int32 {
	return int32(len(s.sleepingPeople) + len(s.awakePeople))
}
