package playlistrepo

import (
	"time"

	"github.com/Chimder/spoti/internal/domain/playlist"
	"github.com/google/uuid"
)

type PlaylistRow struct {
	ID           uuid.UUID `db:"id"`
	OwnerID      uuid.UUID `db:"owner_id"`
	PlaylistName string    `db:"playlist_name"`
	Description  string    `db:"description"`
	Image        string    `db:"image"`
	IsPublic     bool      `db:"is_public"`
	Total        int       `db:"total"`
	CreatedAt    time.Time `db:"created_at"`
}

func (p *PlaylistRow) ToDomain() playlist.Playlist {
	return playlist.Playlist{
		Id:          p.ID,
		Owner:       p.OwnerID,
		Name:        p.PlaylistName,
		Description: p.Description,
		Img:         p.Image,
		Public:      p.IsPublic,
		Total:       uint(p.Total),
	}
}
func PlayListsToDomain(rows []PlaylistRow) []playlist.Playlist {
	result := make([]playlist.Playlist, len(rows))
	for i, row := range rows {
		result[i] = row.ToDomain()
	}
	return result
}
