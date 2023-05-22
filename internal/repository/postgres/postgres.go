package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/amanakin/shortener/internal/domain"
	"github.com/amanakin/shortener/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

const (
	defaultEnabled  = true
	defaultHost     = "localhost"
	defaultPort     = 5432
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultDBName   = "postgres"
)

type Config struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func DefaultConfig() Config {
	return Config{
		Enabled:  defaultEnabled,
		Host:     defaultHost,
		Port:     defaultPort,
		User:     defaultUser,
		Password: defaultPassword,
		DBName:   defaultDBName,
	}
}

type Repo struct {
	pool *pgxpool.Pool
}

func New(config Config) (*Repo, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		config.Host, config.User, config.Password, config.DBName, config.Port)

	slog.Info("info", slog.String("dsn", dsn))
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Repo{
		pool: pool,
	}, nil
}

func (r *Repo) Store(ctx context.Context, link domain.Link) (domain.Link, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return link, err
	}
	defer func() {
		if rErr := tx.Rollback(ctx); rErr != nil {
			err = errors.Join(err, fmt.Errorf("rollback: %w", rErr))
		}
	}()

	var shortened string
	err = r.pool.QueryRow(ctx, "SELECT short_url FROM shortener.urls WHERE original_url = $1",
		link.OriginalURL).Scan(&shortened)

	// Original URL already exists
	if err == nil {
		link.ShortenedURL = shortened
		return link, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return link, fmt.Errorf("select by original_url: %w", err)
	}

	var original string
	err = r.pool.QueryRow(ctx, "SELECT original_url FROM shortener.urls WHERE short_url = $1",
		link.ShortenedURL).Scan(&original)

	// Shortened URL already exists
	if err == nil {
		return link, service.ErrExist
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return link, fmt.Errorf("select by short_url: %w", err)
	}

	_, err = r.pool.Exec(ctx, "INSERT INTO shortener.urls (original_url, short_url) VALUES ($1, $2)",
		link.OriginalURL, link.ShortenedURL)
	if err != nil {
		return link, fmt.Errorf("insert link: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return link, fmt.Errorf("commit: %w", err)
	}

	return link, nil
}

func (r *Repo) Get(ctx context.Context, shortened string) (string, error) {
	var original string
	err := r.pool.QueryRow(ctx, "SELECT original_url FROM shortener.urls WHERE short_url = $1", shortened).Scan(&original)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", service.ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("select by short_url: %w", err)
	}

	return original, nil
}

func (r *Repo) Close(_ context.Context) {
	r.pool.Close()
}
