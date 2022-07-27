package client

import (
	"net/http"

	"github.com/faetools/client"
)

type requestValidator struct {
	cli      client.HTTPRequestDoer
	validate func(req *http.Request) error
}

// NewRequestValidator returns a new request validator.
// Only when the validate function returns no error is the original client used.
// Otherwise, the error of the validator is returned.
func NewRequestValidator(cli client.HTTPRequestDoer, validate func(req *http.Request) error) client.HTTPRequestDoer {
	return &requestValidator{
		cli:      cli,
		validate: validate,
	}
}

func (c *requestValidator) Do(req *http.Request) (*http.Response, error) {
	if err := c.validate(req); err != nil {
		return nil, err
	}

	return c.cli.Do(req)
}
