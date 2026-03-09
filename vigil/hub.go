package vigil

import (
	"context"
)

type Hub struct {
	Egg *ElectricityEgg
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.Egg.Merge(data)
	return nil
}
