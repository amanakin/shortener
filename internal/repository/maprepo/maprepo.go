package maprepo

import (
	"context"
	"sync"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/service"
)

type Repo struct {
	redirects map[string]string
	originals map[string]string
	mu        sync.RWMutex
}

func New() *Repo {
	return &Repo{
		redirects: make(map[string]string),
		originals: make(map[string]string),
	}
}

func (r *Repo) Store(_ context.Context, link domain.Link) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if shortened, ok := r.originals[link.OriginalURL]; ok {
		return domain.Link{
			OriginalURL:  link.OriginalURL,
			ShortenedURL: shortened,
		}, nil
	}

	if _, ok := r.redirects[link.ShortenedURL]; ok {
		return domain.Link{}, service.ErrExist
	}

	r.redirects[link.ShortenedURL] = link.OriginalURL
	r.originals[link.OriginalURL] = link.ShortenedURL

	return link, nil
}

func (r *Repo) Get(_ context.Context, shortened string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if original, ok := r.redirects[shortened]; ok {
		return original, nil
	}

	return "", service.ErrNotFound
}

func (r *Repo) Close(_ context.Context) {}
