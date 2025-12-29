package album

type AlbumTracks struct {
	Href     string
	Limit    int
	Next     interface{}
	Offset   int
	Previous interface{}
	Total    int
	Items    []struct {
		Artists []struct {
			ExternalUrls struct {
				Spotify string
			}
			Href string
			ID   string
			Name string
			Type string
			URI  string
		}
		AvailableMarkets []string
		DiscNumber       int
		DurationMs       int
		Explicit         bool
		ExternalUrls     struct {
			Spotify string
		}
		Href        string
		ID          string
		Name        string
		PreviewURL  interface{}
		TrackNumber int
		Type        string
		URI         string
		IsLocal     bool
	}
}
