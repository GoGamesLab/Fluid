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

func (b *Reservatory) Produce(amount FluidUnit) (FluidUnit, error) {
	if amount <= 0 {
		return 0, nil
	}

	space := b.capacity - b.stored
	if space <= 0 {
		return 0, nil
	}

	add := math.Min(space, amount)
	b.stored += add

	return add, nil
}

func (b *Reservatory) Consume(amount FluidUnit) (FluidUnit, error) {
	if amount <= 0 {
		return 0, nil
	}

	avail := b.stored
	if avail <= 0 {
		return 0, nil
	}

	use := math.Min(avail, amount)
	b.stored -= use

	return use, nil
}

func (b *Reservatory) Demand() FluidUnit { return b.lossPerTick }

func (b *Reservatory) Peek() FluidUnit { return b.stored }

func (b *Reservatory) Capacity() FluidUnit { return b.capacity }

func (b *Reservatory) Stored() FluidUnit { return b.stored }

// Apply tick loss aplica perda por tick; chamada externamente pelo simulador/tick loop.
func (b *Reservatory) Update() {
	if b.lossPerTick <= 0 || b.stored <= 0 {
		return
	}

	l := math.Min(b.lossPerTick, b.stored)
	b.stored -= l
}
