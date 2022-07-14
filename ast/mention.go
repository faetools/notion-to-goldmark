package ast

import (
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
)

// KindMention is a ast.NodeKind of the Mention node.
var KindMention = ast.NewNodeKind("Mention")

// A Mention represents a mention in Notion.
type Mention struct {
	ast.BaseInline
	Content *notion.Mention
}

// Kind returns a kind of this node.
func (n *Mention) Kind() ast.NodeKind { return KindMention }

// Dump dumps an AST tree structure to stdout.
// This function completely aimed for debugging.
func (n *Mention) Dump(source []byte, level int) {
	m := n.Content
	kv := map[string]string{"Type": string(m.Type)}

	switch m.Type {
	case notion.MentionTypeDatabase:
		kv["Database ID"] = string(m.Database.Id)
	case notion.MentionTypeDate:
		kv["Date"] = m.Date.String()
	case notion.MentionTypeLinkPreview:
		kv["Preview URL"] = m.LinkPreview.Url
	case notion.MentionTypePage:
		kv["Page ID"] = string(m.Page.Id)
	case notion.MentionTypeUser:
		kv["User"] = m.User.Name
	}

	ast.DumpHelper(n, source, level, kv, nil)
}
