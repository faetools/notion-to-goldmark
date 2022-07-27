package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"sync"

	"github.com/faetools/client"
)

// TODO move to "github.com/faetools/clueo" (see client in latin -> cliens -> alteration of cluēns, present active participle of clueō)

// consider voco: I call, summon, beckon; I invoke, call upon

// FSClient returns responses using files in a filesystem.
// It keeps track of which files were seen.
type FSClient struct {
	fs fs.FS

	mu   sync.Mutex
	seen map[string]bool

	onNotExists func(path string) any
}

// request transform client -> remove query parameters, for example
// clueo.Client -> http.Client

// NewFSClient returns a new client that returns responses using files in a filesystem.
// The other arguments define what kind of body should be returned when a file does not exist.
func NewFSClient(files fs.FS,
	// TODO options
	// - on not exists
	// - on read error
	// - add file extension
	onNotExists func(path string) any,
) (*FSClient, error) {
	c := &FSClient{
		fs:          files,
		seen:        map[string]bool{},
		onNotExists: onNotExists,
	}

	return c, fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		c.seen[path] = false

		return err
	})
}

var header = map[string][]string{
	client.ContentType: {client.MIMEApplicationJSON},
}

func (c *FSClient) readFile(path string) ([]byte, error) {
	f, err := c.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c.mu.Lock()
	c.seen[path] = true // TODO separate client for this
	c.mu.Unlock()

	return io.ReadAll(f)
}

// Do implements client.HTTPRequestDoer.
func (c *FSClient) Do(req *http.Request) (*http.Response, error) {
	path := req.URL.Path[1:] + ".json"

	f, err := c.fs.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return c.respond404(req)
		}

		return nil, err
	}
	defer f.Close()

	c.mu.Lock()
	c.seen[path] = true
	c.mu.Unlock()

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return c.response(http.StatusOK, body, req), nil
}

func (c *FSClient) respond404(req *http.Request) (*http.Response, error) {
	if c.onNotExists == nil {
		return c.response(http.StatusNotFound, nil, req), nil
	}

	body := c.onNotExists(req.URL.Path)

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling %T body: %w", body, err)
	}

	return c.response(http.StatusNotFound, b, req), nil
}

func (c *FSClient) response(status int, body []byte, req *http.Request) *http.Response {
	return &http.Response{
		StatusCode:    status,
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(body)),
		Request:       req,
		ContentLength: int64(len(body)),
	}
}

// Unseen returns all files that were not seen; useful for testing.
func (c *FSClient) Unseen() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	untested := []string{}
	for path, ok := range c.seen {
		if !ok {
			untested = append(untested, path)
		}
	}

	return untested
}
