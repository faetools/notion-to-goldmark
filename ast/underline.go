package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindUnderline is a ast.NodeKind of the Underline node.
var KindUnderline = ast.NewNodeKind("Underline")

// A Underline represents an underline in Notion.
type Underline struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Underline) Kind() ast.NodeKind { return KindUnderline }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Underline) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
