package shortener

import (
	"context"
	"fmt"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/repository"
	"github.com/amanakin/shortener/internal/service"
	"github.com/amanakin/shortener/internal/service/shortener/hashgenerator"
	"github.com/amanakin/shortener/internal/service/shortener/randgenerator"
)

const (
	defaultAlphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	defaultShortLen      = 10
	defaultScheme        = "https"
	defaultHashGenerator = true
)

var (
	defaultAllowedSchemes = []string{"http", "https"}
)

type Config struct {
	Alphabet       string   `yaml:"alphabet"`
	ShortLen       int      `yaml:"short_len"`
	DefaultScheme  string   `yaml:"default_scheme"`
	AllowedSchemes []string `yaml:"allowed_schemes"`
	HashGenerator  bool     `yaml:"hash_generator"`
}

func DefaultConfig() Config {
	return Config{
		Alphabet:       defaultAlphabet,
		ShortLen:       defaultShortLen,
		DefaultScheme:  defaultScheme,
		AllowedSchemes: defaultAllowedSchemes,
		HashGenerator:  defaultHashGenerator,
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

func NewService(repo repository.ShortenerRepo, config Config) *Shortener {
	var gen Generator
	if config.HashGenerator {
		gen = hashgenerator.New([]byte(config.Alphabet), config.ShortLen)
	} else {
		gen = randgenerator.New([]byte(config.Alphabet), config.ShortLen)
	}

	return &Shortener{
		repo:           repo,
		gen:            gen,
		defaultScheme:  config.DefaultScheme,
		allowedSchemes: config.AllowedSchemes,
	}
}

func (s *Shortener) Shorten(ctx context.Context, original string) (domain.Link, bool, error) {
	var err error
	original, err = FixValidateURL(original, s.defaultScheme, s.allowedSchemes)
	if err != nil {
		return domain.Link{}, false, fmt.Errorf("validating URL %q: %w", original, err)
	}

	for {
		shortened := s.gen.Generate(original)
		link := domain.Link{
			OriginalURL:  original,
			ShortenedURL: shortened,
		}

		link, err = s.repo.Store(ctx, link)
		switch err {
		case nil:
			created := link.ShortenedURL == shortened
			return link, created, nil
		case service.ErrExist:
			continue
		default:
			return domain.Link{}, false, fmt.Errorf("repository store: %w", err)
		}
	}
}

func (s *Shortener) Resolve(ctx context.Context, shortened string) (string, error) {
	original, err := s.repo.Get(ctx, shortened)
	if err != nil {
		return "", fmt.Errorf("repository get: %w", err)
	}
	return original, nil
}
