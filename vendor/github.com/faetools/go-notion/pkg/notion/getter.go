package notion

import "context"

var _ Getter = (*Client)(nil)

// Getter is any client that can get notion documents.
type Getter interface {
	// GetNotionPage return the notion page or an error.
	GetNotionPage(ctx context.Context, id Id) (*Page, error)
	// GetAllBlocks returns all blocks of a given page or block.
	GetAllBlocks(ctx context.Context, id Id) (Blocks, error)
	// GetNotionDatabase returns the notion database or an error.
	GetNotionDatabase(ctx context.Context, id Id) (*Database, error)
	// GetAllDatabaseEntries returns all database entries or an error.
	GetAllDatabaseEntries(ctx context.Context, id Id) (Pages, error)
}
