package pdf

import (
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergeFiles(inputs []string, output string) error {
	if len(inputs) < 2 {
		return fmt.Errorf("informe pelo menos 2 arquivos PDF para unir (recebido: %d)", len(inputs))
	}

	for _, in := range inputs {
		if _, err := os.Stat(in); err != nil {
			return fmt.Errorf("arquivo %q: %w", in, err)
		}
	}

	// ATENÇÃO: verificar a assinatura de api.MergeCreateFile, pois pode varia entre versões
	// doc: `go doc github.com/pdfcpu/pdfcpu/pkg/api MergeCreateFile`
	// caso o build quebre após atualizar a dependência.
	if err := api.MergeCreateFile(inputs, output, false, nil); err != nil {
		return fmt.Errorf("unindo PDFs: %w", err)
	}

	return nil
}
