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
	"github.com/Chimder/spoti/internal/handler/http/otel"
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

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
	////////////////////
	// repo := postgres.NewRepository(dbConn)
	// meiliSearchRepo := meilisearchrepo.NewMeiliRepository(meiliConn)
	// redisCache := rediscache.NewRedisCache(redisConn)
	/////////////
	// artistSrv := service.NewArtistService(repo, redisCache, meiliSearchRepo)
	// userSrv := service.NewUserService(repo, redisCache, meiliSearchRepo)
	// playlistSrv := service.NewPlaylistService(repo, redisCache, meiliSearchRepo)
	// trackSrv := service.NewTrackService(repo, redisCache, meiliSearchRepo)
	// albumSrv := service.NewAlbumService(repo, redisCache, meiliSearchRepo)
	/////////////
	// elasticConn := elastic.NewElasticDB(cfg.ElasticSearchUrl)
	// event := scheduler.NewEventWorker(ctx, dbConn, clkhConn)
	// event.Start()
	// defer event.Stop()
	///////////////////
	///////////////////
	// go func() {
	// 	mux := http.NewServeMux()
	// 	mux.Handle("/metrics", promhttp.Handler())

	// 	log.Info().Msg("Prom metrics running :2112/metrics")
	// 	if err := http.ListenAndServe(":2112", mux); err != nil {
	// 		log.Fatal().Err(err).Msg("failed to serve metrics")
	// 	}
	// }()
	// ///////////////////

	// grpcSrv := grpcserver.NewServer(artistSrv, userSrv, playlistSrv, trackSrv, albumSrv)
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Panic().Err(err).Msg("Failed to listen")
	// 	return
	// }

	// go func() {
	// 	log.Info().Msg("gRPC is running on :50051")
	// 	if err := grpcSrv.Serve(lis); err != nil {
	// 		log.Error().Err(err).Msg("gRPC server error")
	// 	}
	// }()

	// <-ctx.Done()
	// log.Info().Msg("Shut down gRPC...")

	// grpcSrv.GracefulStop()
	// log.Info().Msg("gRPC stopped")

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
