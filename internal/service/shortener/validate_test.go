package shortener

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateURL(t *testing.T) {
	cases := []struct {
		name           string
		u              *url.URL
		allowedSchemes []string
		valid          bool
	}{
		{
			name:           "valid url",
			u:              &url.URL{Scheme: "http", Host: "google.com"},
			allowedSchemes: []string{"http", "https"},
			valid:          true,
		},
		{
			name:           "not allowed scheme",
			u:              &url.URL{Scheme: "ftp", Host: "google.com"},
			allowedSchemes: []string{"http", "https"},
			valid:          false,
		},
		{
			name:           "empty host",
			u:              &url.URL{Scheme: "http"},
			allowedSchemes: []string{"http"},
			valid:          false,
		},
		{
			name:           "nil url",
			u:              nil,
			allowedSchemes: []string{"http"},
			valid:          false,
		},
		{
			name:           "no allowed schemes",
			u:              &url.URL{Scheme: "http", Host: "google.com"},
			allowedSchemes: nil,
			valid:          false,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			valid := validateURL(tCase.u, tCase.allowedSchemes)
			require.Equal(t, tCase.valid, valid)
		})
	}
}

func TestFixValidateURL(t *testing.T) {
	defaultScheme := "https"
	allowedSchemes := []string{"http", "https"}

	cases := []struct {
		name     string
		rawURL   string
		fixedURL string
		errExp   error
	}{
		{
			name:     "add schema to domain",
			rawURL:   "google.com",
			fixedURL: "https://google.com",
			errExp:   nil,
		},
		{
			name:     "valid schema (https), host and port",
			rawURL:   "https://google.com:8080",
			fixedURL: "https://google.com:8080",
			errExp:   nil,
		},
		{
			name:     "valid schema (http), domain, path, params and anchor",
			rawURL:   "http://google.com/some/path?param=1&param=2#anchor",
			fixedURL: "http://google.com/some/path?param=1&param=2#anchor",
			errExp:   nil,
		},
		{
			name:     "invalid schema",
			rawURL:   "ftp://google.com/",
			fixedURL: "",
			errExp:   ErrInvalidURL,
		},
		{
			name:     "empty url",
			rawURL:   "",
			fixedURL: "",
			errExp:   ErrInvalidURL,
		},
		{
			name:     "invalid url format",
			rawURL:   "ht:::///google.com/some/path",
			fixedURL: "",
			errExp:   ErrInvalidURL,
		},
		{
			name:     "empty host",
			rawURL:   "/some/path",
			fixedURL: "",
			errExp:   ErrInvalidURL,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			fixedURL, err := FixValidateURL(tCase.rawURL, defaultScheme, allowedSchemes)
			if tCase.errExp != nil {
				require.ErrorIs(t, err, tCase.errExp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tCase.fixedURL, fixedURL)
			}
		})
	}
}
