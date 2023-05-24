package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amanakin/shortener/internal/handler/grpc"
	"github.com/amanakin/shortener/internal/handler/http"
	"github.com/amanakin/shortener/internal/repository"
	"github.com/amanakin/shortener/internal/repository/maprepo"
	"github.com/amanakin/shortener/internal/repository/postgres"
	"github.com/amanakin/shortener/internal/service"
	"github.com/amanakin/shortener/internal/service/shortener"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HttpConfig      http.Config      `yaml:"http"`
	GrpcConfig      grpc.Config      `yaml:"grpc"`
	PgConfig        postgres.Config  `yaml:"postgres"`
	ShortenerConfig shortener.Config `yaml:"shortener"`
}

func getConfig() (*Config, error) {
	cfg := &Config{
		HttpConfig:      http.DefaultConfig(),
		GrpcConfig:      grpc.DefaultConfig(),
		PgConfig:        postgres.DefaultConfig(),
		ShortenerConfig: shortener.DefaultConfig(),
	}

	configFile := flag.String("c", "", "Path to the YAML configuration file")
	flag.Parse()
	if *configFile == "" {
		return cfg, nil
	}

	file, err := os.Open(*configFile)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	return cfg, nil
}

type Server interface {
	ListenAndServe(ctx context.Context) error
	Stop()
}

func StartServers(logger *slog.Logger, shortenerService service.Shortener, cfg *Config) {
	var servers []Server
	if cfg.HttpConfig.Enabled {
		servers = append(servers, http.New(logger, shortenerService, cfg.HttpConfig))
	}
	if cfg.GrpcConfig.Enabled {
		servers = append(servers, grpc.New(logger, shortenerService, cfg.GrpcConfig))
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	for _, server := range servers {
		go func(server Server) {
			err := server.ListenAndServe(ctx)
			if err != nil {
				logger.Error(err.Error())
			}
		}(server)
	}

	<-ctx.Done()

	for _, server := range servers {
		server.Stop()
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := getConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var repo repository.ShortenerRepo
	if cfg.PgConfig.Enabled {
		repo, err = postgres.New(cfg.PgConfig)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	} else {
		repo = maprepo.New()
	}

	shortenerService := shortener.NewService(repo, cfg.ShortenerConfig)
	StartServers(logger, shortenerService, cfg)
}
