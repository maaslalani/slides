package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/meta"
	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/styles"
	"github.com/spf13/cobra"
)

const (
	delimiter    = "\n---\n"
	altDelimiter = "\n~~~\n"
)

var root = &cobra.Command{
	Use:   "slides <file.md>",
	Short: "Slides is a terminal based presentation tool",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		var content string
		var err error

		if len(args) > 0 {
			content, err = readFile(args[0])
		} else {
			content, err = readStdin()
		}

		if err != nil {
			return err
		}

		content = strings.ReplaceAll(content, altDelimiter, delimiter)
		slides := strings.Split(content, delimiter)

		user, err := user.Current()
		if err != nil {
			return errors.New("could not get current user")
		}

		m, exists := meta.New().ParseHeader(slides[0])
		// If the user specifies a custom configuration options
		// skip the first "slide" since this is all configuration
		if exists {
			slides = slides[1:]
		}

		p := tea.NewProgram(model.Model{
			Slides: slides,
			Page:   0,
			Author: user.Name,
			Date:   time.Now().Format("2006-01-02"),
			Theme:  styles.SelectTheme(m.Theme),
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

func readFile(path string) (string, error) {
	s, err := os.Stat(path)
	if err != nil {
		return "", errors.New("could not read file")
	}
	if s.IsDir() {
		return "", errors.New("can not read directory")
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return "", errors.New("no slides provided")
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}
