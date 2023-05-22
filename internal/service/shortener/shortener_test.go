package shortener

import (
	"context"
	"testing"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/mocks"
	"github.com/amanakin/shortener/internal/service"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestShortener(t *testing.T) {
	cases := []struct {
		name string
		fn   func(t *testing.T, mockRepo *mocks.MockShortenerRepo, mockGen *mocks.MockGenerator)
	}{
		{
			name: "generate collision and retry",
			fn: func(t *testing.T, mockRepo *mocks.MockShortenerRepo, mockGen *mocks.MockGenerator) {
				shortener := Shortener{
					repo:           mockRepo,
					gen:            mockGen,
					defaultScheme:  defaultScheme,
					allowedSchemes: defaultAllowedSchemes,
				}

				collisionLink := domain.Link{
					OriginalURL:  "https://google.com",
					ShortenedURL: "abc",
				}
				newLink := domain.Link{
					OriginalURL:  "https://google.com",
					ShortenedURL: "def",
				}

				first := mockRepo.EXPECT().Store(gomock.Any(), collisionLink).Return(domain.Link{}, service.ErrExist).Times(2)
				second := mockRepo.EXPECT().Store(gomock.Any(), newLink).Return(newLink, nil)
				gomock.InOrder(first, second)

				first = mockGen.EXPECT().Generate(collisionLink.OriginalURL).Return(collisionLink.ShortenedURL).Times(2)
				second = mockGen.EXPECT().Generate(collisionLink.OriginalURL).Return(newLink.ShortenedURL)
				gomock.InOrder(first, second)

				link, created, err := shortener.Shorten(context.Background(), collisionLink.OriginalURL)
				require.NoError(t, err)
				require.True(t, created)
				require.Equal(t, newLink, link)
			},
		},
		{
			name: "original url already exists in repo",
			fn: func(t *testing.T, mockRepo *mocks.MockShortenerRepo, mockGen *mocks.MockGenerator) {
				shortener := Shortener{
					repo:           mockRepo,
					gen:            mockGen,
					defaultScheme:  defaultScheme,
					allowedSchemes: defaultAllowedSchemes,
				}

				wantedLink := domain.Link{
					OriginalURL:  "https://google.com",
					ShortenedURL: "abc",
				}
				oldLink := domain.Link{
					OriginalURL:  "https://google.com",
					ShortenedURL: "def",
				}

				mockRepo.EXPECT().Store(gomock.Any(), wantedLink).Return(oldLink, nil)
				mockGen.EXPECT().Generate(gomock.Any()).Return(wantedLink.ShortenedURL)

				link, created, err := shortener.Shorten(context.Background(), wantedLink.OriginalURL)
				require.NoError(t, err)
				require.False(t, created)
				require.Equal(t, oldLink, link)
			},
		},
		{
			name: "original url is invalid",
			fn: func(t *testing.T, mockRepo *mocks.MockShortenerRepo, mockGen *mocks.MockGenerator) {
				shortener := Shortener{
					repo:           mockRepo,
					gen:            mockGen,
					defaultScheme:  defaultScheme,
					allowedSchemes: defaultAllowedSchemes,
				}

				invalidURL := "invalid url"

				_, _, err := shortener.Shorten(context.Background(), invalidURL)
				require.Error(t, err)
			},
		},
		{
			name: "generate and resolve",
			fn: func(t *testing.T, mockRepo *mocks.MockShortenerRepo, mockGen *mocks.MockGenerator) {
				shortener := Shortener{
					repo:           mockRepo,
					gen:            mockGen,
					defaultScheme:  defaultScheme,
					allowedSchemes: defaultAllowedSchemes,
				}

				link := domain.Link{
					OriginalURL:  "http://google.com",
					ShortenedURL: "abcde",
				}

				mockGen.EXPECT().Generate(link.OriginalURL).Return(link.ShortenedURL)
				mockRepo.EXPECT().Store(gomock.Any(), link).Return(link, nil)
				mockRepo.EXPECT().Get(gomock.Any(), link.ShortenedURL).Return(link.OriginalURL, nil)

				shortenedLink, created, err := shortener.Shorten(context.Background(), link.OriginalURL)
				require.NoError(t, err)
				require.True(t, created)
				require.Equal(t, link, shortenedLink)

				original, err := shortener.Resolve(context.Background(), link.ShortenedURL)
				require.NoError(t, err)
				require.Equal(t, link.OriginalURL, original)
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockShortenerRepo(ctrl)
			mockGen := mocks.NewMockGenerator(ctrl)

			tCase.fn(t, mockRepo, mockGen)
		})
	}
}
