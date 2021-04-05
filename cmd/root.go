package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/model"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "slides <file.md>",
	Short: "Slides is a terminal based presentation tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, _ []string) error {
		p := tea.NewProgram(model.Model{
			Slides: []string{
				"Slide 1",
				"Slide 2",
				"Slide 3",
				"Slide 4",
				"Slide 5",
				"Slide 6",
			},
			Page:   0,
			Author: "Maas Lalani",
			Date:   "2021-04-04",
		})

		p.EnterAltScreen()
		err := p.Start()
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
