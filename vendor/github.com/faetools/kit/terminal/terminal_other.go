// +build !linux,!windows

package terminal

import (
	"fmt"
	"os"
)

// Clear clears the terminal.
func Clear() { fmt.Fprint(os.Stdout, "\033[H\033[2J") }
