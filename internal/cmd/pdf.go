package cmd

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/internal/navigation"
	"github.com/maaslalani/slides/internal/pdf"
	"github.com/muesli/coral"
)

var (
	pdfOutputFilePath string
)

// PDFCmd exports presentation slides to a paginated PDF document.
var PDFCmd = &coral.Command{
	Use:     "pdf <file.md>",
	Aliases: []string{"export"},
	Short:   "Export presentation to a paginated PDF document",
	Args:    coral.ArbitraryArgs,
	RunE: func(cmd *coral.Command, args []string) error {
		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		}

		presentation := model.Model{
			Page:     0,
			Date:     time.Now().Format("2006-01-02"),
			FileName: fileName,
			Search:   navigation.NewSearch(),
		}
		err := presentation.Load()
		if err != nil {
			return err
		}

		outFile := pdfOutputFilePath
		if outFile == "" {
			outFile = DefaultPDFOutput(fileName)
		}

		return pdf.ExportFile(presentation, outFile)
	},
}

// DefaultPDFOutput computes the default PDF output file name based on input file name.
func DefaultPDFOutput(fileName string) string {
	if fileName == "" {
		return "slides.pdf"
	}
	ext := filepath.Ext(fileName)
	if ext != "" {
		return strings.TrimSuffix(fileName, ext) + ".pdf"
	}
	return fileName + ".pdf"
}

func init() {
	PDFCmd.Flags().StringVarP(&pdfOutputFilePath, "output", "o", "", "Output PDF file path")
}
