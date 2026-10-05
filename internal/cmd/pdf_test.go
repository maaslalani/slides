package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPDFOutput(t *testing.T) {
	assert.Equal(t, "slides.pdf", DefaultPDFOutput(""))
	assert.Equal(t, "presentation.pdf", DefaultPDFOutput("presentation.md"))
	assert.Equal(t, "talk.markdown.pdf", DefaultPDFOutput("talk.markdown.md"))
	assert.Equal(t, "notes.pdf", DefaultPDFOutput("notes.txt"))
	assert.Equal(t, "slides.pdf", DefaultPDFOutput("slides"))
}

func TestPDFCmd_Execution(t *testing.T) {
	tmpDir := t.TempDir()
	slideFile := filepath.Join(tmpDir, "presentation.md")
	content := "# Title Slide\n---\n# Slide 2\n"
	err := os.WriteFile(slideFile, []byte(content), 0644)
	require.NoError(t, err)

	outFile := filepath.Join(tmpDir, "output.pdf")
	pdfOutputFilePath = outFile
	defer func() {
		pdfOutputFilePath = ""
	}()

	err = PDFCmd.RunE(PDFCmd, []string{slideFile})
	require.NoError(t, err)

	info, err := os.Stat(outFile)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100))
}
