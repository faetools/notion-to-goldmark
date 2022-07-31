package ast

import (
	"github.com/yuin/goldmark/ast"
)

// KindEmbedSource is a ast.NodeKind of the EmbedSource node.
var KindEmbedSource = ast.NewNodeKind("EmbedSource")

// A EmbedSource represents an embed source in Notion.
type EmbedSource struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *EmbedSource) Kind() ast.NodeKind { return KindEmbedSource }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *EmbedSource) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
