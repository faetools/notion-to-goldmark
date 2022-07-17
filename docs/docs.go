package docs

import (
	"context"
	"errors"
	"fmt"

	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/notion"
)

// Defines values for Type.
const (
	TypePage     = docs.TypePage
	TypeBlocks   = docs.TypeBlocks
	TypeDatabase = docs.TypeDatabase
)

var Skip = docs.Skip

type GetterVisitor interface {
	notion.Getter

	// VisitPage defines what is to be done when visiting a page.
	VisitPage(p *notion.Page) error

	// VisitBlock defines what is to be done when visiting a block.
	VisitBlock(block notion.Block) error

	// DatabaseVisit defines what is to be done when visiting a database.
	VisitDatabase(db *notion.Database) error
}

func Walk(ctx context.Context, v GetterVisitor, tp docs.Type, id notion.Id) error {
	switch tp {
	case TypePage:
		p, err := v.GetNotionPage(ctx, id)
		if err != nil {
			return err
		}

		if err := v.VisitPage(p); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		return Walk(ctx, v, TypeBlocks, id)
	case TypeBlocks:
		blocks, err := v.GetAllBlocks(ctx, id)
		if err != nil {
			return err
		}

		for _, b := range blocks {
			if err := v.VisitBlock(b); err != nil {
				if errors.Is(err, Skip) {
					continue
				}

				return err
			}

			switch b.Type {
			case notion.BlockTypeChildPage:
				if err := Walk(ctx, v, TypePage, notion.Id(b.Id)); err != nil {
					return err
				}
			case notion.BlockTypeChildDatabase:
				// Unfortunately, notion does not tell us if this child database
				// has the same ID as the block ID or if a child database was just referenced.
				//
				// We're still calling Walk, the user will need to filter out such references
				// in their VisitDatabase method.
				if err := Walk(ctx, v, TypeDatabase, notion.Id(b.Id)); err != nil {
					return err
				}
			default:
				if b.HasChildren {
					if err := Walk(ctx, v, TypeBlocks, notion.Id(b.Id)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	case TypeDatabase:
		db, err := v.GetNotionDatabase(ctx, id)
		if err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		if err := v.VisitDatabase(db); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		entries, err := v.GetAllDatabaseEntries(ctx, id)
		if err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		for _, p := range entries {
			if err := Walk(ctx, v, TypePage, notion.Id(p.Id)); err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("unknown object type %q", tp)
	}
}
