package fake

import (
	"fmt"
	"net/http"

	"github.com/faetools/client"
	_client "github.com/faetools/go-notion/pkg/client"
	"github.com/faetools/go-notion/pkg/notion"
)

// NewClient returns a new notion client returning fake results
// and the underlying Doer which you can use to check whether you have
// tested all fake responses.
func NewClient() (*notion.Client, *_client.FSClient, error) {
	fsClient, err := _client.NewFSClient(responses, notFoundResponse)
	if err != nil {
		return nil, nil, err
	}

	c, err := notion.NewDefaultClient("", client.WithHTTPClient(fsClient))

	return c, fsClient, err
}

func notFoundResponse(path string) any {
	return notion.ErrorResponse{
		Code:    fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)),
		Message: fmt.Sprintf("no response found for %s", path),
		Object:  "error",
		Status:  http.StatusNotFound,
	}
}
