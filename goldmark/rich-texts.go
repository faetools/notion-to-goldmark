package goldmark

import (
	"fmt"

	n_ast "github.com/faetools/notion-to-goldmark/ast"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

// toNodeRichTextWithoutAnnotations returns the node without annotation information
func toNodeRichTextWithoutAnnotations(t notion.RichText) ast.Node {
	// NOTE: we are ignoring t.Href for now because it seems to be just
	// a repeat of t.Text.Link.Url, t.Mention.LinkPreview.Url etc.

	switch t.Type {
	case notion.RichTextTypeText:
		return toNodeText(t.Text)
	case notion.RichTextTypeEquation:
		return toNodeEquation(t.Equation)
	case notion.RichTextTypeMention:
		return &n_ast.Mention{Content: t.Mention}
	default:
		panic(fmt.Sprintf("invalid RichText of type %q", t.Type))
	}
}

// we want to wrap the nodes
// - n1,n2,n3->n1,n2,n3
// - n1(bold,italic),n2(bold),n3(bold,italic)->bold{italic{n1},n2,italic{n3}}
// - n1(bold,italic),n2(italic),n3(bold,italic)->italic{bold{n1},n2,bold{n3}}
// - n1(bold),n2(bold,italic),n3(italic)->bold{n1,italic{n2}},italic{n3}

// TODO we might have to work with empty spaces:
// - trim them and add back in without formatting

// annotations with extra methods
type annotations notion.Annotations

// node returns the first node representing an annotation
// this should only be called on annotations that only have annotation
func (a annotations) node() ast.Node {
	switch {
	case a.Bold:
		return ast.NewEmphasis(2)
	case a.Italic:
		return ast.NewEmphasis(1)
	case a.Underline:
		return ast.NewEmphasis(3)
	case a.Strikethrough:
		return extast.NewStrikethrough()
	case a.Code:
		return ast.NewCodeSpan()
	case a.Color != notion.ColorDefault && a.Color != "":
		return &n_ast.Color{Color: a.Color}
	default:
		return nil
	}
}

type annotationWrappers []*annotationWrapper

type annotationWrapper struct {
	// represents first all annotations
	// then objects will be wrapped so in the end it's just one annotation
	ann annotations

	// a node that is wrapped
	// should only be present if there are no annotations or sub wrappers
	node ast.Node

	// the objects that this object wraps
	subs annotationWrappers
}

func toNodeRichTexts(rts notion.RichTexts) (ns []ast.Node) {
	// construct the initial wrappers
	ws := make(annotationWrappers, len(rts))

	for i, rt := range rts {
		n := toNodeRichTextWithoutAnnotations(rt)
		ws[i] = &annotationWrapper{node: n, ann: annotations(rt.Annotations)}
	}

	// merge siblings that have the same annotation
	ws = ws.mergeSiblings()

	// transform each wrapper into a node
	return ws.toNodes()
}

// mergeIfSameAnnotation merges two wrappers if they have the same annotation
func mergeIfSameAnnotation(w1, w2 *annotationWrapper) *annotationWrapper {
	ann := annotations{Color: notion.ColorDefault}

	for _, v := range []struct {
		this, next, merged *bool
	}{
		{&w1.ann.Bold, &w2.ann.Bold, &ann.Bold},
		{&w1.ann.Italic, &w2.ann.Italic, &ann.Italic},
		{&w1.ann.Code, &w2.ann.Code, &ann.Code},
		{&w1.ann.Strikethrough, &w2.ann.Strikethrough, &ann.Strikethrough},
		{&w1.ann.Underline, &w2.ann.Underline, &ann.Underline},
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
		if w.ann.node() != nil || w.node != nil {
			// w has annotation information
			subs = append(subs, w)
			continue
		}

		// w is just a wrapper
		subs = append(subs, w.subs...)
	}

	return subs
}

// remove removes an element from the slice or array.
func remove[T any](collection []T, i int) []T {
	copy(collection[i:], collection[i+1:])

	end := len(collection) - 1
	var zero T
	collection[end] = zero

	return collection[:end]
}

// mergeSiblings merges siblings if they have the same annotations
func (ws annotationWrappers) mergeSiblings() annotationWrappers {
	for i, l := 0, len(ws)-1; i < l; i++ {
		if merger := mergeIfSameAnnotation(ws[i], ws[i+1]); merger != nil {
			ws[i] = merger       // write the merger
			ws = remove(ws, i+1) // remove i+1

			return ws.mergeSiblings()
		}
	}

	return ws
}

func (ws annotationWrappers) toNodes() []ast.Node {
	out := make([]ast.Node, len(ws))

	for i, wr := range ws {
		out[i] = wr.toNode()
	}

	return out
}

func (a annotations) wrap(n ast.Node) ast.Node {
	var w ast.Node

	switch {
	case a.Bold:
		a.Bold = false
		w = ast.NewEmphasis(2)
	case a.Italic:
		a.Italic = false
		w = ast.NewEmphasis(1)
	case a.Underline:
		a.Underline = false
		w = ast.NewEmphasis(3)
	case a.Strikethrough:
		a.Strikethrough = false
		w = extast.NewStrikethrough()
	case a.Code:
		a.Code = false
		w = ast.NewCodeSpan()
	case a.Color != notion.ColorDefault && a.Color != "":
		a.Color = notion.ColorDefault
		w = &n_ast.Color{Color: a.Color}
	default:
		return n
	}

	w.AppendChild(w, n)

	return a.wrap(w) // call recursively
}

func (w *annotationWrapper) toNode() ast.Node {
	if w.node != nil {
		return w.ann.wrap(w.node)
	}

	n := w.ann.node()

	for _, sub := range w.subs {
		n.AppendChild(n, sub.toNode())
	}

	return n
}
