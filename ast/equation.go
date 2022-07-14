package ast

import "github.com/yuin/goldmark/ast"

// KindEquation is a ast.NodeKind of the Equation node.
var KindEquation = ast.NewNodeKind("Equation")

// A Equation represents a equation in Notion.
type Equation struct {
	ast.BaseInline
	Expression string
}

// Kind returns a kind of this node.
func (n *Equation) Kind() ast.NodeKind { return KindEquation }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Equation) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Expression": n.Expression}, nil)
}
