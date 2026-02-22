package recording

import (
	"context"

	"github.com/google/uuid"
)

type RecordingRepository interface {
	CreateRecording(ctx context.Context, r CreateRecordingReq) (uuid.UUID, error)
	GetRecordingById(ctx context.Context, recordingId string) (Recording, error)
}
