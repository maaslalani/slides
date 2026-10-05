package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd_PDFExportFlag(t *testing.T) {
	tmpDir := t.TempDir()
	slidePath := filepath.Join(tmpDir, "deck.md")
	content := "# Slide One\n---\n# Slide Two\n"
	err := os.WriteFile(slidePath, []byte(content), 0644)
	require.NoError(t, err)

	pdfPath := filepath.Join(tmpDir, "deck.pdf")

	rootCmd.SetArgs([]string{slidePath, "--pdf", pdfPath})
	err = rootCmd.Execute()
	require.NoError(t, err)

	info, err := os.Stat(pdfPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100))

	// Reset flag state
	pdfOutput = ""
}

func TestRootCmd_PDFSubcommand(t *testing.T) {
	tmpDir := t.TempDir()
	slidePath := filepath.Join(tmpDir, "presentation.md")
	content := "# Main Slide\nContent\n"
	err := os.WriteFile(slidePath, []byte(content), 0644)
	require.NoError(t, err)

	outPath := filepath.Join(tmpDir, "presentation.pdf")

	rootCmd.SetArgs([]string{"pdf", slidePath, "-o", outPath})
	err = rootCmd.Execute()
	require.NoError(t, err)

	info, err := os.Stat(outPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100))
}

func TestRootCmd_ExportSubcommandAlias(t *testing.T) {
	tmpDir := t.TempDir()
	slidePath := filepath.Join(tmpDir, "alias.md")
	content := "# Alias Slide\nContent\n"
	err := os.WriteFile(slidePath, []byte(content), 0644)
	require.NoError(t, err)

	outPath := filepath.Join(tmpDir, "alias.pdf")

	rootCmd.SetArgs([]string{"export", slidePath, "-o", outPath})
	err = rootCmd.Execute()
	require.NoError(t, err)

	info, err := os.Stat(outPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100))
}
