package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/internal/navigation"
)

func printError(err error) {
	fmt.Fprintf(os.Stderr, `Error: %s
Usage:
  slides <file.md>

`, err.Error())
}

func main() {
	var err error
	var fileName string

	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	presentation := model.Model{
		Page:     0,
		Date:     time.Now().Format("2006-01-02"),
		FileName: fileName,
		Search:   navigation.NewSearch(),
	}
	err = presentation.Load()
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	p := tea.NewProgram(presentation, tea.WithAltScreen())
	err = p.Start()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}
