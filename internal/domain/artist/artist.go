package artist

type Artist struct {
	ID         string
	Url        string
	Followers  uint64
	Genres     []string
	Image      string
	Name       string
	Popularity uint8
	Type       string
	URI        string
}
