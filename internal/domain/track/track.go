package track

import (
	"github.com/google/uuid"
)

type Track struct {
	Id          uuid.UUID
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
	ISRC        string
	DurationMs  int64
	Popularity  int32
	PlayCount   int64
	AudioURI    string
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
