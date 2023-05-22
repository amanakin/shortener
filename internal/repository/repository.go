package repository

import (
	"context"

	"github.com/amanakin/shortener/internal/domain"
)

type ShortenerRepo interface {
	// Store saves link in repository if there's no such link.
	// If original URL already exists, it must return it.
	// If shortened URL already exists, it must return service.ErrExist.
	// Above rules must be followed in specified order.
	Store(ctx context.Context, link domain.Link) (domain.Link, error)
	// Get gets original URL from shortened.
	// If shortened URL is not found It must return service.ErrNotFound.
	Get(ctx context.Context, shortened string) (string, error)
	Close(ctx context.Context)
}
