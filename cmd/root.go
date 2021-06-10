package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
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

		user, err := user.Current()
		if err != nil {
			return errors.New("could not get current user")
		}

		presentation := model.Model{
			Page:     0,
			Author:   user.Name,
			Date:     time.Now().Format("2006-01-02"),
			FileName: fileName,
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
