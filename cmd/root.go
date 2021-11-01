package cmd

import (
	"fmt"
	"github.com/maaslalani/slides/internal/navigation"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/model"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "slides <file.md>",
	Short: "Slides is a terminal based presentation tool",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		p := tea.NewProgram(presentation, tea.WithAltScreen())
		err = p.Start()
		return err
	},
}

func Execute() {
	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
