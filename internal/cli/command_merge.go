package cli

import (
	"flag"
	"fmt"
	"pdfgen/internal/pdf"
	"strings"
)

func Merge(args []string) error {

	var (
		output    string
		filesFlag string
	)

	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.StringVar(&output, "o", "", "Arquivo de saída")
	fs.StringVar(&filesFlag, "files", "", "Lista de PDFs a unir, separados por vírgula(alternativa aos argumentos posicionais)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var files []string
	if filesFlag != "" {
		for _, f := range strings.Split(filesFlag, ",") {
			if f = strings.TrimSpace(f); f != "" {
				files = append(files, f)
			}
		}
	}

	// Argumentos posicionais (o que sobrar depois das flags) também
	// contam — permite `pdfgen merge -o final.pdf *.pdf` direto no shell.
	files = append(files, fs.Args()...)

	if output == "" {
		return fmt.Errorf("informe o arquivo de saída com -o")
	}

	if err := pdf.MergeFiles(files, output); err != nil {
		return fmt.Errorf("erro ao unir PDFs: %w", err)
	}

	fmt.Printf("PDF gerado: %s\n", output)

	return nil
}
