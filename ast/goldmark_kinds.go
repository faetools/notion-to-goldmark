package ast

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/yuin/goldmark/ast"
// )

// var (
// 	// GoldmarkKindColor is a ast.NodeKind of the GoldmarkColor node.
// 	GoldmarkKindColor = ast.NewNodeKind("Color")

// 	// GoldmarkKindEquation is a ast.NodeKind of the GoldmarkEquation node.
// 	GoldmarkKindEquation = ast.NewNodeKind("Equation")

// 	// GoldmarkKindLinkPreview is a ast.NodeKind of the GoldmarkLinkPreview node.
// 	GoldmarkKindLinkPreview = ast.NewNodeKind("LinkPreview")

// 	// GoldmarkKindUserPerson is a ast.NodeKind of the GoldmarkUserPerson node.
// 	GoldmarkKindUserPerson = ast.NewNodeKind("UserPerson")

// 	// GoldmarkKindUserBot is a ast.NodeKind of the GoldmarkUserBot node.
// 	GoldmarkKindUserBot = ast.NewNodeKind("UserBot")

// 	// GoldmarkKindExternalFile is a ast.NodeKind of the GoldmarkExternalFile node.
// 	GoldmarkKindExternalFile = ast.NewNodeKind("ExternalFile")

// 	// GoldmarkKindNotionFile is a ast.NodeKind of the GoldmarkNotionFile node.
// 	// GoldmarkKindNotionFile = ast.NewNodeKind("NotionFile")

// 	// GoldmarkKindEmbed is a ast.NodeKind of the GoldmarkKindEmbed node.
// 	GoldmarkKindEmbed = ast.NewNodeKind("Embed")
// )

// type (
// 	// A GoldmarkColor struct represents a Notion color.
// 	GoldmarkColor struct {
// 		ast.BaseInline
// 		Color
// 	}

// 	// A GoldmarkEquation struct represents an equation in Notion.
// 	GoldmarkEquation struct{ ast.BaseInline }

// 	// A GoldmarkLinkPreview struct represents a link preview in Notion.
// 	GoldmarkLinkPreview struct {
// 		ast.BaseInline
// 		URL string
// 	}

// 	// A GoldmarkUserPerson struct represents a user in Notion.
// 	GoldmarkUserPerson struct {
// 		ast.BaseInline

// 		// A unique identifier of user.
// 		ID UUID

// 		// User's name, as displayed in Notion.
// 		Name string

// 		// Chosen avatar image.
// 		AvatarURL string

// 		Person
// 	}

// 	// A GoldmarkUserBot struct represents a bot in Notion.
// 	GoldmarkUserBot struct {
// 		ast.BaseInline

// 		// A unique identifier of user.
// 		ID UUID

// 		// User's name, as displayed in Notion.
// 		Name string

// 		// Chosen avatar image.
// 		AvatarURL string

// 		Bot
// 	}

// 	// A GoldmarkExternalFile struct represents an external file in Notion.
// 	GoldmarkExternalFile struct {
// 		ast.BaseInline
// 		URL string
// 	}

// 	// A GoldmarkNotionFile struct represents a file hosted by Notion.
// 	// GoldmarkNotionFile struct {
// 	// 	ast.BaseInline
// 	// 	NotionFile
// 	// }

// 	// A GoldmarkEmbed struct represents an external file in Notion.
// 	GoldmarkEmbed struct {
// 		ast.BaseInline
// 		URL string
// 	}
// )

// // Kind implements Node.Kind.
// func (n *GoldmarkColor) Kind() ast.NodeKind { return GoldmarkKindColor }

// // Kind implements Node.Kind.
// func (n *GoldmarkEquation) Kind() ast.NodeKind { return GoldmarkKindEquation }

// // Kind implements Node.Kind.
// func (n *GoldmarkLinkPreview) Kind() ast.NodeKind { return GoldmarkKindLinkPreview }

// // Kind implements Node.Kind.
// func (n *GoldmarkUserPerson) Kind() ast.NodeKind { return GoldmarkKindUserPerson }

// // Kind implements Node.Kind.
// func (n *GoldmarkUserBot) Kind() ast.NodeKind { return GoldmarkKindUserBot }

// // Kind implements Node.Kind.
// func (n *GoldmarkExternalFile) Kind() ast.NodeKind { return GoldmarkKindExternalFile }

// // Kind implements Node.Kind.
// // func (n *GoldmarkNotionFile) Kind() ast.NodeKind { return GoldmarkKindNotionFile }

// // Dump implements Node.Dump.
// func (n *GoldmarkColor) Dump(_ []byte, level int) {
// 	fmt.Printf("%sGoldmarkColor: %q\n", strings.Repeat("    ", level), n.Color)
// }

// // Dump implements Node.Dump.
// func (n *GoldmarkEquation) Dump(source []byte, level int) {
// 	ast.DumpHelper(n, source, level, nil, nil)
// }

// // Dump implements Node.Dump.
// func (n *GoldmarkLinkPreview) Dump(source []byte, level int) {
// 	fmt.Printf("%sGoldmarkLinkPreview: %q\n", strings.Repeat("    ", level), n.URL)
// }

// // Dump implements Node.Dump.
// func (n *GoldmarkUserPerson) Dump(source []byte, level int) {
// 	fmt.Printf("%sGoldmarkUserPerson: %#v\n", strings.Repeat("    ", level), n.ID)
// }

// // Dump implements Node.Dump.
// func (n *GoldmarkUserBot) Dump(source []byte, level int) {
// 	fmt.Printf("%sGoldmarkUserBot: %#v\n", strings.Repeat("    ", level), n.ID)
// }

// // Dump implements Node.Dump.
// func (n *GoldmarkExternalFile) Dump(source []byte, level int) {
// 	fmt.Printf("%sGoldmarkExternalFile: %q\n", strings.Repeat("    ", level), n.URL)
// }

// // Dump implements Node.Dump.
// // func (n *GoldmarkNotionFile) Dump(source []byte, level int) {
// // 	fmt.Printf("%sGoldmarkNotionFile: %q (expires: %s)\n", strings.Repeat("    ", level), n.Url, n.ExpiryTime)
// // }

// // Kind implements Node.Kind.
// func (n *GoldmarkEmbed) Kind() ast.NodeKind { return GoldmarkKindEmbed }

// // Dump implements Node.Dump.
// func (n *GoldmarkEmbed) Dump(source []byte, level int) {
// 	fmt.Printf("%sGoldmarkEmbed: %q\n", strings.Repeat("    ", level), n.URL)
// }
