package vehicle

import (
	"git.fiblab.net/sim/simulet-go/entity"
	"git.fiblab.net/sim/simulet-go/utils/container"
)

func newNode(key float64, value entity.IVehicle) *container.ListNode[entity.IVehicle, entity.VehicleSideLink] {
	return &container.ListNode[entity.IVehicle, entity.VehicleSideLink]{
		Key:   key,
		Value: value,
	}
}
