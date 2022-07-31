package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindToggleText is a ast.NodeKind of the ToggleText node.
var KindToggleText = ast.NewNodeKind("ToggleText")

// A ToggleText represents a toggle text in Notion.
type ToggleText struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *ToggleText) Kind() ast.NodeKind { return KindToggleText }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *ToggleText) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
