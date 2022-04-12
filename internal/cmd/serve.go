package cmd

import (
    "context"
    "log"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/maaslalani/slides/internal/model"
    "github.com/maaslalani/slides/internal/navigation"
    "github.com/maaslalani/slides/internal/server"
    "github.com/muesli/coral"
)

var (
    host string
    port int
    keyPath string
    err error
    fileName string

    ServeCmd = &coral.Command{
        Use:     "serve",
        Aliases: []string{"server"},
        Short:   "Start an SSH server to run slides",
        Args:    coral.ArbitraryArgs,
        RunE: func(cmd *coral.Command, args []string) error {
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

            s, err := server.NewServer(keyPath, host, port, presentation)
            if err != nil {
                return err
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
                return err
            }
            return err
        },
    }
)

func init() {
    ServeCmd.Flags().StringVar(&keyPath, "keyPath", "slides", "Server private key path")
    ServeCmd.Flags().StringVar(&host, "host", "localhost", "Server host to bind to")
    ServeCmd.Flags().IntVar(&port, "port", 53531, "Server port to bind to")
}
