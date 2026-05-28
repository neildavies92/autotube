package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestEnqueueRequestValidateRequiresType(t *testing.T) {
	t.Parallel()

	err := EnqueueRequest{}.Validate()
	if !errors.Is(err, ErrMissingJobType) {
		t.Fatalf("expected ErrMissingJobType, got %v", err)
	}
}

func TestEnqueueRequestValidateRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	err := EnqueueRequest{
		Type:    JobTypeResearchNiche,
		Payload: json.RawMessage(`{`),
	}.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEnqueueRequestValidateAcceptsValidRequest(t *testing.T) {
	t.Parallel()

	err := EnqueueRequest{
		Type:    JobTypeResearchNiche,
		Payload: json.RawMessage(`{"niche_id":1}`),
	}.Validate()
	if err != nil {
		t.Fatalf("validate request: %v", err)
	}
}

func TestHandlerFuncImplementsHandler(t *testing.T) {
	t.Parallel()

	var handler Handler = HandlerFunc(func(_ context.Context, job Job) (Result, error) {
		if job.Type != JobTypeRenderVideo {
			t.Fatalf("expected render job, got %q", job.Type)
		}
		return Result{Output: json.RawMessage(`{"ok":true}`)}, nil
	})

	result, err := handler.HandleJob(context.Background(), Job{Type: JobTypeRenderVideo})
	if err != nil {
		t.Fatalf("handle job: %v", err)
	}

	if string(result.Output) != `{"ok":true}` {
		t.Fatalf("unexpected output: %s", result.Output)
	}
}

func TestRequireValidStatus(t *testing.T) {
	t.Parallel()

	if err := RequireValidStatus(JobStatusQueued); err != nil {
		t.Fatalf("expected queued status to be valid: %v", err)
	}

	if err := RequireValidStatus("unknown"); !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}
