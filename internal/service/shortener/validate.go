package shortener

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid URL")

// validateURL allows non-empty host and provided schemes.
func validateURL(u *url.URL, allowedSchemes []string) bool {
	if u == nil {
		return false
	}

	if u.Host == "" {
		return false
	}

	for _, allowedScheme := range allowedSchemes {
		if u.Scheme == allowedScheme {
			return true
		}
	}

	return false
}

// FixValidateURL validates URL (see validateURL).
// If raw URL not valid it tries to add defaultScheme and validate again.
func FixValidateURL(rawURL string, defaultScheme string, allowedSchemes []string) (string, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err == nil && validateURL(u, allowedSchemes) {
		return rawURL, nil
	}
	// We got invalid URL, but it is already has schema
	if u != nil && u.Scheme != "" {
		return "", ErrInvalidURL
	}

	rawURL = fmt.Sprintf("%s://%s", defaultScheme, rawURL)
	u, err = url.ParseRequestURI(rawURL)
	if err == nil && validateURL(u, allowedSchemes) {
		return rawURL, nil
	}

	return "", ErrInvalidURL
}
