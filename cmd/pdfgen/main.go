package main

import (
	"fmt"
	"os"
	"pdfgen/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.Help()
		return
	}

	switch os.Args[1] {

	case "create":

		if err := cli.Create(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
			os.Exit(1)
		}

	case "-h", "--help", "help":

		cli.Help()

	default:

		fmt.Fprintf(os.Stderr, "Comando '%s' não reconhecido.\n\n", os.Args[1])
		cli.Help()
		os.Exit(1)
	}
}
