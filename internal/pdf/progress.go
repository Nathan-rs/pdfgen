package pdf

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// progressReporter imprime o andamento da geração no stderr, sobrescrevendo
// a mesma linha (via \r) para não poluir o terminal com milhares de linhas.
//
// Se total <= 0 (contagem prévia não disponível/desativada), mostra só o
// total processado e a taxa, sem percentual.
type progressReporter struct {
	total     int64
	processed atomic.Int64
	startedAt time.Time
	stopCh    chan struct{}
	doneCh    chan struct{}
}

func newProgressReporter(total int) *progressReporter {
	return &progressReporter{
		total:     int64(total),
		startedAt: time.Now(),
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
	}
}

// Add soma n ao total de registros processados. Seguro para chamar de
// múltiplas goroutines simultaneamente.
func (p *progressReporter) Add(n int) {
	p.processed.Add(int64(n))
}

// Run imprime o progresso a cada `interval` até Stop ser chamado.
// Deve rodar em goroutine própria: `go reporter.Run(200 * time.Millisecond)`.
func (p *progressReporter) Run(interval time.Duration) {
	defer close(p.doneCh)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.print()
		case <-p.stopCh:
			p.print()
			fmt.Fprintln(os.Stderr)
			return
		}
	}
}

// Stop sinaliza o fim e bloqueia até a última linha ser impressa.
func (p *progressReporter) Stop() {
	close(p.stopCh)
	<-p.doneCh
}

func (p *progressReporter) print() {
	processed := p.processed.Load()
	elapsed := time.Since(p.startedAt).Seconds()

	rate := 0.0
	if elapsed > 0 {
		rate = float64(processed) / elapsed
	}

	if p.total > 0 {
		pct := float64(processed) / float64(p.total) * 100
		if pct > 100 {
			pct = 100
		}

		eta := "?"
		if rate > 0 && processed < p.total {
			remaining := float64(p.total-processed) / rate
			eta = time.Duration(remaining * float64(time.Second)).Round(time.Second).String()
		} else if processed >= p.total {
			eta = "0s"
		}

		fmt.Fprintf(os.Stderr, "\r  %d/%d registros (%.1f%%) — %.0f reg/s — ETA %s   ",
			processed, p.total, pct, rate, eta)
		return
	}

	fmt.Fprintf(os.Stderr, "\r  %d registros processados — %.0f reg/s   ", processed, rate)
}
