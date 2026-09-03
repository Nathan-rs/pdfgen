# PDFGEN

Aplicação CLI em Go para gerar documentos PDF tabulares a partir de arquivos CSV, com suporte a papel timbrado institucional.

---

## Funcionalidades

* Leitura de arquivos CSV (separador `;`) via streaming — sem carregar o arquivo inteiro em memória.
* Identificação automática do cabeçalho do CSV.
* Geração de PDF em formato A4 com tabela paginada automaticamente.
* Cabeçalho da tabela repetido em cada página.
* Linhas alternadas com fundo cinza para melhor legibilidade.
* Suporte a imagem de timbre (PNG/JPG/JPEG) como plano de fundo em todas as páginas.
* Validação completa dos parâmetros antes da geração.

---

## Uso

```bash
pdfgen create -file <arquivo.csv> -o <saida.pdf> [-c <timbre.png>]
```

### Parâmetros

| Parâmetro | Descrição                                       | Obrigatório |
| --------- | ----------------------------------------------- | :---------: |
| `create`  | Comando de geração do PDF.                      |     Sim     |
| `-file`   | Caminho do arquivo CSV de entrada.              |     Sim     |
| `-o`      | Caminho do arquivo PDF de saída.                |     Sim     |
| `-c`      | Imagem de timbre utilizada como plano de fundo. |     Não     |

### Exemplos

**Geração simples:**

```bash
pdfgen create -file ./access.log.csv -o ./relatorio.pdf
```

**Com papel timbrado:**

```bash
pdfgen create -file ./access.log.csv -o ./relatorio.pdf -c ./assets/timbre.png
```

**Ajuda:**

```bash
pdfgen help
pdfgen --help
```

---

## Requisitos

* Linux ou macOS.
* Go 1.25 ou superior para compilação a partir do código-fonte.

---

## Instalação

### Instalação rápida via `curl`

A forma recomendada de instalação é através do script de instalação:

```bash
curl -fsSL https://raw.githubusercontent.com/Nathan-rs/pdfgen/main/install.sh | sh
```

Após a instalação, verifique:

```bash
pdfgen --help
```

> O script instala o executável `pdfgen` no diretório apropriado para uso no terminal.

### Compilar a partir do código-fonte

```bash
git clone https://github.com/Nathan-rs/pdfgen.git
cd pdfgen
go build -o pdfgen ./cmd/pdfgen
```

### Executar diretamente com Go

```bash
go run ./cmd/pdfgen create -file dados.csv -o relatorio.pdf
```

---

## Estrutura do Projeto

```text
pdfgen/
├── cmd/
│   └── pdfgen/
│       └── main.go            # Entrypoint — roteamento de comandos
├── internal/
│   ├── cli/
│   │   ├── create.go          # Comando "create": parsing de flags e orquestração
│   │   └── help.go            # Comando "help": exibição de ajuda
│   ├── csv/
│   │   └── reader.go          # Pacote reservado para futuras abstrações de leitura
│   ├── pdf/
│   │   └── document.go        # Geração do PDF: layout, tabela, paginação e timbre
│   └── validate/
│       └── validate.go        # Validação de CSV, saída e imagem de timbre
├── assets/
│   └── timbre.png             # Imagem padrão do papel timbrado
├── go.mod
├── go.sum
├── install.sh                 # Script de instalação
└── README.md
```

---

## Configuração da Página

| Propriedade         |  Valor |
| ------------------- | -----: |
| Formato             |     A4 |
| Largura             | 210 mm |
| Altura              | 297 mm |
| Margem superior     |  55 mm |
| Margem inferior     |  35 mm |
| Margem esquerda     |  10 mm |
| Margem direita      |  10 mm |
| Área útil (largura) | 190 mm |

---

## Formato do CSV

* **Separador:** `;` (ponto e vírgula).
* A primeira linha é tratada como cabeçalho.
* A quantidade de colunas é detectada automaticamente.
* A largura disponível da tabela é distribuída igualmente entre as colunas.

### Exemplo — Log de acesso HTTP

```csv
timestamp;ip;method;path;status;response_time_ms;user_agent;user
2026-09-03 08:15:32;192.168.1.10;GET;/api/users;200;42;Mozilla/5.0;admin
2026-09-03 08:16:04;192.168.1.25;POST;/api/login;401;18;Mozilla/5.0;-
2026-09-03 08:17:21;192.168.1.10;GET;/api/reports;200;137;Chrome/140;joao
2026-09-03 08:18:45;10.0.0.15;GET;/api/users/42;404;23;Firefox/142;maria
```

Esse formato permite gerar, por exemplo, um relatório contendo:

* Data e hora da requisição;
* Endereço IP;
* Método HTTP;
* Endpoint acessado;
* Código HTTP retornado;
* Tempo de resposta;
* Navegador ou cliente utilizado;
* Usuário responsável pela requisição.

---

## Tecnologias

| Tecnologia                                    | Uso                      |
| --------------------------------------------- | ------------------------ |
| [Go 1.25](https://go.dev/)                    | Linguagem de programação |
| [gofpdf](https://github.com/phpdave11/gofpdf) | Geração do documento PDF |

---

## Validações

O sistema valida os parâmetros antes de iniciar o processamento:

| Validação | Regra                                                                                             |
| --------- | ------------------------------------------------------------------------------------------------- |
| `-file`   | Obrigatório. Deve existir, não ser diretório e ter extensão `.csv`.                               |
| `-o`      | Obrigatório. O diretório pai deve existir e a extensão deve ser `.pdf`.                           |
| `-c`      | Opcional. Se informado, deve existir, não ser diretório e ter extensão `.png`, `.jpg` ou `.jpeg`. |

---

## Roadmap

* [ ] Implementar leitor CSV dedicado (`internal/csv/reader.go`).
* [ ] Logger estruturado.
* [ ] Suporte a merge de múltiplos PDFs via `pdfcpu`.
* [ ] Inserção de capa como primeira página.
* [ ] Processamento concorrente com goroutines (Worker Pool / Batch Builder).
* [ ] Compressão e otimização do PDF final.
* [ ] Testes unitários e de integração.
* [ ] Configuração centralizada via YAML/JSON.
* [x] Script de instalação (`install.sh`).

---

## Licença

Este projeto é de uso interno. Consulte o autor para informações sobre licenciamento.
