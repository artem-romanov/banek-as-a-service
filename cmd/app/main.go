package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "baneks.com/internal/api"
	cfg "baneks.com/internal/config"
	memesloader "baneks.com/internal/loaders/memes_loader"
	aiPkg "baneks.com/pkg/ai"
	"baneks.com/pkg/memer"
	"github.com/labstack/echo/v5"
	"golang.org/x/net/proxy"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	config, err := cfg.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Load config error: %v", err)
		return
	}

	logger := setupLogger(&config)
	slog.SetDefault(logger)

	deps, err := initializeDependencies(&config)
	if err != nil {
		slog.Error("dependencies initialization fail", "err", err)
		os.Exit(1)
	}

	server := api.InitializeServer(ctx, logger, config.ApiKey, api.Dependencies{
		Memer:    deps.Memer,
		AiClient: deps.AiClient,
		MemeLoader: deps.MemesLoader,
	})

	serverConfig := echo.StartConfig{
		Address:         ":" + config.Port,
		GracefulTimeout: 4 * time.Second,
		HideBanner:      true,
	}
	if err := serverConfig.Start(ctx, server); err != nil {
		server.Logger.Error("failed to start server", "error", err)
	}
}

type AppDependencies struct {
	Memer    *memer.Memer
	AiClient *aiPkg.AiClient
	MemesLoader memesloader.MemeLoader
}

func initializeDependencies(config *cfg.AppConfig) (*AppDependencies, error) {
	errs := make([]error, 0)

	m, err := memer.NewMemer(20, 20)
	if err != nil {
		errs = append(errs, fmt.Errorf("memer init fail: %w", err))
	}

	aiProvider := aiPkg.NewGroqOss120Provider(config.GroqToken)

	proxyAddr := ""
	if config.Environment == cfg.EnvProd {
		proxyAddr = "127.0.0.1:1080"
	}

	aiClient, err := initAiClient(aiProvider, proxyAddr)
	if err != nil {
		errs = append(errs, fmt.Errorf("aiClient init fail: %w", err))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &AppDependencies{
		Memer:    m,
		AiClient: aiClient,
		MemesLoader: memesloader.NewQablydauMemeLoader(
			&http.Client{
				Timeout: 15 * time.Second,
			},
		),
	}, nil
}
func initAiClient(provider aiPkg.Provider, proxyAddr string) (*aiPkg.AiClient, error) {
	var transport *http.Transport

	if proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}

		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			if cd, ok := dialer.(proxy.ContextDialer); ok {
				return cd.DialContext(ctx, network, addr)
			}
			return dialer.Dial(network, addr)
		}

		transport = &http.Transport{
			DialContext:           dialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			ExpectContinueTimeout: 1 * time.Second,
		}
	} else {
		transport = &http.Transport{
			TLSHandshakeTimeout:   10 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	c := aiPkg.NewAiClient(httpClient, provider, 3)
	return c, nil
}

func setupLogger(config *cfg.AppConfig) *slog.Logger {
	loggerLevel := slog.LevelDebug
	if config.Environment == cfg.EnvProd {
		loggerLevel = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))
}
