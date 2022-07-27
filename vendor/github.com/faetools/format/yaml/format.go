package yaml

import (
	"bytes"
	"strings"

	"github.com/goccy/go-yaml/token"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	tagMerge = "!!merge"
	tagStr   = "!!str"
)

// Format formats the yaml file.
func Format(src []byte) ([]byte, error) {
	// Skip empty files.
	if len(src) == 0 {
		return src, nil
	}

	n := &yaml.Node{}
	if err := yaml.Unmarshal(src, n); err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	formatNode(n)

	b := &bytes.Buffer{}
	enc := yaml.NewEncoder(b)
	enc.SetIndent(1)

	err := enc.Encode(n)

	return b.Bytes(), err
}

func isNeedQuoted(v string) bool {
	return !strings.ContainsAny(v, "\n\r") && token.IsNeedQuoted(v)
}

func formatNode(n *yaml.Node) {
	switch n.Tag {
	case tagMerge:
		n.Tag = "" // Delete this tag.
	case tagStr:
		if isNeedQuoted(n.Value) {
			n.Style = yaml.SingleQuotedStyle
		} else {
			n.Style = yaml.FlowStyle
		}
	}

	for _, el := range n.Content {
		formatNode(el)
	}
}
