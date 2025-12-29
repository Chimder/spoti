package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"spoti/config"
	postgres_db "spoti/internal/infrastructure/postgres"
	httpgin "spoti/internal/interfaces/http"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//		@title			Spoti Api
//		@version		1.0
//		@description	Similar Spotify Api
//	  @BasePath	/
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg := config.LoadEnv()
	SetupLogger(cfg)

	dbconn, err := postgres_db.Conn(ctx, cfg.PostgresUrl)
	if err != nil {
		log.Panic().Msg("Err conn to db")
		return
	}

	log.Info().Any("", dbconn).Msg("d")
	r := httpgin.Init()
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	}()

	log.Info().Msg("Server is running...")
	<-ctx.Done()
	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}
