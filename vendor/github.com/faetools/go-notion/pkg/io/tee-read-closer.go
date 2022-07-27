package io

import "io"

type readCloser struct {
	io.Reader
	io.Closer
}

// TeeReadCloser returns a ReadCloser that writes to w what it reads from rc.
// When the ReadCloser is closed, rc is closed.
// If the writer is also a closer, it will then also be closed.
func TeeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	c, ok := w.(io.Closer)
	if ok {
		c = MultiCloser(rc, c)
	} else {
		c = rc
	}

	return &readCloser{
		Reader: io.TeeReader(rc, w),
		Closer: c,
	}
}
