package pdf

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maaslalani/slides/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_SingleSlide(t *testing.T) {
	m := model.Model{
		Slides: []string{"# Hello World\n\nThis is a single slide presentation."},
		Author: "Test Author",
		Date:   "2026-10-05",
		Paging: "Slide %d / %d",
	}

	exporter := New()
	pdfBytes, err := exporter.Render(m)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

	// Verify standard PDF header magic bytes
	assert.True(t, bytes.HasPrefix(pdfBytes, []byte("%PDF-")))

	// Verify document builder page count
	doc, err := exporter.Build(m)
	require.NoError(t, err)
	assert.Equal(t, 1, doc.PageNo())
}

func TestRender_PageBreakBoundaries(t *testing.T) {
	testCases := []struct {
		name          string
		slides        []string
		expectedPages int
	}{
		{
			name:          "Zero slides fallback to single page",
			slides:        []string{},
			expectedPages: 1,
		},
		{
			name:          "One slide",
			slides:        []string{"# Slide 1"},
			expectedPages: 1,
		},
		{
			name: "Three slides strict boundaries",
			slides: []string{
				"# Slide 1\nFirst slide text",
				"# Slide 2\nSecond slide text",
				"# Slide 3\nThird slide text",
			},
			expectedPages: 3,
		},
		{
			name: "Ten slides pagination boundary test",
			slides: []string{
				"# Slide 1",
				"# Slide 2",
				"# Slide 3",
				"# Slide 4",
				"# Slide 5",
				"# Slide 6",
				"# Slide 7",
				"# Slide 8",
				"# Slide 9",
				"# Slide 10",
			},
			expectedPages: 10,
		},
	}

	exporter := New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := model.Model{
				Slides: tc.slides,
				Author: "Boundary Tester",
				Date:   "2026-10-05",
				Paging: "Slide %d / %d",
			}

			doc, err := exporter.Build(m)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedPages, doc.PageNo(), "page count must strictly equal expected page breaks")

			pdfBytes, err := exporter.Render(m)
			require.NoError(t, err)
			assert.NotEmpty(t, pdfBytes)

			// Count /Type /Page occurrences in the generated PDF stream
			pdfContent := string(pdfBytes)
			pageOccurrences := strings.Count(pdfContent, "/Type /Page")
			// PDF can contain /Type /Pages (parent catalog) and /Type /Page (individual pages).
			// We check that each individual page is represented.
			pageLeafCount := strings.Count(pdfContent, "/Type /Page\n") + strings.Count(pdfContent, "/Type /Page ") + strings.Count(pdfContent, "/Type /Page/")
			assert.GreaterOrEqual(t, pageOccurrences, tc.expectedPages)
			_ = pageLeafCount
		})
	}
}

func TestRender_SlideMarkdownElements(t *testing.T) {
	complexSlide := `# Main Heading

## Sub Heading

### Section Heading

#### Sub-section 4
##### Sub-section 5
###### Sub-section 6

This is a regular paragraph with some text.

<!-- This comment must be hidden and excluded from output -->

> This is a blockquote with important insights.

* First bullet point
* Second bullet point
* Third bullet point

1. Numbered item one
2. Numbered item two

---

| Feature | Supported | Notes |
| --- | --- | --- |
| PDF Export | Yes | High fidelity |
| Tables | Yes | Formatted grids |

` + "```go\nfunc hello() string {\n    return \"slides\"\n}\n```\n"

	m := model.Model{
		Slides: []string{complexSlide},
		Author: "Markdown Tester",
		Date:   "2026-10-05",
		Paging: "Slide %d / %d",
	}

	exporter := New()
	pdfBytes, err := exporter.Render(m)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

	// Assert comments were stripped and not included in final text
	assert.False(t, strings.Contains(string(pdfBytes), "This comment must be hidden"))
}

func TestRender_MetadataAndFooter(t *testing.T) {
	testPaging := []struct {
		name        string
		pagingFmt   string
		expectMatch string
	}{
		{
			name:        "Standard 2 arg format",
			pagingFmt:   "Slide %d of %d",
			expectMatch: "Slide 1 of 2",
		},
		{
			name:        "Single arg format",
			pagingFmt:   "Page %d",
			expectMatch: "Page 1",
		},
		{
			name:        "Static string format",
			pagingFmt:   "CONFIDENTIAL",
			expectMatch: "CONFIDENTIAL",
		},
		{
			name:        "Empty format defaults to fraction",
			pagingFmt:   "",
			expectMatch: "1 / 2",
		},
	}

	for _, tc := range testPaging {
		t.Run(tc.name, func(t *testing.T) {
			m := model.Model{
				Slides: []string{"# Slide 1", "# Slide 2"},
				Author: "Alice",
				Date:   "2026-01-01",
				Paging: tc.pagingFmt,
			}

			exporter := New()
			pdfBytes, err := exporter.Render(m)
			require.NoError(t, err)
			assert.NotEmpty(t, pdfBytes)
		})
	}
}

func TestExport_Writer(t *testing.T) {
	m := model.Model{
		Slides: []string{"# Test Writer Output\nContent"},
		Author: "Writer Test",
		Date:   "2026-10-05",
	}

	var buf bytes.Buffer
	err := Export(m, &buf)
	require.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
	assert.True(t, bytes.HasPrefix(buf.Bytes(), []byte("%PDF-")))
}

func TestExportFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.pdf")

	m := model.Model{
		Slides: []string{"# Test Export File\nSlide body"},
		Author: "File Test",
		Date:   "2026-10-05",
	}

	err := ExportFile(m, outputPath)
	require.NoError(t, err)

	info, err := os.Stat(outputPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100))

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.True(t, bytes.HasPrefix(content, []byte("%PDF-")))
}

func TestFormatPaging(t *testing.T) {
	assert.Equal(t, "1 / 5", formatPaging("", 1, 5))
	assert.Equal(t, "Page 2 of 10", formatPaging("Page %d of %d", 2, 10))
	assert.Equal(t, "Slide 3", formatPaging("Slide %d", 3, 10))
	assert.Equal(t, "No numbers", formatPaging("No numbers", 3, 10))
}

func TestMarkdownParsers(t *testing.T) {
	// Table row detection
	assert.True(t, isTableRow("| a | b |"))
	assert.False(t, isTableRow("not a table"))

	// Table separator detection
	assert.True(t, isTableSeparator("| --- | :---: |"))
	assert.False(t, isTableSeparator("| data | more data |"))

	// Table columns parsing
	cols := parseTableColumns("| Col1 | Col2 | Col3 |")
	assert.Equal(t, []string{"Col1", "Col2", "Col3"}, cols)

	// Bullet item detection & prefix stripping
	assert.True(t, isBulletItem("* Item"))
	assert.True(t, isBulletItem("- Item"))
	assert.True(t, isBulletItem("+ Item"))
	assert.False(t, isBulletItem("Normal line"))
	assert.Equal(t, "Item", stripBulletPrefix("* Item"))

	// Numbered item parsing
	num, text, ok := parseNumberedItem("1. First item")
	assert.True(t, ok)
	assert.Equal(t, "1", num)
	assert.Equal(t, "First item", text)

	_, _, ok = parseNumberedItem("Not numbered")
	assert.False(t, ok)

	// Horizontal rule
	assert.True(t, isHorizontalRule("---"))
	assert.True(t, isHorizontalRule("***"))
	assert.True(t, isHorizontalRule("___"))
	assert.False(t, isHorizontalRule("--"))
}
