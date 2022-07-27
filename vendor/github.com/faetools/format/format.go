package format

import (
	"bytes"
	"path/filepath"

	"github.com/faetools/format/golang"
	"github.com/faetools/format/markdown"
	"github.com/faetools/format/yaml"
	"github.com/faetools/kit/terminal"
	"github.com/logrusorgru/aurora"
	dockerfile "github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	json "github.com/tidwall/pretty"
)

var min = minify.New()

func init() {
	min.AddFunc(".css", css.Minify)
	min.AddFunc(".js", js.Minify)
	min.AddFunc(".html", html.Minify)
}

// Format formats contents according to type.
func Format(path string, src []byte) (out []byte, err error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		src, err = golang.Format(path, src)
	case ".yml", ".yaml":
		src, err = yaml.Format(src)
	case ".md":
		src, err = markdown.Format(src)
	case ".json":
		src = json.PrettyOptions(src, &json.Options{Indent: "  "})
	case "":
		if filepath.Base(path) == "Dockerfile" {
			docker, err := dockerfile.Parse(bytes.NewReader(src))
			if err != nil {
				return nil, errors.Wrap(err, "parsing dockerfile")
			}

			if len(docker.Warnings) > 0 {
				terminal.Println(aurora.Red, "Dockerfile warnings:")
				for _, warning := range docker.Warnings {
					terminal.Println(aurora.Red, "  • ", warning)
				}
			}
		}
	case ".html", ".js", ".css":
		src, err = min.Bytes(ext, src)
	}

	return src, errors.Wrapf(err, "formatting %s", path)
}
