package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/cmd"
	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/internal/navigation"
	"github.com/muesli/coral"
)

var (
	rootCmd = &coral.Command{
		Use:   "slides <file.md>",
		Short: "Terminal based presentation tool",
		Args:  coral.ArbitraryArgs,
		RunE: func(cmd *coral.Command, args []string) error {
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

			switch outputMode {
			case "", "tui":
				p := tea.NewProgram(presentation, tea.WithAltScreen())
				_, err = p.Run()
				return err
			case "stdout":
				return runStdout(presentation, os.Stdout)
			default:
				return fmt.Errorf("unsupported output mode %q", outputMode)
			}
		},
	}

	outputMode string
)

func init() {
	rootCmd.AddCommand(
		cmd.ServeCmd,
	)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVar(&outputMode, "output", "tui", "render output (tui or stdout)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runStdout(presentation model.Model, w io.Writer) error {
	for i := range presentation.Slides {
		rendered, err := presentation.RenderPage(i, 0)
		if err != nil {
			return err
		}
		if i > 0 {
			if _, err := fmt.Fprintln(w, "---"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, rendered); err != nil {
			return err
		}
	}
	return nil
}
