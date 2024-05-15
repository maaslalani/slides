package cmd

import (
	"os"

	"github.com/maaslalani/slides/internal/remote"
	"github.com/muesli/coral"
)

var (
	socketFile string
)

// RemoteSocketCmd is the command for remote controlling a slides session
// using UNIX socket that is being listened by the slides session.
var RemoteSocketCmd = &coral.Command{
	Use:     "socket [flags] command [args]",
	Aliases: []string{"remote"},
	Short:   "Remote Control using listening socket",
	Args:    coral.ArbitraryArgs,
	RunE: func(cmd *coral.Command, args []string) error {
		k := os.Getenv("SLIDES_REMOTE_SOCKET")
		if k != "" {
			socketFile = k
		}
		return nil
	},
}

func init() {
	RemoteSocketCmd.PersistentFlags().StringVar(
		&socketFile, "socketFile", remote.SocketRemoteListenerDefaultPath, "Socket File")
	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "slide-next",
			Short: "Go to the next slide",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.SlideNext()
			},
		},
	)
	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "slide-prev",
			Short: "Go to the previous slide",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.SlidePrevious()
			},
		},
	)
	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "slide-first",
			Short: "Go to the first slide",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.SlideFirst()
			},
		},
	)

	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "slide-last",
			Short: "Go to the last slide",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.SlideLast()
			},
		},
	)

	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "code-exec",
			Short: "Execute Code blocks of current slide in session",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.CodeExec()
			},
		},
	)

	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "code-copy",
			Short: "Execute Code blocks of current slide in session",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.CodeCopy()
			},
		},
	)

	RemoteSocketCmd.AddCommand(
		&coral.Command{
			Use:   "quit",
			Short: "Quit the slides session",
			RunE: func(cmd *coral.Command, args []string) error {
				remote, err := remote.NewSocketRemote(socketFile)
				if err != nil {
					return err
				}
				defer remote.Close()
				return remote.Quit()
			},
		},
	)
}
