package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindLinkPreview is a ast.NodeKind of the LinkPreview node.
var KindLinkPreview = ast.NewNodeKind("LinkPreview")

// A LinkPreview represents a link preview in Notion.
type LinkPreview struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *LinkPreview) Kind() ast.NodeKind { return KindLinkPreview }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *LinkPreview) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
