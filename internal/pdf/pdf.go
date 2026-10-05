package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/go-pdf/fpdf"
	"github.com/maaslalani/slides/internal/code"
	"github.com/maaslalani/slides/internal/model"
)

// Options holds configuration for exporting slides to PDF.
type Options struct {
	PageWidth       float64
	PageHeight      float64
	MarginLeft      float64
	MarginRight     float64
	MarginTop       float64
	MarginBottom    float64
	FontFamily      string
	CodeFontFamily  string
}

// DefaultOptions returns the standard presentation layout options (A4 Landscape).
func DefaultOptions() Options {
	return Options{
		PageWidth:      297.0, // A4 Landscape width in mm
		PageHeight:     210.0, // A4 Landscape height in mm
		MarginLeft:     20.0,
		MarginRight:    20.0,
		MarginTop:      20.0,
		MarginBottom:   25.0,
		FontFamily:     "Helvetica",
		CodeFontFamily: "Courier",
	}
}

// Exporter renders presentation models into paginated PDF documents.
type Exporter struct {
	opts Options
}

// New creates a new Exporter with default options.
func New() *Exporter {
	return &Exporter{
		opts: DefaultOptions(),
	}
}

// NewWithOptions creates a new Exporter with custom options.
func NewWithOptions(opts Options) *Exporter {
	return &Exporter{
		opts: opts,
	}
}

// Render compiles the presentation model into PDF byte stream.
func (e *Exporter) Render(m model.Model) ([]byte, error) {
	doc, err := e.Build(m)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, fmt.Errorf("could not generate PDF buffer: %w", err)
	}

	return buf.Bytes(), nil
}

// Build creates the fpdf.Fpdf instance for the presentation model.
func (e *Exporter) Build(m model.Model) (*fpdf.Fpdf, error) {
	doc := fpdf.NewCustom(&fpdf.InitType{
		UnitStr: "mm",
		Size:    fpdf.SizeType{Wd: e.opts.PageWidth, Ht: e.opts.PageHeight},
	})
	doc.SetMargins(e.opts.MarginLeft, e.opts.MarginTop, e.opts.MarginRight)
	doc.SetAutoPageBreak(false, e.opts.MarginBottom)

	tr := doc.UnicodeTranslatorFromDescriptor("")

	slides := m.Slides
	if len(slides) == 0 {
		slides = []string{""}
	}

	totalSlides := len(slides)

	for slideIdx, rawSlide := range slides {
		doc.AddPage()

		// Prepare slide content: strip comments and normalize tabs
		slideText := code.HideComments(rawSlide)
		slideText = strings.ReplaceAll(slideText, "\t", "    ")

		e.renderSlideContent(doc, tr, slideText)
		e.renderFooter(doc, tr, m, slideIdx, totalSlides)
	}

	if doc.Error() != nil {
		return nil, doc.Error()
	}

	return doc, nil
}

func (e *Exporter) renderSlideContent(doc *fpdf.Fpdf, tr func(string) string, content string) {
	usableWidth := e.opts.PageWidth - e.opts.MarginLeft - e.opts.MarginRight
	contentBottom := e.opts.PageHeight - e.opts.MarginBottom

	lines := strings.Split(content, "\n")
	idx := 0
	numLines := len(lines)

	for idx < numLines {
		if doc.GetY() >= contentBottom {
			break
		}

		line := lines[idx]
		trimmed := strings.TrimSpace(line)

		// 1. Code block handling
		if strings.HasPrefix(trimmed, "```") {
			codeLines := []string{}
			idx++
			for idx < numLines && !strings.HasPrefix(strings.TrimSpace(lines[idx]), "```") {
				codeLines = append(codeLines, lines[idx])
				idx++
			}
			if idx < numLines {
				idx++ // consume closing ```
			}
			e.renderCodeBlock(doc, tr, codeLines, usableWidth, contentBottom)
			continue
		}

		// 2. Table handling
		if isTableRow(trimmed) {
			tableLines := []string{}
			for idx < numLines && isTableRow(strings.TrimSpace(lines[idx])) {
				tableLines = append(tableLines, strings.TrimSpace(lines[idx]))
				idx++
			}
			e.renderTable(doc, tr, tableLines, usableWidth, contentBottom)
			continue
		}

		// 3. Blank line
		if trimmed == "" {
			doc.SetY(doc.GetY() + 3.0)
			idx++
			continue
		}

		// 4. Horizontal rule
		if isHorizontalRule(trimmed) {
			y := doc.GetY() + 2.0
			doc.SetDrawColor(210, 214, 220)
			doc.SetLineWidth(0.4)
			doc.Line(e.opts.MarginLeft, y, e.opts.MarginLeft+usableWidth, y)
			doc.SetY(y + 3.0)
			idx++
			continue
		}

		// 5. Headings
		if strings.HasPrefix(trimmed, "#") {
			e.renderHeading(doc, tr, trimmed, usableWidth)
			idx++
			continue
		}

		// 6. Blockquote
		if strings.HasPrefix(trimmed, ">") {
			quoteText := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
			e.renderBlockquote(doc, tr, quoteText, usableWidth)
			idx++
			continue
		}

		// 7. Bullet lists
		if isBulletItem(trimmed) {
			bulletText := stripBulletPrefix(trimmed)
			e.renderBulletItem(doc, tr, bulletText, usableWidth)
			idx++
			continue
		}

		// 8. Numbered lists
		if num, text, ok := parseNumberedItem(trimmed); ok {
			e.renderNumberedItem(doc, tr, num, text, usableWidth)
			idx++
			continue
		}

		// 9. Standard paragraph text
		e.renderParagraph(doc, tr, line, usableWidth)
		idx++
	}
}

func (e *Exporter) renderHeading(doc *fpdf.Fpdf, tr func(string) string, line string, usableWidth float64) {
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	text := strings.TrimSpace(line[level:])

	var fontSize float64
	var lineHeight float64
	var spaceAfter float64

	switch level {
	case 1:
		fontSize = 24.0
		lineHeight = 11.0
		spaceAfter = 4.0
		doc.SetTextColor(17, 24, 39)
	case 2:
		fontSize = 18.0
		lineHeight = 9.0
		spaceAfter = 3.5
		doc.SetTextColor(31, 41, 55)
	case 3:
		fontSize = 15.0
		lineHeight = 8.0
		spaceAfter = 3.0
		doc.SetTextColor(55, 65, 81)
	case 4:
		fontSize = 13.0
		lineHeight = 7.0
		spaceAfter = 2.5
		doc.SetTextColor(75, 85, 99)
	case 5:
		fontSize = 11.5
		lineHeight = 6.0
		spaceAfter = 2.0
		doc.SetTextColor(107, 114, 128)
	default:
		fontSize = 10.5
		lineHeight = 5.5
		spaceAfter = 2.0
		doc.SetTextColor(107, 114, 128)
	}

	doc.SetFont(e.opts.FontFamily, "B", fontSize)
	doc.MultiCell(usableWidth, lineHeight, tr(text), "", "L", false)
	doc.SetY(doc.GetY() + spaceAfter)
}

func (e *Exporter) renderCodeBlock(doc *fpdf.Fpdf, tr func(string) string, codeLines []string, usableWidth, contentBottom float64) {
	lineHeight := 4.6
	padding := 3.0
	blockHeight := float64(len(codeLines))*lineHeight + padding*2

	currentY := doc.GetY()
	if currentY+blockHeight > contentBottom {
		blockHeight = contentBottom - currentY
		if blockHeight < padding*2 {
			return
		}
	}

	// Code block background
	doc.SetFillColor(245, 246, 248)
	doc.SetDrawColor(220, 224, 230)
	doc.SetLineWidth(0.3)
	doc.Rect(e.opts.MarginLeft, currentY, usableWidth, blockHeight, "DF")

	doc.SetFont(e.opts.CodeFontFamily, "", 9.5)
	doc.SetTextColor(36, 41, 46)

	textY := currentY + padding
	for _, cl := range codeLines {
		if textY+lineHeight > currentY+blockHeight {
			break
		}
		doc.SetXY(e.opts.MarginLeft+padding, textY)
		doc.CellFormat(usableWidth-padding*2, lineHeight, tr(cl), "", 0, "L", false, 0, "")
		textY += lineHeight
	}

	doc.SetY(currentY + blockHeight + 3.0)
}

func (e *Exporter) renderBlockquote(doc *fpdf.Fpdf, tr func(string) string, text string, usableWidth float64) {
	startX := e.opts.MarginLeft
	startY := doc.GetY()

	doc.SetFont(e.opts.FontFamily, "I", 11.0)
	doc.SetTextColor(75, 85, 99)

	indent := 6.0
	doc.SetX(startX + indent)
	doc.MultiCell(usableWidth-indent, 6.0, tr(text), "", "L", false)
	endY := doc.GetY()

	// Left vertical bar
	doc.SetDrawColor(156, 163, 175)
	doc.SetLineWidth(1.2)
	doc.Line(startX+1.0, startY, startX+1.0, endY)

	doc.SetY(endY + 2.0)
}

func (e *Exporter) renderBulletItem(doc *fpdf.Fpdf, tr func(string) string, text string, usableWidth float64) {
	startX := e.opts.MarginLeft
	startY := doc.GetY()

	// Draw clean bullet
	doc.SetFillColor(79, 70, 229)
	bulletRadius := 0.75
	doc.Circle(startX+2.5, startY+3.0, bulletRadius, "F")

	indent := 6.0
	doc.SetXY(startX+indent, startY)
	doc.SetFont(e.opts.FontFamily, "", 11.0)
	doc.SetTextColor(31, 41, 55)
	doc.MultiCell(usableWidth-indent, 5.8, tr(text), "", "L", false)
	doc.SetY(doc.GetY() + 1.0)
}

func (e *Exporter) renderNumberedItem(doc *fpdf.Fpdf, tr func(string) string, num, text string, usableWidth float64) {
	startX := e.opts.MarginLeft
	startY := doc.GetY()

	doc.SetFont(e.opts.FontFamily, "B", 11.0)
	doc.SetTextColor(79, 70, 229)
	numPrefix := num + "."
	doc.SetXY(startX, startY)
	doc.CellFormat(7.0, 5.8, numPrefix, "", 0, "L", false, 0, "")

	indent := 8.0
	doc.SetXY(startX+indent, startY)
	doc.SetFont(e.opts.FontFamily, "", 11.0)
	doc.SetTextColor(31, 41, 55)
	doc.MultiCell(usableWidth-indent, 5.8, tr(text), "", "L", false)
	doc.SetY(doc.GetY() + 1.0)
}

func (e *Exporter) renderTable(doc *fpdf.Fpdf, tr func(string) string, tableLines []string, usableWidth, contentBottom float64) {
	rows := [][]string{}
	for _, l := range tableLines {
		if isTableSeparator(l) {
			continue
		}
		cols := parseTableColumns(l)
		if len(cols) > 0 {
			rows = append(rows, cols)
		}
	}

	if len(rows) == 0 {
		return
	}

	colCount := 0
	for _, r := range rows {
		if len(r) > colCount {
			colCount = len(r)
		}
	}
	if colCount == 0 {
		return
	}

	colWidth := usableWidth / float64(colCount)
	rowHeight := 6.5

	for rowIdx, row := range rows {
		if doc.GetY()+rowHeight > contentBottom {
			break
		}

		isHeader := (rowIdx == 0)
		if isHeader {
			doc.SetFont(e.opts.FontFamily, "B", 10.0)
			doc.SetFillColor(243, 244, 246)
			doc.SetTextColor(17, 24, 39)
		} else {
			doc.SetFont(e.opts.FontFamily, "", 10.0)
			if rowIdx%2 == 1 {
				doc.SetFillColor(255, 255, 255)
			} else {
				doc.SetFillColor(249, 250, 251)
			}
			doc.SetTextColor(55, 65, 81)
		}

		doc.SetDrawColor(209, 213, 219)
		doc.SetLineWidth(0.2)

		startX := e.opts.MarginLeft
		startY := doc.GetY()

		for colIdx := 0; colIdx < colCount; colIdx++ {
			cellText := ""
			if colIdx < len(row) {
				cellText = row[colIdx]
			}
			doc.SetXY(startX+float64(colIdx)*colWidth, startY)
			doc.CellFormat(colWidth, rowHeight, tr(cellText), "1", 0, "L", true, 0, "")
		}
		doc.SetY(startY + rowHeight)
	}

	doc.SetY(doc.GetY() + 3.0)
}

func (e *Exporter) renderParagraph(doc *fpdf.Fpdf, tr func(string) string, line string, usableWidth float64) {
	doc.SetFont(e.opts.FontFamily, "", 11.0)
	doc.SetTextColor(31, 41, 55)
	cleanLine := stripInlineCode(line)
	doc.MultiCell(usableWidth, 5.8, tr(cleanLine), "", "L", false)
	doc.SetY(doc.GetY() + 1.2)
}

func (e *Exporter) renderFooter(doc *fpdf.Fpdf, tr func(string) string, m model.Model, slideIdx, totalSlides int) {
	footerY := e.opts.PageHeight - 15.0
	lineY := footerY - 4.0
	usableWidth := e.opts.PageWidth - e.opts.MarginLeft - e.opts.MarginRight

	// Footer separator line
	doc.SetDrawColor(228, 228, 231)
	doc.SetLineWidth(0.3)
	doc.Line(e.opts.MarginLeft, lineY, e.opts.MarginLeft+usableWidth, lineY)

	// Left footer: Author + Date
	doc.SetFont(e.opts.FontFamily, "", 9.0)
	doc.SetTextColor(113, 113, 122)

	leftParts := []string{}
	if m.Author != "" {
		leftParts = append(leftParts, m.Author)
	}
	if m.Date != "" {
		leftParts = append(leftParts, m.Date)
	}
	leftText := strings.Join(leftParts, "  |  ")

	doc.SetXY(e.opts.MarginLeft, footerY)
	doc.CellFormat(usableWidth/2, 5.0, tr(leftText), "", 0, "L", false, 0, "")

	// Right footer: Paging
	pagingStr := formatPaging(m.Paging, slideIdx+1, totalSlides)
	doc.SetXY(e.opts.MarginLeft+usableWidth/2, footerY)
	doc.CellFormat(usableWidth/2, 5.0, tr(pagingStr), "", 0, "R", false, 0, "")
}

func formatPaging(fmtPattern string, currentPage, totalPages int) string {
	switch strings.Count(fmtPattern, "%d") {
	case 2:
		return fmt.Sprintf(fmtPattern, currentPage, totalPages)
	case 1:
		return fmt.Sprintf(fmtPattern, currentPage)
	default:
		if fmtPattern != "" {
			return fmtPattern
		}
		return fmt.Sprintf("%d / %d", currentPage, totalPages)
	}
}

// Export renders the presentation model and writes the resulting PDF to w.
func Export(m model.Model, w io.Writer) error {
	exporter := New()
	data, err := exporter.Render(m)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// ExportFile renders the presentation model and writes the PDF to targetPath.
func ExportFile(m model.Model, targetPath string) error {
	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("could not create output file %s: %w", targetPath, err)
	}
	defer file.Close()

	return Export(m, file)
}

func isHorizontalRule(line string) bool {
	if line == "---" || line == "***" || line == "___" {
		return true
	}
	return false
}

func isBulletItem(line string) bool {
	return strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "+ ")
}

func stripBulletPrefix(line string) string {
	if strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "+ ") {
		return strings.TrimSpace(line[2:])
	}
	return line
}

var numListRegex = regexp.MustCompile(`^(\d+)\.\s+(.*)$`)

func parseNumberedItem(line string) (string, string, bool) {
	matches := numListRegex.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2], true
	}
	return "", "", false
}

func isTableRow(line string) bool {
	return strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|")
}

func isTableSeparator(line string) bool {
	clean := strings.ReplaceAll(line, "|", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, ":", "")
	clean = strings.TrimSpace(clean)
	return clean == ""
}

func parseTableColumns(line string) []string {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil
	}
	// Trim first and last empty slices from outer pipes
	if len(parts) > 1 && parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	cols := make([]string, len(parts))
	for i, p := range parts {
		cols[i] = strings.TrimSpace(p)
	}
	return cols
}

func stripInlineCode(s string) string {
	return strings.ReplaceAll(s, "`", "")
}
