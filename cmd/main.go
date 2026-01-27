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
	c "baneks.com/internal/config"
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

	config, err := c.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Load config error: %v", err)
		return
	}

	m, err := memer.NewMemer(20, 20)
	if err != nil {
		log.Fatalf("Load memer error: %v", err)
	}

	// setup logger
	loggerLevel := slog.LevelDebug
	if config.Environment == c.EnvProd {
		loggerLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:05:05.00"))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	server := server.InitializeServer(ctx, logger, config.ApiKey)
	g := server.Group("/api")

	// router init
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)
	memegenerator.InitMemeGeneratorRouter(g, m)

	serverConfig := echo.StartConfig{
		Address:         ":" + config.Port,
		GracefulTimeout: 4 * time.Second,
		HideBanner:      true,
	}
	if err := serverConfig.Start(ctx, server); err != nil {
		server.Logger.Error("failed to start server", "error", err)
	}
}
