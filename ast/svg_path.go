package ast

import "github.com/yuin/goldmark/ast"

// KindSVGPath is a ast.NodeKind of the SVGPath node.
var KindSVGPath = ast.NewNodeKind("SVGPath")

// A SVGPath represents an svg path in Notion.
type SVGPath struct {
	ast.BaseInline
}

// Kind returns a kind of this node.
func (n *SVGPath) Kind() ast.NodeKind { return KindSVGPath }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *SVGPath) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
