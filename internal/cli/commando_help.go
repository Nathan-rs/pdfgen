package cli

import "fmt"

// Help exibe a ajuda geral do aplicativo.
func Help() {

	fmt.Println("PDFGEN")
	fmt.Println()

	fmt.Println("Uso:")
	fmt.Println()

	fmt.Println("  pdfgen create -file <arquivo.csv> -o <saida.pdf> [-c <timbre.png>] [-workers <n>] [-progress=false]")
	fmt.Println("  pdfgen merge  -o <saida.pdf> <arquivo1.pdf> <arquivo2.pdf> [...]")
	fmt.Println()

	fmt.Println("Comandos:")
	fmt.Println()

	fmt.Println("  create     Gera um PDF a partir de um arquivo CSV")
	fmt.Println("  merge      Une vários PDFs em um único arquivo, na ordem informada")
	fmt.Println("  help       Exibe esta ajuda")
	fmt.Println()

	fmt.Println("Opções do create:")
	fmt.Println()

	fmt.Println("  -file      Arquivo CSV de entrada")
	fmt.Println("  -o         Arquivo PDF de saída")
	fmt.Println("  -c         Imagem PNG utilizada como timbre (opcional)")
	fmt.Println("  -workers   Número de goroutines para geração paralela (padrão: núcleos da máquina)")
	fmt.Println("  -progress  Mostrar progresso durante a geração (padrão: true)")
	fmt.Println()

	fmt.Println("Opções do merge:")
	fmt.Println()

	fmt.Println("  -o         Arquivo PDF de saída (obrigatório)")
	fmt.Println("  -files     Lista de PDFs separados por vírgula (alternativa aos argumentos posicionais)")
	fmt.Println()
	fmt.Println("  Obs: as opções devem vir antes dos arquivos. São necessários ao menos 2 PDFs.")
	fmt.Println()

	fmt.Println("Exemplos:")
	fmt.Println()

	fmt.Println("  pdfgen create -file logs.csv -o relatorio.pdf")
	fmt.Println("  pdfgen create -file logs.csv -o relatorio.pdf -c assets/timbre.png")
	fmt.Println("  pdfgen create -file logs.csv -o relatorio.pdf -workers 4")
	fmt.Println("  pdfgen merge -o final.pdf jan.pdf fev.pdf mar.pdf")
	fmt.Println("  pdfgen merge -o final.pdf partes/*.pdf")
	fmt.Println("  pdfgen merge -o final.pdf -files jan.pdf,fev.pdf,mar.pdf")
	fmt.Println()
}
