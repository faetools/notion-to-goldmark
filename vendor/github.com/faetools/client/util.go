package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

const (
	// parameter styles
	form   = "form"
	simple = "simple"

	contentType = "Content-Type"
	json        = "json"

	// MIMEApplicationJSON defines standard type "application/json"
	MIMEApplicationJSON = "application/json"

	// ContentType is the header key to define the type of content.
	ContentType = "Content-Type"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// HTTPRequestDoer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}


// MustParseURL parses the URL. It panics if there is an error.
func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(fmt.Sprintf("OpenAPI contains unparsable operation path: %s", err))
	}

	return u
}

// AddQueryParam adds a certain parameter with its value to the query.
func AddQueryParam(query url.Values, paramName string, value interface{}) error {
	queryFrag, err := runtime.StyleParamWithLocation(form, true, paramName, runtime.ParamLocationQuery, value)
	if err != nil {
		return err
	}

	parsed, err := url.ParseQuery(queryFrag)
	if err != nil {
		return err
	}

	for k, v := range parsed {
		for _, v2 := range v {
			query.Add(k, v2)
		}
	}

	return nil
}

// GetPathParam returns the path parameter value.
func GetPathParam(paramName string, value interface{}) (string, error) {
	p, err := runtime.StyleParamWithLocation(simple, false, paramName, runtime.ParamLocationPath, value)
	switch {
	case err != nil:
		return "", err
	case p == "":
		return "", fmt.Errorf("empty path parameter %q", paramName)
	default:
		return p, nil
	}
}
