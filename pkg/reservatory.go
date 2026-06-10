package fluid

import (
	"math"
)

// Reservatório
type Reservatory struct {
	capacity    FluidUnit
	stored      FluidUnit
	lossPerTick FluidUnit
}

func NewReservatory(cap FluidUnit) *Reservatory {
	if cap < 0 {
		cap = 0
	}

	return &Reservatory{
		capacity: cap,
		stored:   0,
	}
}

func (r *Reservatory) Produce(amount FluidUnit) (FluidUnit, error) {
	if amount <= 0 {
		return 0, nil
	}

	space := r.capacity - r.stored
	if space <= 0 {
		return 0, nil
	}

	add := math.Min(space, amount)
	r.stored += add

	return add, nil
}

func (r *Reservatory) Consume(amount FluidUnit) (FluidUnit, error) {
	if amount <= 0 {
		return 0, nil
	}

	avail := r.stored
	if avail <= 0 {
		return 0, nil
	}

	use := math.Min(avail, amount)
	r.stored -= use

	return use, nil
}

func (r *Reservatory) Demand() FluidUnit { return r.lossPerTick }

func (r *Reservatory) Peek() FluidUnit { return r.stored }

func (r *Reservatory) Capacity() FluidUnit { return r.capacity }

func (r *Reservatory) Stored() FluidUnit { return r.stored }

// Apply tick loss aplica perda por tick; chamada externamente pelo simulador/tick loop.
func (r *Reservatory) Update() {
	if r.lossPerTick <= 0 || r.stored <= 0 {
		return
	}

	l := math.Min(r.lossPerTick, r.stored)
	r.stored -= l
}
