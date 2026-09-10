package pdf

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/phpdave11/gofpdf"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	pageWidth  = 210.0
	pageHeight = 297.0

	marginTop    = 55.0
	marginBottom = 35.0
	marginLeft   = 10.0 // Reduzido para 10mm para acomodar todas as colunas
	marginRight  = 10.0 // Reduzido para 10mm para acomodar todas as colunas

	maxY = pageHeight - marginBottom
)

type Document struct {
	file        *os.File
	reader      *csv.Reader
	pdf         *gofpdf.Fpdf
	headers     []string
	colWidths   []float64
	colX        []float64
	usableWidth float64
	cover       string
	recordCount int
}

// New abre o CSV e prepara um Document para o fluxo original,
// sequencial: New + Save lê o arquivo inteiro e grava o PDF final.
func New(csvFile, cover string) (*Document, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	// BOMOverride detecta e descarta um BOM (UTF-8/UTF-16LE/UTF-16BE) no
	// início do arquivo, comum em CSVs exportados pelo Excel/Windows.
	// Sem BOM, cai no decoder de fallback (UTF-8) sem alterar os bytes.
	decoder := unicode.BOMOverride(unicode.UTF8.NewDecoder())
	r := csv.NewReader(transform.NewReader(f, decoder))
	r.Comma = ';'

	headers, err := r.Read()
	if err != nil {
		f.Close()
		return nil, err
	}

	d := newDocument(headers, cover, 0)
	d.file = f
	d.reader = r

	return d, nil
}

// NewBatch cria um Document que não lê CSV nenhum: os registros chegam
// via AddRecords. Usado pela geração concorrente, onde cada worker tem
// seu próprio Document/gofpdf.Fpdf (gofpdf.Fpdf não é thread-safe, então
// não pode ser compartilhado entre goroutines).
//
// startIndex é o número de registros já escritos em lotes anteriores —
// serve só para o zebrado (linhas alternadas) continuar coerente entre
// as partes depois do merge; não afeta corretude.
func NewBatch(headers []string, cover string, startIndex int) *Document {
	return newDocument(headers, cover, startIndex)
}

func newDocument(headers []string, cover string, startIndex int) *Document {
	d := &Document{
		pdf:         gofpdf.New("P", "mm", "A4", ""),
		headers:     headers,
		cover:       cover,
		usableWidth: pageWidth - marginLeft - marginRight,
		recordCount: startIndex,
	}

	d.calculateColWidths()
	d.pdf.SetMargins(marginLeft, marginTop, marginRight)
	d.pdf.SetAutoPageBreak(false, marginBottom)
	d.newPage()

	return d
}

// Save lê o restante do CSV (via New) e grava o PDF final.
// Mantido para o caminho sequencial / compatibilidade.
func (d *Document) Save(output string) error {
	if d.file != nil {
		defer d.file.Close()
	}

	if d.reader != nil {
		for {
			record, err := d.reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			d.writeRecord(record)
		}
	}

	return d.pdf.OutputFileAndClose(output)
}

// AddRecords renderiza um lote de registros já lidos em memória.
// Usado pela geração concorrente (ver internal/pdf/concurrent.go).
func (d *Document) AddRecords(records [][]string) {
	for _, record := range records {
		d.writeRecord(record)
	}
}

// Close finaliza e grava o PDF. Usado em conjunto com NewBatch/AddRecords.
func (d *Document) Close(output string) error {
	return d.pdf.OutputFileAndClose(output)
}

func (d *Document) writeRecord(record []string) {
	tr := d.pdf.UnicodeTranslatorFromDescriptor("")

	d.recordCount++
	d.pdf.SetFont("Arial", "", 5.5)
	d.pdf.SetTextColor(50, 50, 50)
	d.pdf.SetDrawColor(220, 220, 220)

	const lineH = 4.5
	rowHeight := lineH

	for i, w := range d.colWidths {
		v := ""
		if i < len(record) {
			v = record[i]
		}

		lines := d.pdf.SplitLines([]byte(tr(v)), w-1)
		if h := float64(len(lines)) * lineH; h > rowHeight {
			rowHeight = h
		}
	}

	if d.pdf.GetY()+rowHeight >= maxY {
		d.newPage()
	}

	startY := d.pdf.GetY()
	fill := d.recordCount%2 == 0
	if fill {
		d.pdf.SetFillColor(245, 245, 245)
	} else {
		d.pdf.SetFillColor(255, 255, 255)
	}
	d.pdf.Rect(marginLeft, startY, d.usableWidth, rowHeight, "F")

	for i, w := range d.colWidths {
		v := ""
		if i < len(record) {
			v = record[i]
		}
		d.pdf.Rect(d.colX[i], startY, w, rowHeight, "D")
		d.pdf.SetXY(d.colX[i]+0.5, startY)
		d.pdf.MultiCell(w-0.5, lineH, tr(v), "", "L", false)
	}
	d.pdf.SetY(startY + rowHeight)
}

func (d *Document) addBackground() {
	if d.cover == "" {
		return
	}
	d.pdf.Image(d.cover, 0, 0, pageWidth, pageHeight, false, "", 0, "")
}

func (d *Document) drawTableHeader() {
	d.pdf.SetFont("Arial", "B", 6.5)
	d.pdf.SetFillColor(0, 168, 89)
	d.pdf.SetTextColor(255, 255, 255)
	d.pdf.SetDrawColor(200, 200, 200)

	headerHeight := d.headerHeight()

	x := marginLeft
	y := d.pdf.GetY()

	for i, h := range d.headers {
		d.drawHeaderCell(x, y, d.colWidths[i], headerHeight, h)

		x += d.colWidths[i]
	}

	d.pdf.SetY(y + headerHeight)

	// restaura as configurações usadas pelo corpo da tabela
	d.pdf.SetFont("Arial", "", 5.5)
	d.pdf.SetTextColor(50, 50, 50)
	d.pdf.SetDrawColor(220, 220, 220)

}

func (d *Document) newPage() {
	d.pdf.AddPage()
	d.addBackground()

	d.pdf.SetFont("Arial", "B", 14)
	d.pdf.SetTextColor(30, 30, 30)

	d.pdf.SetXY(marginLeft, marginTop-15)

	d.drawTableHeader()
}

func (d *Document) calculateColWidths() {
	n := len(d.headers)
	base := d.usableWidth / float64(n)
	d.colWidths = make([]float64, n)
	d.colX = make([]float64, n)

	for i := range d.colWidths {
		d.colWidths[i] = base
	}

	d.colX[0] = marginLeft

	for i := 1; i < n; i++ {
		d.colX[i] = d.colX[i-1] + d.colWidths[i-1]
	}
}

// Calcula a altura do header do csv
func (d *Document) headerHeight() float64 {
	tr := d.pdf.UnicodeTranslatorFromDescriptor("")

	const (
		lineH   = 3.5
		padding = 1.0
	)

	height := lineH + (padding * 2)

	for i, h := range d.headers {
		if i >= len(d.colWidths) {
			break
		}

		lines := d.pdf.SplitLines(
			[]byte(tr(h)),
			d.colWidths[i]-(padding*2),
		)

		headerH := float64(len(lines))*lineH + (padding * 2)

		if headerH > height {
			height = headerH
		}
	}

	return height
}

// Quebra em varias linhas quando o header for muito grande
// Deixando um padding maior quando o header for muito grande
func (d *Document) drawHeaderCell(
	x, y, width, height float64,
	text string,
) {
	tr := d.pdf.UnicodeTranslatorFromDescriptor("")

	const (
		lineH   = 3.5
		padding = 1.0
	)

	d.pdf.Rect(x, y, width, height, "FD")

	lines := d.pdf.SplitLines(
		[]byte(tr(text)),
		width-(padding*2),
	)

	textHeight := float64(len(lines)) * lineH
	textY := y + (height-textHeight)/2

	d.pdf.SetXY(x+padding, textY)

	d.pdf.MultiCell(
		width-(padding*2),
		lineH,
		tr(text),
		"",
		"C",
		false,
	)
}
