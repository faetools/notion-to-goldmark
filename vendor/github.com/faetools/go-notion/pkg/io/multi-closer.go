package io

import "io"

type multiCloser struct {
	closers []io.Closer
}

func MultiCloser(writers ...io.Closer) io.Closer {
	allClosers := make([]io.Closer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiCloser); ok {
			allClosers = append(allClosers, mw.closers...)
		} else {
			allClosers = append(allClosers, w)
		}
	}
	return &multiCloser{allClosers}
}

func (t *multiCloser) Close() error {
	for _, w := range t.closers {
		if err := w.Close(); err != nil {
			return err
		}
	}
	return nil
}
