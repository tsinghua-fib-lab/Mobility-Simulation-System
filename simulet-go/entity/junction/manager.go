package junction

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

type JunctionManager struct {
	data      map[int32]*Junction
	junctions []*Junction

	lanesInJunction []*lane.Lane
}

func NewManager() *JunctionManager {
	return &JunctionManager{
		data:            make(map[int32]*Junction),
		junctions:       make([]*Junction, 0),
		lanesInJunction: make([]*lane.Lane, 0),
	}
}

func (m *JunctionManager) Init(junctions []*mapv2.Junction) {
	for _, junction := range junctions {
		m.data[junction.Id] = NewJunction(junction)
	}
	m.junctions = lo.Values(m.data)
}

func (m *JunctionManager) InitLanes(laneManager *lane.LaneManager) {
	for _, junction := range m.data {
		junction.InitLanes(laneManager)
		m.lanesInJunction = append(m.lanesInJunction, lo.Values(junction.lanes)...)
	}
}

func (m *JunctionManager) Get(id int32) (*Junction, error) {
	if junction, ok := m.data[id]; !ok {
		return nil, fmt.Errorf("no id %d in junction data", id)
	} else {
		return junction, nil
	}
}

func (m *JunctionManager) Prepare() {
	parallel.GoFor(m.junctions, BLOCK_SIZE, func(j *Junction) { j.Prepare() })
}

func (m *JunctionManager) Update(stepInterval float64) {
	parallel.GoFor(m.junctions, BLOCK_SIZE, func(j *Junction) { j.Update(stepInterval) })
}
