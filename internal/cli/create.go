package cli

import (
	"flag"
	"fmt"
	"pdfgen/internal/pdf"
	"pdfgen/internal/validate"
)

// Create executa o comando:
//
//	pdfgen create -file dados.csv -o relatorio.pdf [-c timbre.png]
func Create(args []string) error {
	var (
		csvFile string
		output  string
		cover   string
	)

	fs := flag.NewFlagSet("create", flag.ContinueOnError)

	fs.StringVar(&csvFile, "file", "", "Arquivo CSV de entrada")
	fs.StringVar(&output, "o", "", "Arquivo PDF de saída")
	fs.StringVar(&cover, "c", "", "Imagem PNG utilizada como timbre (opcional)")

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

	document, err := pdf.New(csvFile, cover)
	if err != nil {
		return fmt.Errorf("erro ao criar documento: %w", err)
	}

	if err := document.Save(output); err != nil {
		return fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return nil
}
