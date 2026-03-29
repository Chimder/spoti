package grpc

import (
	albumv1 "github.com/Chimder/spoti/internal/gen/album/v1"
	artistv1 "github.com/Chimder/spoti/internal/gen/artist/v1"
	playlistv1 "github.com/Chimder/spoti/internal/gen/playlist/v1"
	trackv1 "github.com/Chimder/spoti/internal/gen/track/v1"
	userv1 "github.com/Chimder/spoti/internal/gen/user/v1"
	"github.com/Chimder/spoti/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(
	artistsrv *service.ArtistService,
	usersrv *service.UserService,
	playlistsrv *service.PlaylistService,
	tracksrv *service.TrackService,
	albumsrv *service.AlbumService,
) *grpc.Server {
	srv := grpc.NewServer(
	//(OTel, recovery, auth)
	// grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor(), ...)
	)
	artistHandler := NewArtistHandler(artistsrv)
	playlistHandler := NewPlaylistHandler(playlistsrv)
	userHandler := NewUserHandler(usersrv)
	trackHandler := NewTrackHandler(tracksrv)
	AlbumHandler := NewAlbumHandler(albumsrv)

	artistv1.RegisterPlaylistServiceServer(srv, artistHandler)
	playlistv1.RegisterPlaylistServiceServer(srv, playlistHandler)
	userv1.RegisterUserServiceServer(srv, userHandler)
	trackv1.RegisterTrackServiceServer(srv, trackHandler)
	albumv1.RegisterAlbumServiceServer(srv, AlbumHandler)

	reflection.Register(srv)

	return srv
}
