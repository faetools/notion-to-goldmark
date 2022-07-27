package cgtools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/faetools/format"
	"github.com/faetools/format/yaml"
	"github.com/faetools/kit/terminal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// NewOsGenerator returns a new generator at the current folder.
func NewOsGenerator() *Generator { return NewGenerator(afero.NewOsFs()) }

// NewGenerator returns a new generator for any filesystem.
func NewGenerator(fs afero.Fs) *Generator { return &Generator{fs: fs} }

// Generator can generate files.
type Generator struct{ fs afero.Fs }

// Write reads from the reader and writes the content to the file.
func (g Generator) Write(path string, r io.Reader, opts ...Option) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading from the reader: %w", err)
	}

	return g.WriteBytes(path, bytes, opts...)
}

// WriteBytes writes the bytes to the given path, creating any directories and overwriting existing files.
func (g Generator) WriteBytes(path string, content []byte, opts ...Option) (err error) {
	o := getOptions(opts)

	if !o.skipFormat {
		// Format contents according to type.
		content, err = format.Format(path, content)
		if err != nil {
			return fmt.Errorf("formatting %s: %w", path, err)
		}
	}

	return g.writeBytes(path, content, o)
}

func (g Generator) writeBytes(path string, content []byte, o *options) error {
	// Check if we need to write.
	current, readErr := afero.ReadFile(g.fs, path)

	if bytes.Equal(content, current) {
		if viper.GetBool("verbose") {
			terminal.Printf(aurora.Green, "  • %v is unchanged\n", path)
		}

		return nil
	}

	newFile := errors.Is(readErr, os.ErrNotExist)
	if newFile {
		if err := g.MkdirAll(filepath.Dir(path)); err != nil {
			return fmt.Errorf("making directories %s: %w", filepath.Dir(path), err)
		}
	}

	f, err := g.fs.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, o.perm)
	if err != nil {
		return fmt.Errorf("opening or creating file %s: %w", path, err)
	}
	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return fmt.Errorf("writing to file %s: %w", path, err)
	}

	if newFile {
		terminal.Printf(aurora.Green, "  • %v generated\n", path)
	} else {
		terminal.Printf(aurora.Green, "  • %v regenerated\n", path)
	}

	if o.modTime.IsZero() {
		return nil
	}

	return g.fs.Chtimes(path, o.modTime, o.modTime)
}

// MkdirAll creates all folders.
func (g *Generator) MkdirAll(dir string) error {
	return g.fs.MkdirAll(dir, os.ModePerm)
}

// WriteTemplate writes the template to the given path, creating any directories and overwriting existing files.
func (g Generator) WriteTemplate(path string, tpl *template.Template, data interface{}) error {
	b := &bytes.Buffer{}
	if err := tpl.Execute(b, data); err != nil {
		return fmt.Errorf("executing template %s: %w", tpl.Name(), err)
	}

	return g.WriteBytes(path, b.Bytes())
}

// WriteYAML writes an interface as a yaml file.
func (g Generator) WriteYAML(path string, v interface{}, opts ...yaml.EncodeOption) error {
	b, err := yaml.Encode(v, opts...)
	if err != nil {
		return fmt.Errorf("encoding %T to YAML: %w", v, err)
	}

	return g.writeBytes(path, b, defaultOptions())
}

// WriteJSON writes an interface as a JSON file.
func (g Generator) WriteJSON(path string, v interface{}) error {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encoding %T to JSON: %w", v, err)
	}

	return g.writeBytes(path, b.Bytes(), defaultOptions())
}
