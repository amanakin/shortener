package handler

import (
	"context"
	"fmt"

	"github.com/amanakin/shortener/internal/handler/grpc/api"
	"github.com/amanakin/shortener/internal/service"
)

type ShortenerHandler struct {
	api.UnimplementedShortenerServer

	Shortener service.Shortener
}

func NewShortener(shortener service.Shortener) *ShortenerHandler {
	return &ShortenerHandler{
		Shortener: shortener,
	}
}

func (s *ShortenerHandler) Shorten(ctx context.Context, req *api.ShortenRequest) (*api.ShortenResponse, error) {
	link, created, err := s.Shortener.Shorten(ctx, req.Url)
	if err != nil {
		return nil, fmt.Errorf("shorten: %w", err)
	}

	return &api.ShortenResponse{
		Original:  link.OriginalURL,
		Shortened: link.ShortenedURL,
		Created:   created,
	}, nil
}

func (s *ShortenerHandler) Resolve(ctx context.Context, req *api.ResolveRequest) (*api.ResolveResponse, error) {
	original, err := s.Shortener.Resolve(ctx, req.Shortened)
	if err != nil {
		return nil, fmt.Errorf("resolve: %w", err)
	}

	return &api.ResolveResponse{
		Original: original,
	}, nil
}
