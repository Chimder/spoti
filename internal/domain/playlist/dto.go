package playlist

import "github.com/google/uuid"

type CreatePlaylistReq struct {
	OwnerId      uuid.UUID
	PlaylistName string
	Description  string
	Image        string
	IsPublic     bool
}

type UpdatePlaylistReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Public      *bool   `json:"public"`
}
