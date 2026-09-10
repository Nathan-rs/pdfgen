package pdf

// Geração concorrente: uma goroutine lê o CSV sequencialmente (I/O é
// barato, não é o gargalo) e distribui lotes de linhas para N workers.
// Cada worker tem seu próprio Document/gofpdf.Fpdf — que é o que permite
// paralelizar, já que gofpdf.Fpdf não é thread-safe e não pode ser
// escrito por mais de uma goroutine ao mesmo tempo.
//
// Cada worker grava seu lote como um PDF parcial em disco; no final os
// parciais são unidos, na ordem original, via pdfcpu.

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// defaultBatchSize: lotes pequenos demais aumentam o overhead do merge
// final (mais arquivos parciais); lotes grandes demais reduzem o
// paralelismo (menos lotes que workers disponíveis). Vale calibrar via
// benchmark com o tamanho real dos seus CSVs.
const defaultBatchSize = 3000

type batch struct {
	index int
	rows  [][]string
}

// GenerateConcurrent lê csvFile, renderiza as linhas em paralelo usando
// `workers` goroutines e grava o resultado final em output.
// Se showProgress for true, imprime o andamento no stderr.
func GenerateConcurrent(csvFile, output, cover string, workers int, showProgress bool) error {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("abrindo csv: %w", err)
	}
	defer f.Close()

	// BOMOverride detecta e descarta um BOM (UTF-8/UTF-16LE/UTF-16BE) no
	// início do arquivo, comum em CSVs exportados pelo Excel/Windows.
	decoder := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	reader := csv.NewReader(transform.NewReader(f, decoder))
	reader.Comma = ';'

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("lendo cabeçalho: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "pdfgen-parts-*")
	if err != nil {
		return fmt.Errorf("criando diretório temporário: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Contagem prévia (best-effort) só para exibir % e ETA na barra de
	// progresso. É uma passagem rápida contando linhas físicas do
	// arquivo — não faz parsing de CSV, então em campos com quebra de
	// linha dentro de aspas o número pode ficar levemente impreciso.
	// Isso NÃO afeta a geração em si, só a exibição do progresso.
	var reporter *progressReporter
	if showProgress {
		total, err := countRecords(csvFile)
		if err != nil {
			total = 0 // segue sem percentual/ETA, só mostra o contador
		}
		reporter = newProgressReporter(total)
		go reporter.Run(200 * time.Millisecond)
	}

	batches := make(chan batch, workers*2)
	errCh := make(chan error, workers+1)

	var (
		wg      sync.WaitGroup
		partsMu sync.Mutex
		parts   []string
	)

	// --- workers: cada um com seu próprio Document isolado ---
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for b := range batches {
				partPath := filepath.Join(tmpDir, fmt.Sprintf("part-%06d.pdf", b.index))

				doc := NewBatch(headers, cover, b.index*defaultBatchSize)
				doc.AddRecords(b.rows)
				if err := doc.Close(partPath); err != nil {
					select {
					case errCh <- fmt.Errorf("worker %d, lote %d: %w", workerID, b.index, err):
					default:
					}
					return
				}

				if reporter != nil {
					reporter.Add(len(b.rows))
				}

				partsMu.Lock()
				parts = append(parts, partPath)
				partsMu.Unlock()
			}
		}(w)
	}

	// --- produtor: lê o CSV sequencialmente e distribui lotes ---
	go func() {
		defer close(batches)
		idx := 0
		rows := make([][]string, 0, defaultBatchSize)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				select {
				case errCh <- fmt.Errorf("lendo linha do csv: %w", err):
				default:
				}
				return
			}
			rows = append(rows, record)
			if len(rows) >= defaultBatchSize {
				batches <- batch{index: idx, rows: rows}
				idx++
				rows = make([][]string, 0, defaultBatchSize)
			}
		}
		if len(rows) > 0 {
			batches <- batch{index: idx, rows: rows}
		}
	}()

	wg.Wait()
	close(errCh)

	if reporter != nil {
		reporter.Stop()
	}

	if err, ok := <-errCh; ok && err != nil {
		return err
	}

	if len(parts) == 0 {
		return fmt.Errorf("csv sem registros para gerar PDF")
	}

	// nomes com zero-padding (%06d) preservam a ordem original das linhas
	sort.Strings(parts)

	if showProgress {
		fmt.Fprintf(os.Stderr, "Unindo %d PDFs parciais...\n", len(parts))
	}

	// ATENÇÃO: a assinatura de api.MergeCreateFile varia entre versões
	// do pdfcpu (parâmetro de configuração mudou de tipo em versões
	// recentes). Confira com `go doc github.com/pdfcpu/pdfcpu/pkg/api MergeCreateFile`.
	if err := api.MergeCreateFile(parts, output, false, nil); err != nil {
		return fmt.Errorf("unindo PDFs parciais: %w", err)
	}

	if showProgress {
		fmt.Fprintf(os.Stderr, "PDF gerado: %s\n", output)
	}

	return nil
}

// countRecords faz uma contagem rápida de linhas físicas do arquivo,
// só para alimentar a barra de progresso (não é usada na geração real,
// que usa encoding/csv corretamente). Desconta 1 linha do cabeçalho.
func countRecords(csvFile string) (int, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024) // tolera linhas longas

	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if count > 0 {
		count-- // desconta o cabeçalho
	}
	return count, nil
}
