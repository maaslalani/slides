//go:build windows
// +build windows

package coral

import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

var preExecHookFn = preExecHook

func preExecHook(c *Command) {
	if MousetrapHelpText != "" && mousetrap.StartedByExplorer() {
		c.Print(MousetrapHelpText)
		if MousetrapDisplayDuration > 0 {
			time.Sleep(MousetrapDisplayDuration)
		} else {
			c.Println("Press return to continue...")
			fmt.Scanln()
		}
		os.Exit(1)
	}
}
