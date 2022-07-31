package goldmark

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

const (
	blockColorPrefix = "block-color-"

	attrClass   = "class"
	attrID      = "id"
	attrType    = "type"
	attrHref    = "href"
	attrViewBox = "viewBox"
	attrStyle   = "style"
	attrD       = "d"
	attrPoints  = "points"
)

var (
	attrSep = []byte{' '}

	classBlockColorTeal        = []byte(blockColorPrefix + "teal")
	classToDoList              = []byte("to-do-list")
	classCheckbox              = []byte("checkbox")
	classCheckboxOn            = []byte("checkbox-on")
	classCheckboxOff           = []byte("checkbox-off")
	classCheckboxTextChecked   = []byte("to-do-children-checked")
	classCheckboxTextUnchecked = []byte("to-do-children-unchecked")
	classBulletedList          = []byte("bulleted-list")
	classNumberedList          = []byte("numbered-list")
	classToggle                = []byte("toggle")
	classCode                  = []byte("code")
	classLinkToPage            = []byte("link-to-page")
	classCollectionContent     = []byte("collection-content")
	classCollectionTitle       = []byte("collection-title")
	classIcon                  = []byte("icon")
	classPropertyIcon          = []byte("property-icon")
	classURLValue              = []byte("url-value")

	viewBoxStandard = []byte("0 0 14 14")
	viewBoxStatus   = []byte("0 0 16 16")
	viewBoxRollup   = []byte("0 0 18 18")

	styleSVG = []byte("width:14px;height:14px;display:block;fill:rgba(55, 53, 47, 0.45);flex-shrink:0;-webkit-backface-visibility:hidden")
)

func setClasses(n ast.Node, c notion.Color, classes ...[]byte) {
	switch c {
	case notion.ColorDefault, "": // ignore
	case notion.ColorGreen:
		// a green background color is actually teal
		classes = append([][]byte{classBlockColorTeal}, classes...)
	default:
		classes = append([][]byte{[]byte(blockColorPrefix + c)}, classes...)
	}

	n.SetAttributeString(attrClass, bytes.Join(classes, attrSep))
}

func appendRichTexts(n ast.Node, txts notion.RichTexts) {
	for _, txt := range toNodeRichTexts(txts) {
		n.AppendChild(n, txt)
	}
}

func getDir(name string, id notion.UUID) string {
	return string(
		util.URLEscape(
			[]byte(
				fmt.Sprintf("%s %s", name, strings.ReplaceAll(string(id), "-", "")),
			),
			true))
}
