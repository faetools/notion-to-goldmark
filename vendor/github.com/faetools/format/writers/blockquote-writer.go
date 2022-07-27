package writers

import (
	"bytes"
	"io"
)

// NewBlockquoteWriter returns a new writer that adds '>' each empty line or lines starting with '>' and '> ' on other lines.
func NewBlockquoteWriter(base io.Writer) Writer {
	bqw := &blockquoteWriter{base: Upgrade(base)}
	bqw.Writer = NewStateWriter(bqw.stateStart)

	return bqw
}

type blockquoteWriter struct {
	Writer        // The underlying state writer.
	base   Writer // The base writer.
}

func (w *blockquoteWriter) stateStart(p []byte) (StateFunc, []byte, int, error) {
	// Start with writing '>'.
	return w.stateAfterGreaterThan, p, 0, w.base.WriteByte(greaterThan)
}

func (w *blockquoteWriter) stateNewLine(p []byte) (StateFunc, []byte, int, error) {
	if err := w.base.WriteByte(newLine); err != nil {
		return w.stateNewLine, p, 0, err
	}

	return w.stateAfterGreaterThan, p[1:], 1, w.base.WriteByte(greaterThan)
}

func (w *blockquoteWriter) stateAfterGreaterThan(p []byte) (StateFunc, []byte, int, error) {
	var err error

	switch p[0] {
	case newLine:
		return w.stateNewLine, p, 0, nil
	case greaterThan:
	default:
		// Only write space if next token is not greaterThan.
		err = w.base.WriteByte(emptySpace)
	}

	return w.stateText, p, 0, err
}

func (w *blockquoteWriter) stateText(p []byte) (StateFunc, []byte, int, error) {
	// Write normally until new line.
	indexNewLine := bytes.IndexByte(p, newLine)
	if indexNewLine == -1 {
		// Write normally.
		size, err := w.base.Write(p)
		return w.stateText, p[size:], size, err
	}

	// Write normally until k.
	size, err := w.base.Write(p[:indexNewLine])

	return w.stateNewLine, p[indexNewLine:], size, err
}
