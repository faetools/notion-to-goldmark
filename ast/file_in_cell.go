package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindFileInCell is a ast.NodeKind of the FileInCell node.
var KindFileInCell = ast.NewNodeKind("FileInCell")

// A FileInCell represents a file in a cell in in Notion.
type FileInCell struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *FileInCell) Kind() ast.NodeKind { return KindFileInCell }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *FileInCell) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
