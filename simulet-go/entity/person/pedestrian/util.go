package pedestrian

import (
	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/utils/container"
)

func newNode(key float64, value entity.IPedestrian) *container.ListNode[entity.IPedestrian, struct{}] {
	return &container.ListNode[entity.IPedestrian, struct{}]{
		Key:   key,
		Value: value,
	}
}
