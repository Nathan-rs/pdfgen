package cli

import "fmt"

// Help exibe a ajuda geral do aplicativo.
func Help() {

	fmt.Println("PDFGEN")
	fmt.Println()

	fmt.Println("Uso:")
	fmt.Println()

	fmt.Println("  pdfgen create -file <arquivo.csv> -o <saida.pdf> [-c <timbre.png>]")
	fmt.Println()

	fmt.Println("Comandos:")
	fmt.Println()

	fmt.Println("  create     Gera um PDF a partir de um arquivo CSV")
	fmt.Println("  help       Exibe esta ajuda")
	fmt.Println()

	fmt.Println("Exemplos:")
	fmt.Println()

	fmt.Println("  pdfgen create -file logs.csv -o relatorio.pdf")
	fmt.Println("  pdfgen create -file logs.csv -o relatorio.pdf -c assets/timbre.png")
	fmt.Println()
}
