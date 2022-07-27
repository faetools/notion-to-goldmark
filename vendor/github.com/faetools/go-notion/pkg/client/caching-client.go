package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/faetools/client"
)

type cachingDoer struct {
	cli client.HTTPRequestDoer

	mu    sync.Mutex
	cache map[string]*cachedResponse
}

// NewCachingClient is a wrapper for a HTTPRequestDoer that caches all responses.
func NewCachingClient(cli client.HTTPRequestDoer) client.HTTPRequestDoer {
	return &cachingDoer{
		cli:   cli,
		cache: map[string]*cachedResponse{},
	}
}

type cachedResponse struct {
	body []byte
	resp *http.Response
}

func (c *cachingDoer) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	url := req.URL.String()

	if cached, ok := c.cache[url]; ok {
		return cached.clone(), nil
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	cached := &cachedResponse{body: body, resp: resp}
	c.cache[url] = cached

	return cached.clone(), nil
}

func (r cachedResponse) clone() *http.Response {
	return &http.Response{
		Status:           r.resp.Status,
		StatusCode:       r.resp.StatusCode,
		Proto:            r.resp.Proto,
		ProtoMajor:       r.resp.ProtoMajor,
		ProtoMinor:       r.resp.ProtoMinor,
		Header:           r.resp.Header,
		Body:             io.NopCloser(bytes.NewBuffer(r.body)),
		ContentLength:    r.resp.ContentLength,
		TransferEncoding: r.resp.TransferEncoding,
		Close:            r.resp.Close,
		Uncompressed:     r.resp.Uncompressed,
		Trailer:          r.resp.Trailer,
		Request:          r.resp.Request,
		TLS:              r.resp.TLS,
	}
}
