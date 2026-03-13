package httpgin

import (
	// _ "csTrade/docs"
	// "csTrade/internal/handlers/middleware"

	"context"
	"spoti/internal/repository/clickhouse"
	meilisearchrepo "spoti/internal/repository/meilisearch"
	"spoti/internal/repository/postgres"
	rediscache "spoti/internal/repository/redis"
	"spoti/internal/service"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(ctx context.Context, dbConn *pgxpool.Pool, clkhConn driver.Conn,
	// rdsConn *redis.Client, elasticConn *elasticsearch.TypedClient) *gin.Engine {
	rdsConn *redis.Client, meiliConn meilisearch.ServiceManager) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

	userHandler := NewUserHandler(userService)
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
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("/:userId/playlists/:playlistId/follow", userHandler.FollowPlaylist)
		users.DELETE("/:userId/playlists/:playlistId/follow", userHandler.UnfollowPlaylist)
		users.POST("/:userId/artists/:artistId/follow", userHandler.FollowArtist)
		users.DELETE("/:userId/artists/:artistId/follow", userHandler.UnfollowArtist)
	}

	albums := r.Group("/albums")
	{
		albums.POST("", albumHandler.CreateAlbum)
		albums.GET("/:id", albumHandler.GetAlbumJson)
		albums.GET("", albumHandler.GetAlbumsByIds)
		albums.GET("/:id/tracks", albumHandler.GetAlbumTracks)
		albums.GET("/new", albumHandler.GetNewReleases)

		albums.GET("/users/:userId", albumHandler.GetUserSavedAlbums)
		albums.PUT("/users/:userId", albumHandler.SaveAlbumsForCurrentUser)
		albums.DELETE("/users/:userId", albumHandler.RemoveAlbumsFromCurrentUser)
		albums.GET("/users/:userId/contains", albumHandler.CheckUsersSavedAlbums)
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
