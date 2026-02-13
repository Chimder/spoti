package track

import (
	"github.com/google/uuid"
)

type Track struct {
	ID          uuid.UUID
	AlbumID     uuid.UUID
	RecordingID uuid.UUID
	Name        string
	Number      int16
	DiscNumber  int16
	Explicit    bool
	IsPlayable  bool
	Type        string
	URI         string
	IsLocal     bool
}

// type TrackWithRecording struct {
// 	ID          uuid.UUID
// 	AlbumID     uuid.UUID
// 	RecordingID uuid.UUID
// 	Name        string
// 	Number      int16
// 	DiscNumber  int16
// 	Explicit    bool
// 	IsPlayable  bool
// 	Type        string
// 	URI         string
// 	IsLocal     bool
// }
