package registry

import (
	"context"
	"net/http"

	"github.com/friendly-fhir/fhenix/pkg/registry/internal/auth"
	"golang.org/x/oauth2"
)

// Authentication is an interface for providing oauth2 authentication to the
// registry.
type Authentication = auth.Authentication

// NoAuthentication returns an [Authentication] option that provides no
// real authentication. This is the default behavior if no authentication
// is specified.
func NoAuthentication() Authentication {
	return auth.AuthenticationFunc(func(context.Context) (auth.Client, error) {
		return http.DefaultClient, nil
	})
}

// StaticTokenSource returns an [Authentication] that always returns the same token.
func StaticTokenSource(token string) Authentication {
	return TokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
}

// StaticTokenSource returns an [Authentication] that always returns the same
// token.
func TokenSource(ts oauth2.TokenSource) Authentication {
	return auth.AuthenticationFunc(func(ctx context.Context) (auth.Client, error) {
		return oauth2.NewClient(ctx, ts), nil
	})
}
