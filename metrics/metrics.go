package pool

import "time"

// Metrics almacena las métricas del un elemento.
type Metrics struct {
	// Inicio
	ini time.Time
	// Transacciones por segundo
	tps float32

	// Tranacciones por minuto
	tpm float32
}

func (m *Metrics) timeup() time.Duration {
	return time.Since(m.ini)
}

func (m *Metrics) start() {
	m.ini = time.Now()
}
