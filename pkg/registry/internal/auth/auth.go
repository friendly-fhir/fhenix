package auth

import (
	"context"
	"net/http"
)

// Client is an interface that models the [http.Client] struct type.
//
// This enables dependency-injecting the client around for different purposes,
// which is leveraged internally in this package so that the registrytest
// package can provide non-networked clients for testing.
type Client interface {
	// Do sends an HTTP request and returns an HTTP response, following
	Do(req *http.Request) (*http.Response, error)
}

// Authentication is an interface for providing oauth2 authentication to the
// registry.
type Authentication interface {
	client(ctx context.Context) (Client, error)
}

// AuthenticationFunc is a function type that implements the [Authentication],
// so that other packages in this project can provide custom clients.
type AuthenticationFunc func(ctx context.Context) (Client, error)

func (a AuthenticationFunc) client(ctx context.Context) (Client, error) {
	return a(ctx)
}

// LoadClient loads a client from the given [Authentication].
func LoadClient(ctx context.Context, auth Authentication) (Client, error) {
	return auth.client(ctx)
}
