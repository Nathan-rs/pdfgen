package pdf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// progressReporter imprime o andamento da geração no stderr, sobrescrevendo
// a mesma linha (via \r) para não poluir o terminal com milhares de linhas.
//
// Se total <= 0 (contagem prévia não disponível/desativada), mostra só o
// total processado e a taxa, sem percentual/ETA.
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

	if p.total <= 0 {
		fmt.Fprintf(os.Stderr, "\r  %s registros processados — %s reg/s   ",
			formatNumber(processed), formatNumber(int64(rate)))
		return
	}

	pct := float64(processed) / float64(p.total) * 100
	if pct > 100 {
		pct = 100
	}

	// Estimativa linear: assume que o ritmo médio (rate) vai se manter
	// até o fim. Cobre só a fase de renderização dos lotes — o merge
	// final dos PDFs parciais roda depois e não entra nesse cálculo.
	eta := "?"
	switch {
	case processed >= p.total:
		eta = "0s"
	case rate > 0:
		remaining := float64(p.total-processed) / rate
		eta = time.Duration(remaining * float64(time.Second)).Round(time.Second).String()
	}

	fmt.Fprintf(os.Stderr, "\r  %s/%s registros (%.1f%%) — %s reg/s — ETA %s   ",
		formatNumber(processed), formatNumber(p.total), pct,
		formatNumber(int64(rate)), eta)
}

// formatNumber formata um inteiro com separador de milhar (padrão
// pt-BR, ponto): 1234567 -> "1.234.567". Só para exibição — não afeta
// nenhum cálculo interno.
func formatNumber(n int64) string {
	s := strconv.FormatInt(n, 10)

	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}

	var out []byte
	for i, d := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, d)
	}

	if neg {
		return "-" + string(out)
	}
	return string(out)
}
