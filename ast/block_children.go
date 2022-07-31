package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindBlockChildren is a ast.NodeKind of the BlockChildren node.
var KindBlockChildren = ast.NewNodeKind("BlockChildren")

// A BlockChildren represents a block children in Notion.
type BlockChildren struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *BlockChildren) Kind() ast.NodeKind { return KindBlockChildren }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *BlockChildren) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
