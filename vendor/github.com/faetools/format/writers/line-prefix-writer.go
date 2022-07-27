package writers

import (
	"bytes"
	"io"
)

type linePrefixWriter struct {
	Writer // The underlying state writer.

	base   Writer
	prefix []byte
}

// NewLinePrefixWriter returns a writer that adds a prefix to each non-empty line.
func NewLinePrefixWriter(base io.Writer, prefix []byte) Writer {
	if len(prefix) == 0 {
		return Upgrade(base)
	}

	pw := &linePrefixWriter{base: Upgrade(base), prefix: prefix}
	pw.Writer = NewStateWriter(pw.stateStart)

	return pw
}

func (w *linePrefixWriter) stateStart(p []byte) (StateFunc, []byte, int, error) {
	if p[0] == newLine {
		return w.stateNewLine, p, 0, nil
	}

	// Start with writing the prefix.
	_, err := w.base.Write(w.prefix)

	return w.stateText, p, 0, err
}

func (w *linePrefixWriter) stateText(p []byte) (StateFunc, []byte, int, error) {
	k := bytes.IndexByte(p, newLine)
	if k == -1 {
		// Write normally.
		size, err := w.base.Write(p)
		return w.stateText, p[size:], size, err
	}

	// Write normally until new line.
	size, err := w.base.Write(p[:k])

	return w.stateNewLine, p[k:], size, err
}

func (w *linePrefixWriter) stateNewLine(p []byte) (StateFunc, []byte, int, error) {
	size, err := w.base.Write(arrayNewLine)
	return w.stateStart, p[size:], size, err
}
