package main

import (
	_ "embed"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/maaslalani/slides/internal/cmd"
	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/internal/navigation"
	"github.com/maaslalani/slides/internal/pdf"
	"github.com/muesli/coral"
)

var (
	pdfOutput string

	rootCmd = &coral.Command{
		Use:   "slides <file.md>",
		Short: "Terminal based presentation tool",
		Args:  coral.ArbitraryArgs,
		RunE: func(c *coral.Command, args []string) error {
			var err error
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
			err = presentation.Load()
			if err != nil {
				return err
			}

			if pdfOutput != "" {
				outFile := pdfOutput
				if outFile == "default" {
					if len(args) > 1 {
						outFile = args[1]
					} else {
						outFile = cmd.DefaultPDFOutput(fileName)
					}
				}
				return pdf.ExportFile(presentation, outFile)
			}

			p := tea.NewProgram(presentation)
			_, err = p.Run()
			return err
		},
	}
)

func init() {
	rootCmd.AddCommand(
		cmd.ServeCmd,
		cmd.PDFCmd,
	)
	rootCmd.Flags().StringVarP(&pdfOutput, "pdf", "p", "", "Export presentation to PDF")
	rootCmd.Flags().Lookup("pdf").NoOptDefVal = "default"
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
