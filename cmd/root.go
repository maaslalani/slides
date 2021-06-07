package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/model"
	"github.com/spf13/cobra"
)

const (
	delimiter    = "\n---\n"
	altDelimiter = "\n~~~\n"
)

var root = &cobra.Command{
	Use:   "slides <file.md>",
	Short: "Slides is a terminal based presentation tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		f := args[0]

		s, err := os.Stat(f)
		if err != nil {
			return errors.New("could not read file")
		}
		if s.IsDir() {
			return errors.New("must pass a file")
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			return errors.New("could not read file")
		}

		content := string(b)
		content = strings.ReplaceAll(content, altDelimiter, delimiter)
		slides := strings.Split(content, delimiter)

		user, err := user.Current()
		if err != nil {
			return errors.New("could not get current user")
		}

		p := tea.NewProgram(model.Model{
			Slides: slides,
			Page:   0,
			Author: user.Name,
			Date:   s.ModTime().Format("2006-01-03"),
		}, tea.WithAltScreen())

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
