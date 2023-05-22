package shortener

import (
	"context"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/repository"
	"github.com/amanakin/shortener/internal/service"
	"github.com/amanakin/shortener/internal/service/shortener/randgenerator"
)

const (
	defaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	defaultShortLen = 10
	defaultScheme   = "https"
)

var (
	defaultAllowedSchemes = []string{"http", "https"}
)

type Config struct {
	Alphabet       string   `yaml:"alphabet"`
	ShortLen       int      `yaml:"short_len"`
	DefaultScheme  string   `yaml:"default_scheme"`
	AllowedSchemes []string `yaml:"allowed_schemes"`
}

func DefaultConfig() Config {
	return Config{
		Alphabet:       defaultAlphabet,
		ShortLen:       defaultShortLen,
		DefaultScheme:  defaultScheme,
		AllowedSchemes: defaultAllowedSchemes,
	}
}

type Generator interface {
	Generate(input string) string
}

type Shortener struct {
	repo           repository.ShortenerRepo
	gen            Generator
	defaultScheme  string
	allowedSchemes []string
}

func New(repo repository.ShortenerRepo, config Config) *Shortener {
	return &Shortener{
		repo:           repo,
		gen:            randgenerator.New([]byte(config.Alphabet), config.ShortLen),
		defaultScheme:  config.DefaultScheme,
		allowedSchemes: config.AllowedSchemes,
	}
}

func (s *Shortener) Shorten(ctx context.Context, original string) (domain.Link, bool, error) {
	var err error
	original, err = FixValidateURL(original, s.defaultScheme, s.allowedSchemes)
	if err != nil {
		return domain.Link{}, false, err
	}

	for {
		shortened := s.gen.Generate(original)

		var link domain.Link
		link, err = s.repo.Store(ctx, domain.Link{
			OriginalURL:  original,
			ShortenedURL: shortened,
		})

		switch err {
		case nil:
			created := link.ShortenedURL == shortened
			return link, created, nil
		case service.ErrExist:
			continue
		default:
			return domain.Link{}, false, err
		}
	}
}

func (s *Shortener) Resolve(ctx context.Context, url string) (string, error) {
	return s.repo.Get(ctx, url)
}
