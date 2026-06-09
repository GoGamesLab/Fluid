package fluid

import (
	"math"
)

type FluidUnit = float64

// Fonte/sumidouro genérico de energia
type FluidSource interface {
	Produce(amount FluidUnit) (produced FluidUnit, err error) // tenta produzir até amount
	Peek() FluidUnit                                          // quantidade disponível (instantânea)
}

// Consome energia
type FluidSink interface {
	Consume(amount FluidUnit) (consumed FluidUnit, err error) // tenta consumir até amount
	Demand() FluidUnit                                        // demanda desejada (p.ex. por tick)
}

// Fonte e sumidouro combinados (baterias, nodes)
type FluidNode interface {
	FluidSource
	FluidSink
	Capacity() FluidUnit
	Stored() FluidUnit
}

// Hint para processos que requerem energia
type FluidRequirement struct {
	Amount   FluidUnit // em unidades base
	TypeHint string    // "heat", "electric", "kinetic", "nuclear" -- usado para escolher conversor
	MinTemp  float64   // opcional (para heat)
	MaxTemp  float64
}

// Converters registry and nodes
type FluidManager struct {
	nodes map[string]FluidNode
}

func NewFluidManager() *FluidManager {
	return &FluidManager{
		nodes: make(map[string]FluidNode),
	}
}

func (m *FluidManager) RegisterNode(id string, n FluidNode) { m.nodes[id] = n }

// Satisfaz um requisito: tenta prover Amount do TypeHint
// Retorna quantidade efetivamente fornecida (base units) e byproducts gerados
func (m *FluidManager) SatisfyRequirement(req FluidRequirement, preferNodeIDs []string) (provided FluidUnit, err error) {
	need := req.Amount // in base units
	if need <= 0 {
		return 0, nil
	}

	// build preferred set for O(1) checks
	preferred := make(map[string]struct{}, len(preferNodeIDs))
	for _, p := range preferNodeIDs {
		preferred[p] = struct{}{}
	}

	// strategy: iterate preferred nodes first, try to consume stored base units
	var totalProvided FluidUnit
	for _, id := range preferNodeIDs {
		node, exists := m.nodes[id]
		if !exists {
			continue
		}

		avail := node.Peek()
		if avail <= 0 {
			continue
		}

		want := math.Min(need-totalProvided, avail)
		consumed, err := node.Consume(want)
		if err != nil {
			continue
		}
		totalProvided += consumed

		if totalProvided >= need {
			break
		}
	}

	// If still lacking, attempt to produce from nodes that can Produce (reactors),
	// skipping the preferred nodes already tried.
	if totalProvided < need {
		for id, node := range m.nodes {
			if _, isPref := preferred[id]; isPref {
				continue
			}

			// try produce
			toProduce := need - totalProvided
			prod, pErr := node.Produce(toProduce)
			if pErr != nil || prod <= 0 {
				continue
			}
			totalProvided += prod
			if totalProvided >= need {
				break
			}
		}
	}

	return totalProvided, nil
}
