package cmd

import (
	"github.com/muesli/coral"
)

// RemoteCmd is the command for remote controlling a slides session.
// It exposes the slides control to external processes.
var RemoteCmd = &coral.Command{
	Use:     "remote",
	Aliases: []string{"remote"},
	Short:   "Remote control slides session",
	Args:    coral.NoArgs,
}

func init() {
	RemoteCmd.AddCommand(RemoteSocketCmd)
}
