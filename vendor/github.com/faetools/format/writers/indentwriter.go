package writers

import (
	"bytes"
	"io"
)

// NewIndentWriter creates a new indent writer that indents non-empty lines with indent number of tabs.
func NewIndentWriter(base io.Writer, indent int) Writer {
	return NewLinePrefixWriter(base, bytes.Repeat(arrayTab, indent))
}
