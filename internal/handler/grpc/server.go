package grpc

import (
	artistv1 "github.com/Chimder/spoti/internal/gen/artist/v1"
	"github.com/Chimder/spoti/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(artserv *service.ArtistService) *grpc.Server {
	srv := grpc.NewServer(
	//(OTel, recovery, auth)
	// grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor(), ...)
	)
	artistHandler := NewArtistHandler(artserv)
	artistv1.RegisterPlaylistServiceServer(srv, artistHandler)

	reflection.Register(srv)

	return srv
}
