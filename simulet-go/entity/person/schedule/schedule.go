package schedule

import (
	tripv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/trip/v2"
	"git.fiblab.net/sim/simulet-go/utils"
)

type Schedule struct {
	base            []*tripv2.Schedule
	scheduleIndex   int32   // 当前schedule下标
	tripIndex       int32   // 当前trip下标
	loopCount       int32   // schedule循环计数器
	lastTripEndTime float64 // 上次trip结束时间
}

func NewSchedule() *Schedule {
	return &Schedule{
		base: make([]*tripv2.Schedule, 0),
	}
}

func (s *Schedule) Base() []*tripv2.Schedule {
	return s.base
}

func (s *Schedule) NextTrip(time float64) bool {
	if len(s.base) == 0 {
		return false
	}
	schedule := s.base[s.scheduleIndex]
	s.lastTripEndTime = time
	if s.tripIndex++; s.tripIndex == int32(len(schedule.Trips)) {
		s.tripIndex = 0
		if s.loopCount++; schedule.LoopCount > 0 && s.loopCount >= schedule.LoopCount {
			s.loopCount = 0
			if s.scheduleIndex++; s.scheduleIndex == int32(len(s.base)) {
				s.base = make([]*tripv2.Schedule, 0)
				s.scheduleIndex = 0
				return false
			} else {
				if waitTime := s.base[s.scheduleIndex].WaitTime; waitTime != nil {
					s.lastTripEndTime += *waitTime
				} else if departureTime := s.base[s.scheduleIndex].DepartureTime; departureTime != nil {
					s.lastTripEndTime = *departureTime
				}
			}
		}
	}
	return true
}

func (s *Schedule) GetTrip() *tripv2.Trip {
	return s.base[s.scheduleIndex].Trips[s.tripIndex]
}

func (s *Schedule) Set(base []*tripv2.Schedule, time float64) {
	if len(base) == 0 {
		return
	}
	s.base = append(base, s.base...)
	s.scheduleIndex, s.tripIndex, s.loopCount = 0, 0, 0
	if lastDepartureTime := base[0].DepartureTime; lastDepartureTime != nil {
		s.lastTripEndTime = *lastDepartureTime
	} else if waitTime := base[0].WaitTime; waitTime != nil {
		s.lastTripEndTime = time + *waitTime
	} else {
		s.lastTripEndTime = time
	}
}

func (s *Schedule) Empty() bool {
	return len(s.base) == 0
}

func (s *Schedule) GetDepartureTime() float64 {
	if len(s.base) == 0 {
		//没有日程则返回∞
		return utils.INF
	}
	trip := s.GetTrip()
	if departureTime := trip.DepartureTime; departureTime != nil {
		if s.loopCount != 0 {
			log.Warning("departure time used in loop")
		}
		return *departureTime
	}
	if waitTime := trip.WaitTime; waitTime != nil {
		return s.lastTripEndTime + *waitTime
	} else {
		return s.lastTripEndTime
	}
}
