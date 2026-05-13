package httpgin

import (
	// _ "csTrade/docs"
	// "csTrade/internal/handlers/middleware"

	"context"

	"github.com/Chimder/spoti/internal/handler/http/middleware"
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
	"github.com/Chimder/spoti/internal/service"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(ctx context.Context, dbConn *pgxpool.Pool, clkhConn driver.Conn,
	// rdsConn *redis.Client, elasticConn *elasticsearch.TypedClient) *gin.Engine {
	rdsConn *redis.Client, meiliConn meilisearch.ServiceManager) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(otelgin.Middleware("spoti-api"))
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics", "/healthz"},
	}))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	repo := postgres.NewRepository(dbConn)
	// _ = elastic.NewElasticRepository(elasticConn)
	meiliSearchRepo := meilisearchrepo.NewMeiliRepository(meiliConn)
	clickHouseRepo := clickhouse.NewListeningEventRepo(clkhConn)
	redisCache := rediscache.NewRedisCache(rdsConn)

	userService := service.NewUserService(repo, redisCache, meiliSearchRepo)
	albumService := service.NewAlbumService(repo, redisCache, meiliSearchRepo)
	artistService := service.NewArtistService(repo, redisCache, meiliSearchRepo)
	playlistService := service.NewPlaylistService(repo, redisCache, meiliSearchRepo)
	trackService := service.NewTrackService(repo, redisCache, meiliSearchRepo)

	userHandler := NewUserHandler(userService, redisCache)
	albumHandler := NewAlbumHandler(albumService)
	artistHandler := NewArtistHandler(*artistService)
	playlistHandler := NewPlaylistHandler(playlistService)
	trackHandler := NewTrackHandler(*trackService)

	listeningEventHandler := NewListeningEventHandler(clickHouseRepo)

	{
		// r.GET("/swagger", ginSwagger.WrapHandler(swaggerfiles.Handler))
		// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		r.GET("/healthz", func(c *gin.Context) {
			c.String(200, "ok")
		})
	}

	event := r.Group("/event")
	{
		event.GET("/listening", listeningEventHandler.AddEvent)
	}

	users := r.Group("/users")
	{
		users.POST("/singUp", userHandler.CreateUser)
		users.POST("/singIn", userHandler.SingInUser)
		users.POST("/refresh", userHandler.RefreshUserToken)
		auth := users.Group("/")
		auth.Use(middleware.AuthUserMiddleware())
		{
			auth.GET("/:id", userHandler.GetUserByID)
			auth.POST("/:userId/playlists/:playlistId/follow", userHandler.FollowPlaylist)
			auth.DELETE("/:userId/playlists/:playlistId/follow", userHandler.UnfollowPlaylist)
			auth.POST("/:userId/artists/:artistId/follow", userHandler.FollowArtist)
			auth.DELETE("/:userId/artists/:artistId/follow", userHandler.UnfollowArtist)
		}
	}

	albums := r.Group("/albums")
	{
		albums.POST("", albumHandler.CreateAlbum)
		albums.GET("/:id", albumHandler.GetAlbumWithTracks)
		albums.GET("", albumHandler.GetAlbumsByIds)
		albums.GET("/new", albumHandler.GetNewReleases)

		albums.GET("/users/:userId", albumHandler.GetUserSavedAlbums)
		albums.PUT("/users/:userId", albumHandler.SaveAlbumsForCurrentUser)
		albums.DELETE("/users/:userId", albumHandler.RemoveAlbumsFromCurrentUser)
		albums.GET("/users/:userId/contains", albumHandler.CheckUserSavedAlbums)
	}
	artists := r.Group("/artists")
	{
		artists.POST("", artistHandler.CreateArtist)
		artists.GET("/:id", artistHandler.GetArtist)
		artists.GET("", artistHandler.GetArtistsByIDs)
		artists.GET("/:id/albums", artistHandler.GetArtistAlbums)
	}

	playlists := r.Group("/playlists")
	{
		playlists.POST("", playlistHandler.CreatePlaylist)
		playlists.GET("/:id", playlistHandler.GetPlaylistById)
		playlists.PUT("/:id", playlistHandler.UpdatePlaylist)

		playlists.POST("/:playlistId/tracks/:trackId", playlistHandler.AddToPlaylist)
		playlists.DELETE("/:playlistId/tracks/:trackId", playlistHandler.DeleteFromPlaylist)

		playlists.GET("/users/:userId", playlistHandler.GetAllUserPlaylists)
	}

	tracks := r.Group("/tracks")
	{
		tracks.POST("", trackHandler.CreateTrack)
		tracks.GET("/:id", trackHandler.GetTrackById)
		tracks.GET("", trackHandler.GetTracksByIds)
		tracks.GET("/artist/:artistId", trackHandler.GetArtistTracks)
	}

	return r
}
