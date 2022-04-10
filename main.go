package main

import (
	_ "embed"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/model"
	"github.com/maaslalani/slides/internal/navigation"
	"github.com/maaslalani/slides/internal/server"
)

func printError(err error) {
	fmt.Fprintf(os.Stderr, `Error: %s
Usage:
  slides <file.md>
  slides <file.md> -server
`, err.Error())
}

func main() {
	var err error
	var fileName string

	useServer := false

	if len(os.Args) > 1 {
		fileName = os.Args[1]

		if len(os.Args) > 2 && os.Args[2] == "-server" {
			useServer = true
		}
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

	if useServer {
		startServer(presentation)
	} else {
		startTUI(presentation)
	}
}

func startTUI(presentation model.Model) {
	var err error
	p := tea.NewProgram(presentation, tea.WithAltScreen())
	err = p.Start()
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

func startServer(presentation model.Model) {
	var err error
	var host = ""
	var keyPath = "slides"
	var port = 53531

	k := os.Getenv("SLIDES_SERVER_KEY_PATH")
	if k != "" {
	    keyPath = k
	}
	h := os.Getenv("SLIDES_SERVER_HOST")
	if h != "" {
	    host = h
	}
	p := os.Getenv("SLIDES_SERVER_PORT")
	if p != "" {
	    port, _ = strconv.Atoi(p)
	}

	s, err := server.NewServer(keyPath, host, port, presentation)
	if err != nil {
	    log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting Slides server on %s:%d", host, port)
	go func() {
	    if err = s.Start(); err != nil {
	        log.Fatalln(err)
	    }
	}()

	<-done
	log.Print("Stopping Slides server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
	    log.Fatalln(err)
	}
}
