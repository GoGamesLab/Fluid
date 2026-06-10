package fluid

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

// Fonte e sumidouro combinados (reservatórios, nodes)
type FluidNode interface {
	FluidSource
	FluidSink
	Capacity() FluidUnit
	Stored() FluidUnit
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
