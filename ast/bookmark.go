package ast

import "github.com/yuin/goldmark/ast"

// KindBookmark is a ast.NodeKind of the Bookmark node.
var KindBookmark = ast.NewNodeKind("Bookmark")

// A Bookmark represents a bookmark in Notion.
type Bookmark struct {
	ast.BaseInline
	URL string
}

// Kind returns a kind of this node.
func (n *Bookmark) Kind() ast.NodeKind { return KindBookmark }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Bookmark) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"URL": n.URL}, nil)
}
