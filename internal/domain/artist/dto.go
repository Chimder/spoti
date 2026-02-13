package artist

type CreateArtistReq struct {
	Url         string
	Uri         string
	ArtistName string
	Image       string
	Followers   int
	Popularity  int
	Genres      []string
}
