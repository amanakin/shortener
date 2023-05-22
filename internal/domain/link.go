package domain

import "errors"

var (
	// ErrNoURLsLeft is returned when all URLs are used and we can't create more
	ErrNoURLsLeft = errors.New("all URLs are used")
)

type Link struct {
	OriginalURL  string
	ShortenedURL string
}
