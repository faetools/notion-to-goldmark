package ast

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

// KindCheckboxText is a ast.NodeKind of the CheckboxText node.
var KindCheckboxText = ast.NewNodeKind("CheckboxText")

// A CheckboxText represents a callout text in Notion.
type CheckboxText struct {
	ast.BaseInline
	Checked bool
}

// Kind returns a kind of this node.
func (n *CheckboxText) Kind() ast.NodeKind { return KindCheckboxText }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *CheckboxText) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Checked": fmt.Sprintf("%t", n.Checked),
	}, nil)
}
