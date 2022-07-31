package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindStatus is a ast.NodeKind of the Status node.
var KindStatus = ast.NewNodeKind("Status")

// A Status represents a status in Notion.
type Status struct {
	ast.BaseInline
	Data *notion.SelectValue
}

// Kind returns a kind of this node.
func (n *Status) Kind() ast.NodeKind { return KindStatus }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Status) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
