package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scraper/internal/services/fetcher"
	http2 "scraper/pkg/http"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	config "scraper/configs"
	"scraper/internal/services/sender"
	"scraper/pkg/env"
	"scraper/pkg/health"
	"scraper/pkg/logger"
)

func main() {
	var (
		wg      sync.WaitGroup
		port    = env.Get("PORT", config.Port)
		isDebug = env.Get("DEBUG", "false") == "true"
		timeout = env.Get("TIMEOUT", "10s")
		mux     = http.NewServeMux()
		server  = http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		}
	)

	_ = timeout

	lg, logSync := logger.MustNew(config.ServiceName, isDebug)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	health.Healthz(mux)
	health.Readyz(ctx, mux, time.Second*15)

	go func() {
		wg.Add(1)
		defer wg.Done()
		lg.Info("starting http server...", zap.String("http_port", port))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error("http server error", zap.Error(err))
		}

		lg.Debug("http server has been successfully stopped")
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		<-ctx.Done()

		lg.Debug("calling http server shutdown...")
		err := server.Shutdown(ctx)
		if err != nil {
			lg.Error("unexpected http server error", zap.Error(err))
		}
	}()

	client, err := http2.NewHttpClientWithCookies(time.Second * 15)
	if err != nil {
		lg.Error("http client error", zap.Error(err))
		cancel()
	}

	dec := fetcher.NewSoupDecorator(client)
	senderservice := sender.New(dec, lg)

	ticker := time.NewTicker(time.Minute * 60)
	go func() {
		<-ctx.Done()

		lg.Debug("calling ticker stop ...")

		ticker.Stop()
	}()

	go func() {
		lg.Info("starting to fetch movies...")

		wg.Add(1)
		defer wg.Done()

		err := senderservice.Broadcast(ctx)
		if err != nil {
			lg.Error("broadcast error", zap.Error(err))
		}

		for {
			select {
			case <-ctx.Done():
				lg.Debug("ticker has been successfully stopped")

				return
			case <-ticker.C:
				// fixme: pipe here
				err = senderservice.Broadcast(ctx)
				if err != nil {
					lg.Error("broadcast error", zap.Error(err))
				}
			}
		}
	}()

	<-ctx.Done()

	wg.Wait()

	lg.Info("application has finished its work")
	_ = logSync()
}
