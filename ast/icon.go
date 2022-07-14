package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindIcon is a ast.NodeKind of the Icon node.
var KindIcon = ast.NewNodeKind("Icon")

// A Icon represents a icon in Notion.
type Icon struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Icon) Kind() ast.NodeKind { return KindIcon }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Icon) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
