package docs

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/faetools/go-notion/pkg/notion"
)

// Type defines the type of document.
type Type string

// Defines values for Type.
const (
	TypePage            Type = "page"
	TypeBlocks          Type = "blocks"
	TypeDatabase        Type = "database"
	TypeDatabaseEntries Type = "database_entries"
)

// Skip is used as a return value from a Visitor to indicate that
// the page or database named in the call is to be skipped.
// It is not returned as an error by any other function.
var Skip = errors.New("skip this page or database") //nolint:go-lint

// Walk traverses notion documents.
func Walk(ctx context.Context, v Visitor, tp Type, id notion.Id) error {
	switch tp {
	case TypePage:
		if err := v.VisitPage(ctx, id); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return fmt.Errorf("visiting page %q: %w", id, err)
		}

		return Walk(ctx, v, TypeBlocks, id)
	case TypeBlocks:
		blocks, err := v.VisitBlocks(ctx, id)
		if err != nil {
			return fmt.Errorf("visiting block children of %q: %w", id, err)
		}

		for _, b := range blocks {
			switch b.Type {
			case notion.BlockTypeChildPage:
				if err := Walk(ctx, v, TypePage, notion.Id(b.Id)); err != nil {
					return fmt.Errorf("walking block child page of %q: %w", id, err)
				}
			case notion.BlockTypeChildDatabase:
				if err := v.VisitDatabase(ctx, notion.Id(b.Id)); err != nil {
					if errors.Is(err, Skip) {
						continue
					}

					// Unfortunately, notion does not tell us if a child database
					// has the same ID as the block ID or if a child database was merely referenced.
					//
					// We're still calling Walk, and check here for existence of a database with the ID.
					apiErr := &notion.Error{}
					if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound {
						continue
					}

					return fmt.Errorf("visiting block child database of %q: %w", id, err)
				}

				if err := Walk(ctx, v, TypeDatabaseEntries, notion.Id(b.Id)); err != nil {
					return fmt.Errorf("walking child database entries of %q: %w", id, err)
				}
			default:
				if b.HasChildren {
					if err := Walk(ctx, v, TypeBlocks, notion.Id(b.Id)); err != nil {
						return fmt.Errorf("walking children a block in %q: %w", id, err)
					}
				}
			}
		}

		return nil
	case TypeDatabase:
		if err := v.VisitDatabase(ctx, id); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return fmt.Errorf("visiting database %q: %w", id, err)
		}

		return Walk(ctx, v, TypeDatabaseEntries, id)
	case TypeDatabaseEntries:
		entries, err := v.VisitDatabaseEntries(ctx, id)
		if err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return fmt.Errorf("visiting database entries of %q: %w", id, err)
		}

		for _, p := range entries {
			if err := Walk(ctx, v, TypePage, notion.Id(p.Id)); err != nil {
				return fmt.Errorf("walking entry in database %q: %w", id, err)
			}
		}

		return nil
	default:
		return fmt.Errorf("unknown object type %q", tp)
	}
}
