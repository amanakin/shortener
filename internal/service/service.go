package service

import (
	"context"
	"errors"

	"github.com/amanakin/shortener/internal/domain"
)

var (
	// ErrExist is returned when shortened URL already exists.
	ErrExist = errors.New("shortened URL exists")
	// ErrNotFound is returned when URL is not found
	ErrNotFound = errors.New("URL not found")
)

type Shortener interface {
	// Shorten creates short URL from origin URL, and returns if already created
	Shorten(ctx context.Context, original string) (domain.Link, bool, error)
	// Resolve gets origin URL from previously shortened
	Resolve(ctx context.Context, shortened string) (string, error)
}
