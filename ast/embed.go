package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindEmbed is a ast.NodeKind of the Embed node.
var KindEmbed = ast.NewNodeKind("Embed")

// A Embed represents an embedding in Notion.
type Embed struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Embed) Kind() ast.NodeKind { return KindEmbed }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Embed) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
