package maprepo

import (
	"context"
	"testing"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/service"
	"github.com/stretchr/testify/require"
)

func TestMapRepo(t *testing.T) {
	t.Run("return link if it's not in repo", func(t *testing.T) {
		repo := New()

		link := domain.Link{
			OriginalURL:  "https://google.com",
			ShortenedURL: "123",
		}

		storedLink, err := repo.Store(context.Background(), link)
		require.NoError(t, err)
		require.Equal(t, link, storedLink)
	})

	t.Run("return stored link if original URL in repo", func(t *testing.T) {
		repo := New()

		link := domain.Link{
			OriginalURL:  "https://google.com",
			ShortenedURL: "123",
		}

		storedLink, err := repo.Store(context.Background(), link)
		require.NoError(t, err)
		require.Equal(t, link, storedLink)

		link.ShortenedURL = "456"
		newStoredLink, err := repo.Store(context.Background(), link)
		require.NoError(t, err)
		require.Equal(t, storedLink, newStoredLink)
	})

	t.Run("return ErrExist if shortened URL in repo", func(t *testing.T) {
		repo := New()

		link := domain.Link{
			OriginalURL:  "https://google.com",
			ShortenedURL: "123",
		}

		storedLink, err := repo.Store(context.Background(), link)
		require.NoError(t, err)
		require.Equal(t, link, storedLink)

		link.OriginalURL = "some.other.url"
		_, err = repo.Store(context.Background(), link)
		require.ErrorIs(t, err, service.ErrExist)
	})

	t.Run("resolve returns what was stored", func(t *testing.T) {
		repo := New()

		link := domain.Link{
			OriginalURL:  "https://google.com",
			ShortenedURL: "123",
		}

		storedLink, err := repo.Store(context.Background(), link)
		require.NoError(t, err)
		require.Equal(t, link, storedLink)

		original, err := repo.Get(context.Background(), "123")
		require.NoError(t, err)
		require.Equal(t, "https://google.com", original)
	})
}
