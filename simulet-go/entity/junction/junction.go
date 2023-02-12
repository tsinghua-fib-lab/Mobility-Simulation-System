package junction

import (
	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/entity/junction/trafficlight"
	"git.fiblab.net/sim/simulet-go/entity/lane"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	lightv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/traffic_light/v2"
)

type Junction struct {
	id           int32
	laneIDs      []int32
	trafficLight ITrafficLight        // 信号灯模块
	lanes        map[int32]*lane.Lane // 车道id->车道指 针映射表
}

func NewJunction(base *mapv2.Junction) *Junction {
	return &Junction{
		id:      base.Id,
		laneIDs: base.LaneIds,
		lanes:   make(map[int32]*lane.Lane),
	}
}

func (j *Junction) InitLanes(laneManager *lane.LaneManager) {
	lanes := make([]entity.ILaneTrafficLightSetter, 0)
	for _, laneID := range j.laneIDs {
		if lane, err := laneManager.Get(laneID); err != nil {
			log.Panic(err)
		} else {
			lane.SetParentJunctionWhenInit(j)
			j.lanes[laneID] = lane
			lanes = append(lanes, lane)
		}
	}
	j.trafficLight = trafficlight.NewLocalTrafficLight(j.id, lanes)
}

func (j *Junction) Prepare() {
	j.trafficLight.Prepare()
}

func (j *Junction) Update(stepInterval float64) {
	j.trafficLight.Update(stepInterval)
}

func (j *Junction) ID() int32 {
	return j.id
}

func (j *Junction) Lanes() map[int32]*lane.Lane {
	return j.lanes
}

func (j *Junction) SetTrafficLight(tl *lightv2.TrafficLight) error {
	return j.trafficLight.Set(tl)
}
