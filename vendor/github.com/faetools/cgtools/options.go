package cgtools

import (
	"os"
	"time"
)

// defaultMode is the default mode to be used when writing.
const defaultMode = 0o644

type options struct {
	perm       os.FileMode
	modTime    time.Time
	skipFormat bool
}

func defaultOptions() *options {
	return &options{perm: defaultMode}
}

// Option sets an option.
type Option func(*options)

// FileMode overrides the default file mode.
func FileMode(perm os.FileMode) Option {
	return func(o *options) { o.perm = perm }
}

// ModTime manually sets the mod time after changing a file.
func ModTime(modTime time.Time) Option {
	return func(o *options) { o.modTime = modTime }
}

// SkipFormat skips formatting.
func SkipFormat(o *options) { o.skipFormat = true }

func getOptions(opts []Option) *options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return o
}
