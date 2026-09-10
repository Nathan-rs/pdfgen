package cli

import (
	"flag"
	"fmt"
	"runtime"

	"pdfgen/internal/pdf"
	"pdfgen/internal/validate"
)

// Create executa o comando:
//
//	pdfgen create -file dados.csv -o relatorio.pdf [-c timbre.png] [-workers N] [-progress=false]
func Create(args []string) error {
	var (
		csvFile      string
		output       string
		cover        string
		workers      int
		showProgress bool
	)

	fs := flag.NewFlagSet("create", flag.ContinueOnError)

	fs.StringVar(&csvFile, "file", "", "Arquivo CSV de entrada")
	fs.StringVar(&output, "o", "", "Arquivo PDF de saída")
	fs.StringVar(&cover, "c", "", "Imagem PNG utilizada como timbre (opcional)")
	fs.IntVar(&workers, "workers", runtime.NumCPU(), "Número de goroutines para geração paralela")
	fs.BoolVar(&showProgress, "progress", true, "Mostrar progresso durante a geração")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validações
	if err := validate.CSV(csvFile); err != nil {
		return err
	}

	if err := validate.Output(output); err != nil {
		return err
	}

	if err := validate.Cover(cover); err != nil {
		return err
	}

	if err := pdf.GenerateConcurrent(csvFile, output, cover, workers, showProgress); err != nil {
		return fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return nil
}
