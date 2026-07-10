package validate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CSV valida o arquivo CSV informado.
func CSV(path string) error {

	if strings.TrimSpace(path) == "" {
		return errors.New("o parâmetro -file é obrigatório")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("arquivo CSV não encontrado: %s", path)
	}

	if info.IsDir() {
		return fmt.Errorf("%s é um diretório", path)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".csv" {
		return fmt.Errorf("arquivo deve possuir extensão .csv")
	}

	return nil
}

// Output valida o caminho do PDF de saída.
func Output(path string) error {

	if strings.TrimSpace(path) == "" {
		return errors.New("o parâmetro -o é obrigatório")
	}

	dir := filepath.Dir(path)

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("diretório de saída não existe: %s", dir)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s não é um diretório", dir)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".pdf" {
		return fmt.Errorf("o arquivo de saída deve possuir extensão .pdf")
	}

	return nil
}

// Cover valida a imagem utilizada como timbre.
func Cover(path string) error {

	if strings.TrimSpace(path) == "" {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("arquivo de timbre não encontrado: %s", path)
	}

	if info.IsDir() {
		return fmt.Errorf("%s é um diretório", path)
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".png", ".jpg", ".jpeg":
		return nil
	default:
		return fmt.Errorf("imagem de timbre deve ser PNG, JPG ou JPEG")
	}
}
