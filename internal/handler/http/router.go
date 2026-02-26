package httpgin

import (
	// _ "csTrade/docs"
	// "csTrade/internal/handlers/middleware"

	"context"
	"spoti/internal/repository/clickhouse"
	"spoti/internal/repository/postgres"
	"spoti/internal/service"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(ctx context.Context, dbConn *pgxpool.Pool, clkhConn driver.Conn) *gin.Engine {
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
	clickHouseRepo := clickhouse.NewListeningEventRepo(clkhConn)

	userService := service.NewUserService(repo)
	albumService := service.NewAlbumService(repo)
	artistService := service.NewArtistService(repo)
	playlistService := service.NewPlaylistService(repo)
	trackService := service.NewTrackService(repo)

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
		users.POST("/:userId/artists/:artistId/follow", userHandler.FollowPlaylist)
		users.DELETE("/:userId/artists/:artistId/follow", userHandler.UnfollowArtist)
	}

	albums := r.Group("/albums")
	{
		albums.POST("", albumHandler.CreateAlbum)
		albums.GET("/:id", albumHandler.GetAlbumJson)
		albums.GET("", albumHandler.GetAlbumsByIds)
		albums.GET("/:id/tracks", albumHandler.GetAlbumTracks)

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
