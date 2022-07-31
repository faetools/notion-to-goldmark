package goldmark

import (
	"fmt"

	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/samber/lo"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

const mergeAnnotations = false

// TODO we might have to work with empty spaces:
// - trim them and add back in without formatting

func toNodeRichTexts(rts notion.RichTexts) (ns []ast.Node) {
	// construct the initial wrappers
	var ws annotationWrappers = lo.Map(rts, func(rt notion.RichText, _ int) *annotationWrapper {
		return newAnnotationWrapper(rt)
	})

	// merge siblings that have the same annotation
	if mergeAnnotations {
		ws = ws.mergeSiblings()
	}

	// transform each wrapper into a node
	return lo.Map(ws, func(w *annotationWrapper, _ int) ast.Node { return w.toNode() })
}

type annotationWrappers []*annotationWrapper

type annotationWrapper struct {
	// represents first all annotations
	// then objects will be wrapped so in the end it's just one annotation
	ann notion.Annotations

	// a node that is wrapped
	// should only be present if there are no annotations or sub wrappers
	node ast.Node

	// the objects that this object wraps
	subs annotationWrappers
}

func newAnnotationWrapper(t notion.RichText) *annotationWrapper {
	// NOTE: we are ignoring t.Href for now because it seems to be just
	// a repeat of t.Text.Link.Url, t.Mention.LinkPreview.Url etc.

	wr := &annotationWrapper{ann: t.Annotations}

	switch t.Type {
	case notion.RichTextTypeText:
		wr.node = toNodeText(t.Text)
	case notion.RichTextTypeEquation:
		wr.node = toNodeEquation(t.Equation)
	case notion.RichTextTypeMention:
		wr.node = &n_ast.Mention{Content: t.Mention}
	default:
		panic(fmt.Sprintf("invalid RichText of type %q", t.Type))
	}

	return wr
}

// mergeSiblings merges siblings if they have the same annotations
func (ws annotationWrappers) mergeSiblings() annotationWrappers {
	for i, l := 0, len(ws)-1; i < l; i++ {
		if merger := mergeIfSameAnnotation(ws[i], ws[i+1]); merger != nil {
			ws[i] = merger // write the merger

			// remove i+1
			copy(ws[i+1:], ws[i+2:])
			ws = ws[:len(ws)-1]

			return ws.mergeSiblings()
		}
	}

	return ws
}

// mergeIfSameAnnotation merges two wrappers if they have the same annotation
func mergeIfSameAnnotation(w1, w2 *annotationWrapper) *annotationWrapper {
	ann := notion.Annotations{Color: notion.ColorDefault}

	for _, v := range []struct {
		this, next, merged *bool
	}{
		{&w1.ann.Bold, &w2.ann.Bold, &ann.Bold},
		{&w1.ann.Underline, &w2.ann.Underline, &ann.Underline},
		{&w1.ann.Italic, &w2.ann.Italic, &ann.Italic},
		{&w1.ann.Code, &w2.ann.Code, &ann.Code},
		{&w1.ann.Strikethrough, &w2.ann.Strikethrough, &ann.Strikethrough},
	} {
		if *v.this && *v.next {
			// set the flag for the new wrapper
			*v.merged = true

			// disable the flag for the old wrappers
			*v.this = false
			*v.next = false

			return &annotationWrapper{
				ann:  ann,
				subs: merge(w1, w2),
			}
		}
	}

	switch w1.ann.Color {
	case notion.ColorDefault, "": // no color, ignore
	case w2.ann.Color: // same color
		ann.Color = w1.ann.Color
		w1.ann.Color = notion.ColorDefault
		w2.ann.Color = notion.ColorDefault

		return &annotationWrapper{
			ann:  ann,
			subs: merge(w1, w2),
		}
	}

	return nil
}

// merge merges two wrappers so they can form the subwrappers to a new wrapper
func merge(w1, w2 *annotationWrapper) annotationWrappers {
	// make sure all wrappers contain either node or annotation
	subs := make(annotationWrappers, 0, 2)

	for _, w := range []*annotationWrapper{w1, w2} {
		if w.node != nil {
			// w has a node
			subs = append(subs, w)
			continue
		}

		// w is just a wrapper
		subs = append(subs, w.subs...)
	}

	return subs
}

func (w *annotationWrapper) toNode() ast.Node {
	if w.node != nil {
		return wrapInAnnotation(w.ann, w.node)
	}

	return wrapInAnnotation(w.ann, lo.Map(w.subs, func(sub *annotationWrapper, _ int) ast.Node {
		return sub.toNode()
	})...)
}

func wrapInAnnotation(a notion.Annotations, children ...ast.Node) ast.Node {
	var w ast.Node

	switch {
	case a.Bold:
		a.Bold = false
		w = ast.NewEmphasis(2)
	case a.Underline:
		a.Underline = false
		w = &n_ast.Underline{}
	case a.Italic:
		a.Italic = false
		w = ast.NewEmphasis(1)
	case a.Code:
		a.Code = false
		w = ast.NewCodeSpan()
	case a.Strikethrough:
		a.Strikethrough = false
		w = extast.NewStrikethrough()
	case a.Color != notion.ColorDefault && a.Color != "":
		c := a.Color
		a.Color = notion.ColorDefault
		w = &n_ast.Color{Color: c}
	default:
		return children[0]
	}

	for _, child := range children {
		w.AppendChild(w, child)
	}

	return wrapInAnnotation(a, w)
}
