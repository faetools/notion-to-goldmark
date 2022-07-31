package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindUser is a ast.NodeKind of the User node.
var KindUser = ast.NewNodeKind("User")

// A User represents a Notion user.
type User struct {
	ast.BaseInline
	Data notion.User
}

// Kind returns a kind of this node.
func (n *User) Kind() ast.NodeKind { return KindUser }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *User) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
