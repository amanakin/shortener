package grpc

import (
	"fmt"
	"net"
	"strconv"

	"github.com/amanakin/shortener/internal/handler/grpc/api"
	"github.com/amanakin/shortener/internal/handler/grpc/handler"
	"github.com/amanakin/shortener/internal/service"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

import "context"

const (
	defaultHost    = "localhost"
	defaultPort    = 8081
	defaultEnabled = true
)

type Config struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

func DefaultConfig() Config {
	return Config{
		Host:    defaultHost,
		Port:    defaultPort,
		Enabled: defaultEnabled,
	}
}

type Server struct {
	config Config

	srv       *grpc.Server
	shortener *handler.ShortenerHandler
	logger    *slog.Logger
}

type GrpcLogger struct {
	logger *slog.Logger
}

func (l *GrpcLogger) Log(keyvals ...interface{}) error {
	l.logger.Info("info", keyvals...)
	return nil
}

func New(logger *slog.Logger, shortener service.Shortener, config Config) *Server {
	logger = logger.WithGroup("grpc")
	log := &GrpcLogger{logger}

	return &Server{
		config: config,
		srv: grpc.NewServer(
			grpc.ChainUnaryInterceptor(kit.UnaryServerInterceptor(log))),
		shortener: handler.NewShortener(shortener),
		logger:    logger,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	var lc net.ListenConfig
	lsn, err := lc.Listen(ctx, "tcp", net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port)))
	if err != nil {
		return fmt.Errorf("tcp listen: %w", err)
	}

	api.RegisterShortenerServer(s.srv, s.shortener)

	s.logger.Info("grpc server listening",
		slog.String("host", s.config.Host),
		slog.Int("port", s.config.Port))
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
