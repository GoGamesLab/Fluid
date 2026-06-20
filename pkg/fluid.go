package fluid

type FluidUnit = float64

type FluidSource interface {
	Produce(amount FluidUnit) FluidUnit
	Peek() FluidUnit
}

type FluidSink interface {
	Consume()
	Demand() FluidUnit
}

type FluidManager struct {
	FluidSources map[string]FluidSource
	FluidSinks   map[string]FluidSink
}

func NewFluidManager() *FluidManager {
	return &FluidManager{
		FluidSources: make(map[string]FluidSource),
		FluidSinks:   make(map[string]FluidSink),
	}
}

func (m *FluidManager) RegisterSource(id string, n FluidSource) { m.FluidSources[id] = n }
func (m *FluidManager) RegisterSink(id string, n FluidSink)     { m.FluidSinks[id] = n }

func (m *FluidManager) Update() {
	for _, sink := range m.FluidSinks {
		need := sink.Demand()
		if need == 0 {
			continue
		}

		got := FluidUnit(0)

		for _, src := range m.FluidSources {
			if got >= need {
				break
			}

			remaining := need - got
			produced := src.Produce(remaining)
			got += produced
		}

		if got+1e-9 >= need {
			sink.Consume()
		}
	}
}
