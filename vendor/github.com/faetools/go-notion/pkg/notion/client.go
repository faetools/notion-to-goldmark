package notion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/faetools/client"
	"github.com/google/uuid"
)

const (
	versionHeader  = "Notion-Version"
	version        = "2022-02-22"
	maxPageSizeInt = 100
)

var maxPageSize PageSize = maxPageSizeInt

// NewDefaultClient returns a new client with the default options.
func NewDefaultClient(bearer string, opts ...client.Option) (*Client, error) {
	opts = append([]client.Option{
		client.WithBearer(bearer),
		client.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set(versionHeader, version)
			return nil
		}),
	}, opts...)

	return NewClient(opts...)
}

// Error ensures responses with an error fulfill the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%d %s: %s - %s", e.Status, http.StatusText(e.Status), e.Code, e.Message)
}

// GetNotionPage return the notion page or an error.
func (c Client) GetNotionPage(ctx context.Context, id Id) (*Page, error) {
	resp, err := c.GetPage(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// UpdateNotionPage updates the notion page or returns an error.
func (c Client) UpdateNotionPage(ctx context.Context, p Page) (*Page, error) {
	resp, err := c.UpdatePage(ctx, Id(p.Id), UpdatePageJSONRequestBody(p))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// GetNotionDatabase returns the notion database or an error.
func (c Client) GetNotionDatabase(ctx context.Context, id Id) (*Database, error) {
	resp, err := c.GetDatabase(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// GetAllDatabaseEntries returns all database entries or an error.
func (c Client) GetAllDatabaseEntries(ctx context.Context, id Id) (Pages, error) {
	return c.GetDatabaseEntries(ctx, id, nil, nil)
}

// GetDatabaseEntries return filtered and sorted database entries or an error.
func (c Client) GetDatabaseEntries(ctx context.Context, id Id, filter *Filter, sorts *Sorts) (Pages, error) {
	entries := Pages{}

	var cursor *UUID
	for {
		resp, err := c.QueryDatabase(ctx, id,
			QueryDatabaseJSONRequestBody{
				Filter:      filter,
				PageSize:    maxPageSizeInt,
				Sorts:       sorts,
				StartCursor: cursor,
			})
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		case http.StatusTooManyRequests:
			return nil, resp.JSON429
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		entries = append(entries, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return entries, nil
		}

		cursor = (*UUID)(resp.JSON200.NextCursor)
	}
}

func ensureDatabaseIsValid(db *Database) {
	// set mandatory values
	db.Object = "database"
	if db.Parent != nil && db.Parent.Type == "" {
		db.Parent.Type = "page_id"
	}

	// initialize properties
	if db.Properties == nil {
		db.Properties = PropertyMetaMap{"Title": TitleProperty}
		return
	}

	// make sure a title property is present
	for _, prop := range db.Properties {
		if prop.Title != nil {
			return
		}
	}

	db.Properties["Title"] = TitleProperty
}

// CreateNotionDatabase creates a notion database or returns an error.
func (c Client) CreateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	ensureDatabaseIsValid(&db)

	// create a UUID for the new database
	db.Id = UUID(uuid.NewString())

	resp, err := c.CreateDatabase(ctx, CreateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// UpdateNotionDatabase updates a notion database or returns an error.
func (c Client) UpdateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	// can't be present when updating
	db.Parent = nil
	db.CreatedTime = nil

	ensureDatabaseIsValid(&db)

	resp, err := c.UpdateDatabase(ctx, Id(db.Id), UpdateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// ListAllUsers returns all users in the workspace.
func (c Client) ListAllUsers(ctx context.Context) (Users, error) {
	users := Users{}

	var cursor *StartCursor
	for {
		resp, err := c.ListUsers(ctx, &ListUsersParams{
			PageSize:    &maxPageSize,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		case http.StatusTooManyRequests:
			return nil, resp.JSON429
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		users = append(users, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return users, nil
		}

		cursor = (*StartCursor)(resp.JSON200.NextCursor)
	}
}

// GetAllBlocks returns all blocks of a given page or block.
func (c Client) GetAllBlocks(ctx context.Context, id Id) (Blocks, error) {
	blocks := Blocks{}

	var cursor *StartCursor
	for {
		resp, err := c.GetBlocks(ctx, id, &GetBlocksParams{
			PageSize:    &maxPageSize,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("getting blocks for %s: %w", id, err)
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		blocks = append(blocks, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return blocks, nil
		}

		cursor = (*StartCursor)(resp.JSON200.NextCursor)
	}
}
