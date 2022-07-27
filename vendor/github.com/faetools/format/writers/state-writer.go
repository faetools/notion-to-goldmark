package writers

type (
	// StateFunc is a function that writes and returns another state for the state writer.
	StateFunc func([]byte) (StateFunc, []byte, int, error)

	stateWriter struct{ state StateFunc }
)

var _ Writer = (*stateWriter)(nil)

// NewStateWriter returns a new writer that changes its state based on the given state function.
func NewStateWriter(start StateFunc) Writer {
	return &stateWriter{state: start}
}

func (w *stateWriter) Write(p []byte) (size int, err error) {
	var written int

	for {
		if len(p) == 0 {
			return
		}

		w.state, p, written, err = w.state(p)
		size += written

		if err != nil {
			return
		}
	}
}

func (w *stateWriter) WriteByte(b byte) (err error) {
	var p []byte
	w.state, p, _, err = w.state([]byte{b})

	switch {
	case err != nil:
		return err
	case len(p) > 0:
		// State has changed, do again.
		return w.WriteByte(p[0])
	default:
		return nil
	}
}

func (w *stateWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}
