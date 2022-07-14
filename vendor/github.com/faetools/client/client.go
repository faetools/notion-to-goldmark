package client

import (
	"net/http"
	"net/url"
)

// Client defines a standard API client.
//
// Define an alias for this type with methods in order to use it.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	BaseURL *url.URL

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HTTPRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// NewClient creates a new Client, with reasonable defaults.
func NewClient(opts ...Option) (*Client, error) {
	// create a client
	c := Client{}

	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&c); err != nil {
			return nil, err
		}
	}

	// use http.Client, if client not already present
	if c.Client == nil {
		c.Client = &http.Client{}
	}

	return &c, nil
}
