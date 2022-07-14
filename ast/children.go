package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindChildren is a ast.NodeKind of the Children node.
var KindChildren = ast.NewNodeKind("Children")

// A Children represents children in Notion.
// Usually, children will have an additional indentation.
type Children struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *Children) Kind() ast.NodeKind { return KindChildren }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Children) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
