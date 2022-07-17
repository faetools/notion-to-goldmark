package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindCallout is a ast.NodeKind of the Callout node.
var KindCallout = ast.NewNodeKind("Callout")

// A Callout represents a callout in Notion.
type Callout struct {
	ast.BaseInline
}

// NewCallout returns a new callout with the icon as the first child.
func NewCallout(icon notion.Icon) *Callout {
	n := &Callout{}
	switch icon.Type {
	case notion.IconTypeEmoji:
		n.AppendChild(n, ast.NewString([]byte(*icon.Emoji)))
	case notion.IconTypeExternal:
		link := ast.NewLink()

		link.Destination = []byte(icon.URL())
		n.AppendChild(n, ast.NewImage(link))
	}

	return n
}

// Kind returns a kind of this node.
func (n *Callout) Kind() ast.NodeKind { return KindCallout }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Callout) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
