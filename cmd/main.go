package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "baneks.com/internal/api"
	"baneks.com/internal/api/baneks"
	memegenerator "baneks.com/internal/api/meme_generator"
	"baneks.com/internal/api/memes"
	"baneks.com/internal/api/middlewares"
	"baneks.com/internal/config"
	"baneks.com/pkg/memer"
	"github.com/labstack/echo/v5"
)

func main() {
	globalCtx := context.Background()

	ctx, cancel := signal.NotifyContext(
		globalCtx,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	config, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Load config error: %v", err)
		return
	}
	guard := middlewares.New(config.ApiKey)
	m, err := memer.NewMemer(20, 20)
	if err != nil {
		log.Fatalf("Load memer error: %v", err)
	}

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// TODO: distinct between prod and dev,
		// set error for prod, debug for dev
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:05:05.00"))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	server := server.InitializeServer(ctx, logger)
	g := server.Group("/api")

	// global middlewares init
	g.Use(guard.GuardWithSecretMiddleware)

	// router init
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)
	memegenerator.InitMemeGeneratorRouter(g, m)

	serverConfig := echo.StartConfig{
		Address:         ":8888",
		GracefulTimeout: 1 * time.Second,
		HideBanner:      true,
	}
	if err := serverConfig.Start(ctx, server); err != nil {
		server.Logger.Error("failed to start server", "error", err)
	}
}
