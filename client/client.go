package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/pkg/errors"
)

func New(ctx context.Context, apiKey, apiEndpoint, version string) (*ClientWithResponses, error) {
	bearerTokenProvider, bearerTokenProviderErr := securityprovider.NewSecurityProviderBearerToken(apiKey)
	if bearerTokenProviderErr != nil {
		return nil, bearerTokenProviderErr
	}

	base := cleanhttp.DefaultClient()

	// The generated client won't turn validation errors into actual errors, so we do this
	// inside of a generic middleware.
	base.Transport = Wrap(cleanhttp.DefaultTransport(), func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		if err == nil && resp.StatusCode > 299 {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("status %d: no response body", resp.StatusCode)
			}

			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
		}

		return resp, err
	})

	client, err := NewClientWithResponses(
		apiEndpoint,
		WithHTTPClient(base),
		WithRequestEditorFn(bearerTokenProvider.Intercept),
		// Add a user-agent so we can tell which version these requests came from.
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("user-agent", fmt.Sprintf("catalog-importer/%s", version))
			return nil
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	return client, nil
}

// RoundTripperFunc wraps a function to implement the RoundTripper interface, allowing
// easy wrapping of existing round-trippers.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Wrap allows easy wrapping of an existing RoundTripper with a function that can
// optionally call the original, or do its own thing.
func Wrap(next http.RoundTripper, apply func(req *http.Request, next http.RoundTripper) (*http.Response, error)) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return apply(req, next)
	})
}
