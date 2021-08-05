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

var userName string

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

		if userName == "" {
			user, err := user.Current()
			if err != nil {
				return errors.New("could not get current user")
			}
			userName = user.Name
		}

		presentation := model.Model{
			Page:     0,
			Author:   userName,
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
	root.Flags().StringVarP(&userName, "username", "u", "", "Custom user name to show in the footer")
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
