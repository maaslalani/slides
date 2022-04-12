package main

import (
	_ "embed"
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
		Args: coral.ArbitraryArgs,
		RunE: func(cmd *coral.Command, args []string) error {
			var err error
			var fileName string

			if len(args) > 1 {
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

			p := tea.NewProgram(presentation, tea.WithAltScreen())
			err = p.Start()
			return err
		},
	}
)

func init() {
	rootCmd.AddCommand(
		cmd.ServeCmd,
	)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
