package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chimder/spoti/config"
	httpgin "github.com/Chimder/spoti/internal/handler/http"
	"github.com/Chimder/spoti/internal/otel"
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
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

	otelShutdown, err := otel.Setup(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("err setup OTel")
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			log.Error().Err(err).Msg("OTel shutdown error")
		}
	}()

	dbConn, err := postgres.NewConn(ctx, cfg.PostgresUrl)
	if err != nil {
		log.Panic().Msg("Err conn to postgres")
		return
	}

	clkhConn, err := clickhouse.Conn(ctx)
	if err != nil {
		log.Panic().Msg("Err conn to clickhouse")
		return
	}
	redisConn := rediscache.Conn(cfg.RedisUrl)

	meiliConn := meilisearchrepo.NewMeiliDB(cfg.MeiliSearchUrl)
	// elasticConn := elastic.NewElasticDB(cfg.ElasticSearchUrl)

	// event := scheduler.NewEventWorker(ctx, dbConn, clkhConn)
	// event.Start()
	// defer event.Stop()

	r := httpgin.Init(ctx, dbConn, clkhConn, redisConn, meiliConn)
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
