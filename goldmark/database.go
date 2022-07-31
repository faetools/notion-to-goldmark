package goldmark

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/faetools/go-notion/pkg/notion"
	n_ast "github.com/faetools/notion-to-goldmark/ast"
	"github.com/samber/lo"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	extast "github.com/yuin/goldmark/extension/ast"
)

var propertyValueTypes = map[notion.PropertyType][]byte{
	notion.PropertyTypeCheckbox:       []byte("typesCheckbox"),
	notion.PropertyTypeCreatedBy:      []byte("typesCreatedBy"),
	notion.PropertyTypeLastEditedBy:   []byte("typesCreatedBy"),
	notion.PropertyTypeCreatedTime:    []byte("typesCreatedAt"),
	notion.PropertyTypeDate:           []byte("typesDate"),
	notion.PropertyTypeEmail:          []byte("typesEmail"),
	notion.PropertyTypeFiles:          []byte("typesFile"),
	notion.PropertyTypeFormula:        []byte("typesFormula"),
	notion.PropertyTypeLastEditedTime: []byte("typesLastEditedTime"),
	notion.PropertyTypeMultiSelect:    []byte("typesMultipleSelect"),
	notion.PropertyTypeNumber:         []byte("typesNumber"),
	notion.PropertyTypePeople:         []byte("typesPerson"),
	notion.PropertyTypePhoneNumber:    []byte("typesPhoneNumber"),
	notion.PropertyTypeRelation:       []byte("typesRelation"),
	notion.PropertyTypeRichText:       []byte("typesText"),
	notion.PropertyTypeRollup:         []byte("typesRollup"),
	notion.PropertyTypeSelect:         []byte("typesSelect"),
	notion.PropertyTypeStatus:         []byte("typesStatus"),
	notion.PropertyTypeTitle:          []byte("typesTitle"),
	notion.PropertyTypeUrl:            []byte("typesUrl"),
}

type tableCollector struct {
	p        *pageCollector
	root     string
	props    notion.PropertyMetaMap
	propKeys []string
}

func (p *pageCollector) getTable(id notion.Id) (*extast.Table, error) {
	db, err := p.cli.GetNotionDatabase(p.ctx, id)
	if err != nil {
		return nil, err
	}

	keys := lo.Keys(db.Properties)

	// the title first, the rest are sorted alphabetically
	sort.SliceStable(keys, func(i, j int) bool {
		if db.Properties[keys[j]].Type == notion.PropertyTypeTitle {
			return false
		}

		return db.Properties[keys[i]].Type == notion.PropertyTypeTitle || keys[i] < keys[j]
	})

	c := &tableCollector{
		p:        p,
		root:     getDir(db.Title.Content(), db.Id),
		props:    db.Properties,
		propKeys: keys,
	}

	table := extast.NewTable()
	setClasses(table, "", classCollectionContent)

	table.AppendChild(table, c.tableHeader())

	entries, err := p.cli.GetAllDatabaseEntries(p.ctx, id)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(entries, func(i, j int) bool {
		t1 := entries[i].Title()
		t2 := entries[j].Title()

		if t1 == "" {
			return false
		}

		return t2 == "" || t1 < t2
	})

	for _, entry := range entries {
		row, err := c.tableRow(entry)
		if err != nil {
			return nil, err
		}

		table.AppendChild(table, row)
	}

	return table, nil
}

func (c *tableCollector) tableHeader() *extast.TableHeader {
	headerRow := extast.NewTableRow(nil)

	for _, name := range c.propKeys {
		prop := c.props[name]

		cell := extast.NewTableCell()
		headerRow.AppendChild(headerRow, cell)

		icon := &n_ast.PropertyIcon{}
		cell.AppendChild(cell, icon)
		setClasses(icon, "", classIcon, classPropertyIcon)

		svg := &n_ast.SVG{}
		icon.AppendChild(icon, svg)
		svg.SetAttributeString(attrViewBox, getViewBox(prop.Type))
		svg.SetAttributeString(attrStyle, styleSVG)
		svg.SetAttributeString(attrClass, propertyValueTypes[prop.Type])
		svg.AppendChild(svg, getSVGContent(prop.Type))

		cell.AppendChild(cell, newString(name))
	}

	return extast.NewTableHeader(headerRow)
}

func (c *tableCollector) tableRow(entry notion.Page) (*extast.TableRow, error) {
	row := extast.NewTableRow(nil)
	row.SetAttributeString(attrID, []byte(entry.Id))

	for _, propName := range c.propKeys {
		prop := entry.Properties[propName]

		cell := extast.NewTableCell()

		id, err := url.QueryUnescape(prop.Id)
		if err != nil {
			return nil, err
		}

		cell.SetAttributeString(attrClass, []byte(fmt.Sprintf("cell-%s", id)))

		cellContent, err := c.toNodesPropertyValue(entry, propName, prop)
		if err != nil {
			return nil, err
		}

		for _, content := range cellContent {
			cell.AppendChild(cell, content)
		}

		row.AppendChild(row, cell)
	}

	return row, nil
}

var numberFormats = map[notion.NumberConfigFormat]string{
	notion.NumberConfigFormatNumber: "%v",
	notion.NumberConfigFormatEuro:   "€%.2f",

	// TODO test these
	notion.NumberConfigFormatBaht:             "฿%.2f",
	notion.NumberConfigFormatCanadianDollar:   "$%.2f",
	notion.NumberConfigFormatChileanPeso:      "$%.2f",
	notion.NumberConfigFormatColombianPeso:    "$%.2f",
	notion.NumberConfigFormatDanishKrone:      "kr%.2f",
	notion.NumberConfigFormatDirham:           "د.إ%.2f",
	notion.NumberConfigFormatDollar:           "$%.2f",
	notion.NumberConfigFormatForint:           "Ft%.2f",
	notion.NumberConfigFormatFranc:            "CHF%.2f",
	notion.NumberConfigFormatHongKongDollar:   "$%.2f",
	notion.NumberConfigFormatKoruna:           "Kč%.2f",
	notion.NumberConfigFormatKrona:            "kr%.2f",
	notion.NumberConfigFormatLeu:              "lei%.2f",
	notion.NumberConfigFormatLira:             "₺%.2f",
	notion.NumberConfigFormatMexicanPeso:      "$%.2f",
	notion.NumberConfigFormatNewTaiwanDollar:  "NT$%.2f",
	notion.NumberConfigFormatNewZealandDollar: "$%.2f",
	notion.NumberConfigFormatNorwegianKrone:   "kr%.2f",
	notion.NumberConfigFormatNumberWithCommas: "%v",
	notion.NumberConfigFormatPercent:          "%v%%",
	notion.NumberConfigFormatPhilippinePeso:   "₱%.2f",
	notion.NumberConfigFormatPound:            "£%.2f",
	notion.NumberConfigFormatRand:             "R%.2f",
	notion.NumberConfigFormatReal:             "R$%.2f",
	notion.NumberConfigFormatRinggit:          "RM%.2f",
	notion.NumberConfigFormatRiyal:            "﷼%.2f",
	notion.NumberConfigFormatRuble:            "₽%.2f",
	notion.NumberConfigFormatRupee:            "₨%.2f",
	notion.NumberConfigFormatRupiah:           "Rp%.2f",
	notion.NumberConfigFormatShekel:           "₪%.2f",
	notion.NumberConfigFormatWon:              "₩%.2f",
	notion.NumberConfigFormatYen:              "¥%.2f",
	notion.NumberConfigFormatYuan:             "¥%.2f",
	notion.NumberConfigFormatZloty:            "zł%.2f",
}

func (c *tableCollector) toNodesPropertyValue(p notion.Page, propName string, prop notion.PropertyValue) ([]ast.Node, error) {
	var n ast.Node

	switch prop.Type {
	case notion.PropertyTypeTitle:
		n = linkToPage(prop.Title.Content(), p.Id, c.p.root, c.root)
	case notion.PropertyTypeNumber:
		if prop.Number == nil {
			return nil, nil
		}

		format := c.props[propName].Number.Format
		n = newString(fmt.Sprintf(numberFormats[format], prop.GetNumber()))
	case notion.PropertyTypeRelation:
		ids := prop.GetRelation().GetIDs()

		l := len(ids)

		if l == 0 {
			return nil, nil
		}

		nodes := make([]ast.Node, 2*l-1)

		for i, id := range ids {
			p, err := c.p.cli.GetNotionPage(c.p.ctx, notion.Id(id))
			if err != nil {
				return nil, err
			}

			if i != 0 {
				nodes[i*2-1] = newString(", ")
			}

			nodes[i*2] = linkToPage(p.Title(), id, c.p.root, c.root)
		}

		return nodes, nil
	case notion.PropertyTypeCheckbox:
		n = newCheckbox(prop.GetCheckbox(), "")
	case notion.PropertyTypePhoneNumber:
		if prop.PhoneNumber == nil {
			return nil, nil
		}

		n = newURLValue(*prop.PhoneNumber)
	case notion.PropertyTypeRichText:
		return toNodeRichTexts(prop.GetRichText()), nil
	case notion.PropertyTypeFormula:
		switch prop.Formula.Type {
		case notion.FormulaTypeBoolean:
			n = newCheckbox(*prop.Formula.Boolean, "")
		case notion.FormulaTypeDate:
			if prop.Formula.Date == nil {
				return nil, nil
			}

			n = n_ast.NewDate(prop.Formula.Date, true)
		case notion.FormulaTypeNumber:
			if prop.Formula.Number == nil {
				return nil, nil
			}

			n = newString(fmt.Sprintf("%v", *prop.Formula.Number))

		case notion.FormulaTypeString:
			if prop.Formula.String == nil {
				return nil, nil
			}

			n = newString(*prop.Formula.String)
		default:
			return nil, fmt.Errorf("unknown formula type %q", prop.Formula.Type)
		}
	case notion.PropertyTypeStatus:
		status := prop.Status
		if prop.Status == nil {
			status = customStatus
		}

		n = &n_ast.Status{Data: status}
	case notion.PropertyTypeEmail:
		if prop.Email == nil {
			return nil, nil
		}

		n = newURLValue(*prop.Email)
	case notion.PropertyTypeRollup:
		switch prop.Rollup.Type {
		case notion.RollupTypeArray:
			if prop.Rollup.Array == nil || len(*prop.Rollup.Array) == 0 {
				return nil, nil
			}

			switch c.props[propName].Rollup.Function {
			case "show_unique": // Show unique values
				// HACK just the number of items and don't print 1
				// we don't know why notion does not print 1
				if l := len(*prop.Rollup.Array); l > 1 {
					return []ast.Node{newString(fmt.Sprintf("%d", len(*prop.Rollup.Array)))}, nil
				}

				return nil, nil
			case "unique": // Count unique values
			case "median": // Median
			}

			nodes := []ast.Node{}

			for _, el := range *prop.Rollup.Array {
				switch el.Type {
				case notion.RollupArrayItemTypeTitle:
					if el.Title == nil {
						continue
					}

					nodes = append(nodes, toNodeRichTexts(*el.Title)...)
				case notion.RollupArrayItemTypeDate:
					if el.Date != nil {
						nodes = append(nodes, &n_ast.Date{Date: el.Date})
					}
				default:
					nodes = append(nodes, newString(fmt.Sprintf("UNIMPLEMENTED rollup %s", el.Type)))
				}
			}

			return nodes, nil

		// case notion.RollupTypeDate: TODO
		case notion.RollupTypeNumber:
			if prop.Rollup.Number == nil || *prop.Rollup.Number == 0 {
				return nil, nil
			}

			n = newString(fmt.Sprintf("%v", *prop.Rollup.Number))
		// case notion.RollupTypeString: // TODO
		default:
			n = newString(fmt.Sprintf("UNIMPLEMENTED rollup %s", prop.Rollup.Type))
		}
	case notion.PropertyTypeFiles:
		return lo.Map(prop.GetFiles(), func(f notion.File, _ int) ast.Node {
			n := &n_ast.FileInCell{}

			rawURL := f.URL()

			u, _ := url.Parse(rawURL) // TODO don't underscore error

			fileName := filepath.Base(u.Path)

			if u.Host == "s3.us-west-2.amazonaws.com" &&
				strings.HasPrefix(u.Path, "/secure.notion-static.com/") {
				// TODO download

				rawURL = filepath.Join(c.p.root, c.root, getDir(p.Title(), p.Id), fileName)
			}

			link := newLink("", rawURL)

			switch filepath.Ext(fileName) {
			case ".jpg", ".jpeg", ".png":
				img := ast.NewImage(link)
				img.SetAttributeString(attrStyle, []byte("width:20px;max-height:24px"))

				// TODO remove in favor of below code
				if f.Type == notion.FileTypeExternal {
					img.SetAttributeString("alt", []byte(rawURL))
				} else {
					img.SetAttributeString("alt", []byte(fileName))
				}

				txt := ast.NewText()
				txt.Segment = text.NewSegment(0, 0) // TODO
				img.AppendChild(img, txt)

				link.AppendChild(link, img)
			default:
				link.AppendChild(link, newString(fileName))
			}

			n.AppendChild(n, link)

			return n
		}), nil

	case notion.PropertyTypeUrl:
		if prop.Url == nil {
			return nil, nil
		}

		n = newURLValue(*prop.Url)
	case notion.PropertyTypePeople:
		if prop.People == nil {
			return nil, nil
		}

		return lo.Map(*prop.People, func(u notion.User, _ int) ast.Node {
			return &n_ast.User{Data: u}
		}), nil
	case notion.PropertyTypeSelect:
		if prop.Select == nil {
			return nil, nil
		}

		n = &n_ast.Select{Data: prop.Select}
	case notion.PropertyTypeMultiSelect:
		if prop.MultiSelect == nil {
			return nil, nil
		}

		return lo.Map(*prop.MultiSelect, func(sel notion.SelectValue, _ int) ast.Node {
			return &n_ast.Select{Data: &sel}
		}), nil

	case notion.PropertyTypeDate:
		if prop.Date == nil {
			return nil, nil
		}

		n = n_ast.NewDate(prop.Date, false)
	case notion.PropertyTypeCreatedTime:
		n = n_ast.NewDate(&notion.Date{Start: *prop.CreatedTime}, true)
	case notion.PropertyTypeCreatedBy:
		n = &n_ast.User{Data: *prop.CreatedBy}
	case notion.PropertyTypeLastEditedBy:
		n = &n_ast.User{Data: *prop.LastEditedBy}
	default:
		n = newString(fmt.Sprintf("UNIMPLEMENTED %s", prop.Type))
	}

	return []ast.Node{n}, nil
}

func newURLValue(dest string) *ast.Link {
	n := ast.NewLink()
	n.Destination = []byte(dest)

	n.SetAttributeString(attrClass, classURLValue)
	n.AppendChild(n, newString(dest))

	return n
}

// HACK
var customStatus = &notion.SelectValue{
	Name:  "I&#x27;ve taken a look",
	Color: notion.ColorBlue,
}
