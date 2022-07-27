package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/faetools/cgtools"
	"github.com/faetools/client"
	_io "github.com/faetools/go-notion/pkg/io"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

// A FSClientWriter not only uses a filesystem to get responses,
// it also writes the responses to that same filesystem.
type FSClientWriter struct {
	// a client to get responses to requests
	cli client.HTTPRequestDoer
	// a generator to write down the responses in files
	gen *cgtools.Generator
	// the already existing files
	files client.HTTPRequestDoer
	// a wait group for writing files in parallel
	eg *errgroup.Group

	fsClient *FSClient
}

// NewFSClientWriter returns a new FSClientWriter.
func NewFSClientWriter(cli client.HTTPRequestDoer, fs afero.Fs) (*FSClientWriter, error) {
	fsClient, err := NewFSClient(afero.NewIOFS(fs), nil)
	if err != nil {
		return nil, err
	}

	return &FSClientWriter{
		cli: cli,

		// to generate any files we don't have
		gen: cgtools.NewGenerator(fs),

		// cache the responses from the fsClient
		// because the generator could be in the middle of writing
		files: NewCachingClient(fsClient),

		eg:       &errgroup.Group{},
		fsClient: fsClient,
	}, nil
}

// Do fulfils the HTTPRequestDoer interface.
func (c *FSClientWriter) Do(req *http.Request) (*http.Response, error) {
	// check if we already have the response in cache
	resp, err := c.files.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	// get the responose via the client
	resp, err = c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s - got %s", req.Method, req.URL.Path, resp.Status)
	}

	c.eg.Go(c.write(req.URL.Path+".json", InterceptReadCloser(&resp.Body)))

	return resp, nil
}

func (c *FSClientWriter) write(fileName string, r io.Reader) func() error {
	return func() error {
		return c.gen.Write(fileName, r)
	}
}

func InterceptReadCloser(body *io.ReadCloser) io.Reader {
	r, w := io.Pipe()
	*body = _io.TeeReadCloser(*body, w)

	return r
}

// Wait waits for all files to be written.
func (c FSClientWriter) Wait() error { return c.eg.Wait() }

func (c FSClientWriter) Unseen() []string { return c.fsClient.Unseen() }
