package ast

import "github.com/yuin/goldmark/ast"

// KindCaption is a ast.NodeKind of the Caption node.
var KindCaption = ast.NewNodeKind("Caption")

// A Caption represents a caption in Notion.
type Caption struct{ ast.BaseInline }

// Kind returns a kind of this node.
func (n *Caption) Kind() ast.NodeKind { return KindCaption }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Caption) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
