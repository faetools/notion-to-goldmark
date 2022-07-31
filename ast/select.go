package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindSelect is a ast.NodeKind of the Select node.
var KindSelect = ast.NewNodeKind("Select")

// A Select represents a select in Notion.
type Select struct {
	ast.BaseInline
	Data *notion.SelectValue
}

// Kind returns a kind of this node.
func (n *Select) Kind() ast.NodeKind { return KindSelect }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Select) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
