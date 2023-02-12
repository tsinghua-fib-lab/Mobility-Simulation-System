package lane

import (
	"fmt"

	mapv2 "git.fiblab.net/sim/simulet-go/gen/proto/go/wolong/map/v2"
	"git.fiblab.net/sim/simulet-go/utils/parallel"
	"github.com/samber/lo"
)

const (
	BLOCK_SIZE = 2048
)

var (
	Manager = NewManager()
)

type LaneManager struct {
	data  map[int32]*Lane
	lanes []*Lane
}

func NewManager() *LaneManager {
	return &LaneManager{
		data:  make(map[int32]*Lane),
		lanes: make([]*Lane, 0),
	}
}

func (m *LaneManager) Init(lanes []*mapv2.Lane) {
	for _, lane := range lanes {
		m.data[lane.Id] = NewLane(lane)
	}
	m.lanes = lo.Values(m.data)
	parallel.GoFor(m.lanes, BLOCK_SIZE, func(l *Lane) { l.InitLanes(m) })
}

func (m *LaneManager) Get(id int32) (*Lane, error) {
	if lane, ok := m.data[id]; !ok {
		return nil, fmt.Errorf("no id %d in lane data", id)
	} else {
		return lane, nil
	}
}

func (m *LaneManager) Prepare() {
	parallel.GoFor(m.lanes, BLOCK_SIZE, func(l *Lane) { l.Prepare() })
	parallel.GoFor(m.lanes, BLOCK_SIZE, func(l *Lane) { l.Prepare2() })
}
