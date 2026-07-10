package pdf

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/phpdave11/gofpdf"
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

// Largura total disponível para a tabela
// var usableWidth = pageWidth - marginLeft - marginRight

func New(csvFile, cover string) (*Document, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.Comma = ';'

	headers, err := r.Read()
	if err != nil {
		f.Close()
		return nil, err
	}

	d := &Document{
		file:        f,
		reader:      r,
		pdf:         gofpdf.New("P", "mm", "A4", ""),
		headers:     headers,
		cover:       cover,
		usableWidth: pageWidth - marginLeft - marginRight,
	}

	d.calculateColWidths()
	d.pdf.SetMargins(marginLeft, marginTop, marginRight)
	d.pdf.SetAutoPageBreak(false, marginBottom)
	d.newPage()

	return d, nil
}

func (d *Document) Save(output string) error {
	defer d.file.Close()

	tr := d.pdf.UnicodeTranslatorFromDescriptor("")

	for {
		record, err := d.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

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
	return d.pdf.OutputFileAndClose(output)
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
	tr := d.pdf.UnicodeTranslatorFromDescriptor("")
	for i, h := range d.headers {
		d.pdf.CellFormat(d.colWidths[i], 6, tr(h), "1", 0, "C", true, 0, "")
	}
	d.pdf.Ln(6)
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
