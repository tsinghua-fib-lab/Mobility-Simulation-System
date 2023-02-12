package aoi

import (
	"fmt"

	"git.fiblab.net/sim/simulet-go/entity/lane"
	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	"git.fiblab.net/sim/simulet-go/utils/parallel"
	"github.com/samber/lo"
)

const (
	BLOCK_SIZE = 512
)

var (
	Manager = NewManager()
)

type AoiManager struct {
	data map[int32]*Aoi
	aois []*Aoi
}

func NewManager() *AoiManager {
	m := &AoiManager{
		data: make(map[int32]*Aoi),
		aois: make([]*Aoi, 0),
	}
	return m
}

func (m *AoiManager) Init(aois []*mapv2.Aoi) {
	for _, aoi := range aois {
		m.data[aoi.Id] = NewAoi(aoi, m)
	}
	m.aois = lo.Values(m.data)
}

func (m *AoiManager) InitLanes(laneManager *lane.LaneManager) {
	parallel.GoFor(m.aois, BLOCK_SIZE, func(aoi *Aoi) { aoi.InitLanes(laneManager) })
}

func (m *AoiManager) Get(id int32) (*Aoi, error) {
	if aoi, ok := m.data[id]; !ok {
		return nil, fmt.Errorf("no id %d in aoi data", id)
	} else {
		return aoi, nil
	}
}

func (m *AoiManager) Prepare() {
	parallel.GoFor(m.aois, BLOCK_SIZE, func(a *Aoi) { a.Prepare() })
}

func (m *AoiManager) Update(stepInterval float64) {
	parallel.GoFor(m.aois, BLOCK_SIZE, func(a *Aoi) { a.Update(stepInterval) })
}
