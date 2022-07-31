package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindChildDatabase is a ast.NodeKind of the ChildDatabase node.
var KindChildDatabase = ast.NewNodeKind("ChildDatabase")

// A ChildDatabase represents a child database in Notion.
type ChildDatabase struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *ChildDatabase) Kind() ast.NodeKind { return KindChildDatabase }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *ChildDatabase) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
