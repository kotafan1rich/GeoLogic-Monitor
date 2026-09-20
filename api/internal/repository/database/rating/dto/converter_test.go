package dto

import (
	"testing"
	"uuid"
)

func TestToDomainHistoryReturnsEmptyHistory(t *testing.T) {
	t.Parallel()

	trackedLocationID := uuid.New()
	history := ToDomainHistory(trackedLocationID, nil)

	if history == nil {
		t.Fatal("history is nil")
	}
	if history.TrackedLocationID != trackedLocationID {
		t.Fatalf("tracked location ID: got %s, want %s", history.TrackedLocationID, trackedLocationID)
	}
	if history.History == nil || len(history.History) != 0 {
		t.Fatalf("history entries: got %#v, want empty non-nil slice", history.History)
	}
}
