package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JobType string

const (
	JobTypeResearchNiche   JobType = "research_niche"
	JobTypeGenerateDraft   JobType = "generate_draft"
	JobTypeRenderVideo     JobType = "render_video"
	JobTypeUploadVideo     JobType = "upload_video"
	JobTypeSyncAnalytics   JobType = "sync_analytics"
	JobTypeReviewReadiness JobType = "review_readiness"
)

type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCanceled  JobStatus = "canceled"
)

var (
	ErrMissingJobType = errors.New("workflow setup error: missing job type")
	ErrInvalidStatus  = errors.New("workflow setup error: invalid job status")
)

type EnqueueRequest struct {
	Type    JobType
	NicheID *int64
	Payload json.RawMessage
}

func (request EnqueueRequest) Validate() error {
	if request.Type == "" {
		return ErrMissingJobType
	}

	if len(request.Payload) > 0 && !json.Valid(request.Payload) {
		return fmt.Errorf("workflow setup error: invalid job payload JSON")
	}

	return nil
}

type Job struct {
	ID         int64
	Type       JobType
	Status     JobStatus
	NicheID    *int64
	Payload    json.RawMessage
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type Result struct {
	Output json.RawMessage
}

type Event struct {
	JobID     int64
	Type      string
	Message   string
	Payload   json.RawMessage
	CreatedAt time.Time
}

type Queue interface {
	Enqueue(context.Context, EnqueueRequest) (Job, error)
	Next(context.Context, JobType) (Job, error)
	MarkRunning(context.Context, int64) (Job, error)
	MarkSucceeded(context.Context, int64, Result) (Job, error)
	MarkFailed(context.Context, int64, error) (Job, error)
	MarkCanceled(context.Context, int64) (Job, error)
}

type EventSink interface {
	RecordEvent(context.Context, Event) error
}

type Handler interface {
	HandleJob(context.Context, Job) (Result, error)
}

type HandlerFunc func(context.Context, Job) (Result, error)

func (fn HandlerFunc) HandleJob(ctx context.Context, job Job) (Result, error) {
	return fn(ctx, job)
}

func ValidStatus(status JobStatus) bool {
	switch status {
	case JobStatusQueued, JobStatusRunning, JobStatusSucceeded, JobStatusFailed, JobStatusCanceled:
		return true
	default:
		return false
	}
}

func RequireValidStatus(status JobStatus) error {
	if !ValidStatus(status) {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	return nil
}
