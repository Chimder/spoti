package track

import "github.com/google/uuid"

type CreateTrackReq struct {
	AlbumId     uuid.UUID
	RecordingId uuid.UUID
	Name        string
	Number      int16
	DiscNumber  int16
	Explicit    bool
	IsPlayable  bool
	Type        string
	URI         string
	IsLocal     bool
}
