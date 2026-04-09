package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "baneks.com/internal/api"
	"baneks.com/internal/api/ai"
	"baneks.com/internal/api/baneks"
	memegenerator "baneks.com/internal/api/meme_generator"
	"baneks.com/internal/api/memes"
	c "baneks.com/internal/config"
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
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	aiClient, err := initAiClient(config.CerebrusToken, "127.0.0.1:1080")
	if err != nil {
		log.Fatalf("Create aiClient error: %v", err)
	}

	server := api.InitializeServer(ctx, logger, config.ApiKey)
	g := server.Group("/api")

	// router init
	baneks.InitBanekRouter(g)
	memes.InitMemesRouter(g)
	memegenerator.InitMemeGeneratorRouter(g, m)
	ai.InitAiRouter(g, aiClient)

	serverConfig := echo.StartConfig{
		Address:         ":" + config.Port,
		GracefulTimeout: 4 * time.Second,
		HideBanner:      true,
	}
	if err := serverConfig.Start(ctx, server); err != nil {
		server.Logger.Error("failed to start server", "error", err)
	}
}

func initAiClient(token, proxyAddr string) (*aiPkg.AiClient, error) {
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

	transport := &http.Transport{
		DialContext:           dialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return aiPkg.NewAiClient(httpClient, token, 3, nil)
}
