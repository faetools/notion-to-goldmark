package client

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Option allows setting custom parameters during construction.
type Option func(*Client) error

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HTTPRequestDoer) Option {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) Option {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// WithBasicAuth adds basic authentication to all requests.
func WithBasicAuth(username, password string) Option {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	})
}

// WithBearer adds authentication via bearer token.
func WithBearer(bearer string) Option {
	if !strings.HasPrefix(bearer, "Bearer ") {
		bearer = "Bearer " + bearer
	}

	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", bearer)
		return nil
	})
}

// WithToken adds basic authentication via token to all requests.
func WithToken(token string) Option {
	return WithBasicAuth("", token)
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		// ensure the server URL always has a trailing slash
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}

		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.BaseURL = newBaseURL
		return nil
	}
}
