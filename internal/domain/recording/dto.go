package recording

type CreateRecordingReq struct {
	ISRC       string
	DurationMs int64
	AudioUri   string
}
