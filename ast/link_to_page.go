package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindLinkToPage is a ast.NodeKind of the LinkToPage node.
var KindLinkToPage = ast.NewNodeKind("LinkToPage")

// A LinkToPage represents a link to a page or database in Notion.
type LinkToPage struct {
	ast.BaseInline
	Content notion.LinkToPage
}

// NewLinkToPage returns a new link to page node.
func NewLinkToPage(l notion.LinkToPage) ast.Node {
	return &LinkToPage{Content: l}
}

// Kind returns a kind of this node.
func (n *LinkToPage) Kind() ast.NodeKind { return KindLinkToPage }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *LinkToPage) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Type": string(n.Content.Type),
		"ID":   string(n.Content.ID()),
	}, nil)
}
