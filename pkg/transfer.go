package fluid

import (
	"math"
)

type FluidNetwork struct {
	Amount   float64 // Quantidade atual na rede de tubos
	Pressure float32 // Determina a velocidade de escoamento
}

type FluidPipe struct {
	ThroughputPerTick   FluidUnit
	LossPerUnitDistance float64 // frac lost per unit distance
	Distance            float64
}

// Transfer tira fluido de uma FONTE (Source), aplica perdas e entrega a um CONSUMIDOR (Sink)
func (p *FluidPipe) Transfer(from FluidSource, to FluidSink, want FluidUnit) (transferred FluidUnit, err error) {
	if want <= 0 {
		return 0, nil
	}

	// 1. Puxa a fluido disponível da fonte
	got, err := from.(FluidSink).Consume(want) // Cast temporário se for um FluidNode, ou ajuste a interface
	if err != nil || got <= 0 {
		return 0, err
	}

	// 2. Aplica perda por distância (Resistência do cano)
	lossFrac := math.Max(0.0, p.LossPerUnitDistance*p.Distance)
	if lossFrac > 1.0 {
		lossFrac = 1.0
	}

	lossAmount := float64(got) * lossFrac
	effective := FluidUnit(float64(got) - lossAmount)

	// 3. Aplica o teto de vazão (Throughput) do cano
	if p.ThroughputPerTick > 0 && effective > p.ThroughputPerTick {
		// Se estourar o limite, o excesso é descartado ou acumula (aqui limitamos o que é entregue)
		// Em jogos como Factorio, o cano simplesmente gargala a transferência.
		effective = p.ThroughputPerTick
	}

	// 4. Entrega a fluido convertida para o consumidor
	// Como a interface FluidSink usa Consume(), para fins de rede pura,
	// assumimos que injetar fluido em um Sink/Node é feito via Produce se for Node,
	// ou se o 'to' for um acumulador.
	// Nota: Se 'to' for uma máquina consumidora pura, ela precisará de um buffer interno (Reservatório)
	produced, err := to.(FluidSource).Produce(effective)
	if err != nil {
		return 0, err
	}

	return produced, nil
}
