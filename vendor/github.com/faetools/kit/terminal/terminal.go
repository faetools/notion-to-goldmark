package terminal

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// Println prints something in color or after having gone through some other auora transformation.
func Println(transform func(interface{}) aurora.Value, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprintln(a...)))
}

// Print prints something in color or after having gone through some other aurora transformation.
func Print(transform func(interface{}) aurora.Value, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprint(a...)))
}

// Printf prints something in color or after having gone through some other aurora transformation.
func Printf(transform func(interface{}) aurora.Value, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprintf(format, a...)))
}
