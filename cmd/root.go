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
	"github.com/maaslalani/slides/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		theme := styles.Theme
		if viper.GetString("theme") != "" {
			customTheme, err := styles.CustomTheme(viper.GetString("theme"))
			if err != nil {
				return err
			}
			theme = customTheme
		}

		p := tea.NewProgram(model.Model{
			Slides: slides,
			Page:   0,
			Author: user.Name,
			Date:   s.ModTime().Format("2006-01-02"),
			Theme:  theme,
		}, tea.WithAltScreen())

		err = p.Start()
		return err
	},
}

func init() {
	root.PersistentFlags().StringP("theme", "t", "", "custom theme.json")
	_ = viper.BindPFlag("theme", root.PersistentFlags().Lookup("theme"))
}

func Execute() {
	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
