package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/amanakin/shortener/internal/handler/http/handler"
	"github.com/amanakin/shortener/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"golang.org/x/exp/slog"
)

const (
	defaultEnabled   = true
	defaultHost      = "localhost"
	defaultPort      = 8080
	defaultRateLimit = 100
	defaultReadLimit = 1024 * 1024
)

type Config struct {
	Enabled bool `yaml:"enabled"`

	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	RateLimit int   `yaml:"rate_limit"`
	ReadLimit int64 `yaml:"read_limit"`
}

func DefaultConfig() Config {
	return Config{
		Enabled:   defaultEnabled,
		Host:      defaultHost,
		Port:      defaultPort,
		RateLimit: defaultRateLimit,
		ReadLimit: defaultReadLimit,
	}
}

type Server struct {
	config Config

	srv       *http.Server
	shortener *handler.ShortenerHandler
	logger    *slog.Logger
}

func loggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("request_info",
					slog.Group("request",
						slog.Duration("duration", time.Since(start)),
						slog.Int("status", ww.Status()),
						slog.Int("bytes_written", ww.BytesWritten()),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("remote_addr", r.RemoteAddr)))
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func New(logger *slog.Logger, shortener service.Shortener, config Config) *Server {
	logger = logger.WithGroup("http")

	srv := &http.Server{
		Addr: net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
	}

	return &Server{
		config:    config,
		srv:       srv,
		shortener: handler.NewShortener(logger, shortener, config.ReadLimit),
		logger:    logger,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	var lc net.ListenConfig
	lsn, err := lc.Listen(ctx, "tcp", s.config.Host+":"+strconv.Itoa(s.config.Port))
	if err != nil {
		return fmt.Errorf("tcp listen: %w", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RealIP) // set req.RemoteAddr from 'X-Real-IP' or 'X-Forwarded-For'
	router.Use(loggerMiddleware(s.logger))
	router.Use(httprate.LimitAll(s.config.RateLimit, time.Second))
	router.Use(middleware.Recoverer)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})

	s.shortener.Register(router)

	s.srv.Handler = router
	s.logger.Info("http server listening",
		slog.String("host", s.config.Host),
		slog.Int("port", s.config.Port))
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		s.logger.Error("shutdown", slog.String("error", err.Error()))
	}
}
