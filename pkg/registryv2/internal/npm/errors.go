package npm

import "fmt"

var (
	// ErrStatusCode is returned when the registry returns an unexpected status code.
	ErrStatusCode = fmt.Errorf("npm: unexpected status code")

	// ErrBadContentType is returned when the registry returns an unexpected
	// content-type.
	ErrBadContentType = fmt.Errorf("npm: unexpected content-type")

	// ErrBadManifest is returned when the registry returns content that cannot
	// be processed.
	ErrBadManifest = fmt.Errorf("npm: bad manifest content")
)
