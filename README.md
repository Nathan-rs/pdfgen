# PDFGEN

Aplicação CLI em Go para gerar documentos PDF tabulares a partir de arquivos CSV, com suporte a papel timbrado institucional.

---

## Funcionalidades

- Leitura de arquivos CSV (separador `;`) via streaming — sem carregar o arquivo inteiro em memória.
- Identificação automática do cabeçalho do CSV.
- Geração de PDF em formato A4 com tabela paginada automaticamente.
- Cabeçalho da tabela repetido em cada página.
- Linhas alternadas com fundo cinza para melhor legibilidade.
- Suporte a imagem de timbre (PNG/JPG/JPEG) como plano de fundo em todas as páginas.
- Validação completa dos parâmetros antes da geração.

---

## Uso

```bash
pdfgen create -file <arquivo.csv> -o <saida.pdf> [-c <timbre.png>]
```

### Parâmetros

| Parâmetro | Descrição                                        | Obrigatório |
| --------- | ------------------------------------------------ | :---------: |
| `create`  | Comando de geração do PDF.                       |     Sim     |
| `-file`   | Caminho do arquivo CSV de entrada.               |     Sim     |
| `-o`      | Caminho do arquivo PDF de saída.                 |     Sim     |
| `-c`      | Imagem de timbre utilizada como plano de fundo.  |     Não     |

### Exemplos

**Geração simples:**

```bash
pdfgen create -file ./dados.csv -o ./relatorio.pdf
```

**Com papel timbrado:**

```bash
pdfgen create -file ./dados.csv -o ./relatorio.pdf -c ./assets/timbre.png
```

**Ajuda:**

```bash
pdfgen help
pdfgen --help
```

---

## Requisitos

- Go 1.25 ou superior.

## Instalação

### Compilar a partir do código-fonte

```bash
git clone https://github.com/Nathan-rs/pdfgen.git
cd pdfgen
go build -o pdfgen ./cmd/pdfgen
```

### Executar diretamente

```bash
go run ./cmd/pdfgen create -file dados.csv -o relatorio.pdf
```

---

## Estrutura do Projeto

```
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
└── README.md
```

---

## Configuração da Página

| Propriedade | Valor    |
| ----------- | -------: |
| Formato     |       A4 |
| Largura     |   210 mm |
| Altura      |   297 mm |
| Margem superior | 55 mm |
| Margem inferior | 35 mm |
| Margem esquerda | 10 mm |
| Margem direita  | 10 mm |
| Área útil (largura) | 190 mm |

---

## Formato do CSV

- **Separador:** `;` (ponto e vírgula).
- A primeira linha é tratada como cabeçalho.
- Quantidade de colunas é detectada automaticamente — a largura é distribuída igualmente entre elas.

Exemplo:

```csv
id;data;ip;ibge;route;method;id_user;name_user;device_type
1;2026-01-01;192.168.0.1;3550308;/api/users;GET;42;João Silva;desktop
```

---

## Tecnologias

| Tecnologia | Uso |
| ---------- | --- |
| [Go 1.25](https://go.dev/) | Linguagem de programação |
| [gofpdf](https://github.com/phpdave11/gofpdf) | Geração do documento PDF |

---

## Validações

O sistema valida os parâmetros antes de iniciar o processamento:

| Validação | Regra |
| --------- | ----- |
| `-file`   | Obrigatório. Deve existir, não ser diretório e ter extensão `.csv`. |
| `-o`      | Obrigatório. O diretório pai deve existir e a extensão deve ser `.pdf`. |
| `-c`      | Opcional. Se informado, deve existir, não ser diretório e ter extensão `.png`, `.jpg` ou `.jpeg`. |

---

## Roadmap

- [ ] Implementar leitor CSV dedicado (`internal/csv/reader.go`).
- [ ] Logger estruturado.
- [ ] Suporte a merge de múltiplos PDFs via `pdfcpu`.
- [ ] Inserção de capa como primeira página.
- [ ] Processamento concorrente com goroutines (Worker Pool / Batch Builder).
- [ ] Compressão e otimização do PDF final.
- [ ] Testes unitários e de integração.
- [ ] Configuração centralizada via YAML/JSON.
- [ ] Script de instalação (`install.sh`).

---

## Licença

Este projeto é de uso interno. Consulte o autor para informações sobre licenciamento.
