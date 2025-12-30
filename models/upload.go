package models

type UploadTask struct {
	ID                string
	Hash              string
	JobID             string
	StatusURL         string
	FileID            string
	Status            string // pending/processing/completed/failed
	CreatedAt         int64
	LastStatusCheckAt int64
}

const (
	UploadStatusPending    = "pending"
	UploadStatusProcessing = "processing"
	UploadStatusCompleted  = "completed"
	UploadStatusFailed     = "failed"
)
