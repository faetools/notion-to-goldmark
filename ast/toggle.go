package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindToggle is a ast.NodeKind of the Toggle node.
var KindToggle = ast.NewNodeKind("Toggle")

// A Toggle represents a toggle in Notion.
type Toggle struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Toggle) Kind() ast.NodeKind { return KindToggle }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Toggle) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
