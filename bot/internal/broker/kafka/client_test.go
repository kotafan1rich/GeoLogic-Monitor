package kafka

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/broker"
)

type producerStub struct {
	records []*kgo.Record
	err     error
}

func (s *producerStub) ProduceSync(
	_ context.Context,
	records ...*kgo.Record,
) kgo.ProduceResults {
	s.records = append(s.records, records...)

	results := make(kgo.ProduceResults, 0, len(records))
	for _, record := range records {
		results = append(results, kgo.ProduceResult{
			Record: record,
			Err:    s.err,
		})
	}

	return results
}

func TestPublishDLQ(t *testing.T) {
	producer := &producerStub{}
	client := &Client{producer: producer}
	recordErr := errors.New("invalid notification")
	record := &kgo.Record{
		Topic:     "notifications",
		Partition: 2,
		Offset:    42,
		Key:       []byte("notification-key"),
		Value:     []byte(`{"broken":`),
	}

	beforePublish := time.Now().UTC()
	if err := client.publishDLQ(context.Background(), record, recordErr); err != nil {
		t.Fatalf("publishDLQ() error = %v", err)
	}
	afterPublish := time.Now().UTC()

	if len(producer.records) != 1 {
		t.Fatalf("produced records = %d, want 1", len(producer.records))
	}

	got := producer.records[0]
	if got.Topic != "notifications.dlq" {
		t.Errorf("topic = %q, want %q", got.Topic, "notifications.dlq")
	}
	if !bytes.Equal(got.Key, record.Key) {
		t.Errorf("key = %q, want %q", got.Key, record.Key)
	}
	if !bytes.Equal(got.Value, record.Value) {
		t.Errorf("value = %q, want %q", got.Value, record.Value)
	}

	headers := make(map[string]string, len(got.Headers))
	for _, header := range got.Headers {
		headers[header.Key] = string(header.Value)
	}

	wantHeaders := map[string]string{
		"x-original-topic":     record.Topic,
		"x-original-partition": strconv.Itoa(int(record.Partition)),
		"x-original-offset":    strconv.FormatInt(record.Offset, 10),
		"x-attempts":           strconv.Itoa(1 + len(retryBackoffs)),
		"x-error":              recordErr.Error(),
	}
	for key, want := range wantHeaders {
		if got := headers[key]; got != want {
			t.Errorf("header %q = %q, want %q", key, got, want)
		}
	}

	failedAt, err := time.Parse(time.RFC3339, headers["x-failed-at"])
	if err != nil {
		t.Fatalf("parse x-failed-at: %v", err)
	}
	if failedAt.Before(beforePublish.Truncate(time.Second)) ||
		failedAt.After(afterPublish) {
		t.Errorf(
			"x-failed-at = %v, want between %v and %v",
			failedAt,
			beforePublish,
			afterPublish,
		)
	}
}

func TestPublishDLQReturnsProducerError(t *testing.T) {
	producerErr := errors.New("broker unavailable")
	client := &Client{producer: &producerStub{err: producerErr}}

	err := client.publishDLQ(
		context.Background(),
		&kgo.Record{Topic: "notifications"},
		errors.New("handler failed"),
	)
	if !errors.Is(err, producerErr) {
		t.Fatalf("publishDLQ() error = %v, want wrapped %v", err, producerErr)
	}
}

func TestHandleWithRetry(t *testing.T) {
	originalBackoffs := retryBackoffs
	retryBackoffs = [len(retryBackoffs)]time.Duration{}
	t.Cleanup(func() {
		retryBackoffs = originalBackoffs
	})

	t.Run("succeeds immediately", func(t *testing.T) {
		calls := 0
		err := handleWithRetry(
			context.Background(),
			func(context.Context, broker.Message) error {
				calls++
				return nil
			},
			broker.Message{},
			nil,
		)
		if err != nil {
			t.Fatalf("handleWithRetry() error = %v", err)
		}
		if calls != 1 {
			t.Errorf("handler calls = %d, want 1", calls)
		}
	})

	t.Run("succeeds after retry", func(t *testing.T) {
		calls := 0
		err := handleWithRetry(
			context.Background(),
			func(context.Context, broker.Message) error {
				calls++
				if calls < 3 {
					return errors.New("temporary failure")
				}
				return nil
			},
			broker.Message{},
			nil,
		)
		if err != nil {
			t.Fatalf("handleWithRetry() error = %v", err)
		}
		if calls != 3 {
			t.Errorf("handler calls = %d, want 3", calls)
		}
	})

	t.Run("returns last error after all attempts", func(t *testing.T) {
		calls := 0
		lastErr := errors.New("permanent failure")
		err := handleWithRetry(
			context.Background(),
			func(context.Context, broker.Message) error {
				calls++
				return lastErr
			},
			broker.Message{},
			nil,
		)
		if !errors.Is(err, lastErr) {
			t.Fatalf("handleWithRetry() error = %v, want %v", err, lastErr)
		}
		if want := 1 + len(retryBackoffs); calls != want {
			t.Errorf("handler calls = %d, want %d", calls, want)
		}
	})

	t.Run("stops when context is canceled", func(t *testing.T) {
		retryBackoffs[0] = time.Hour
		defer func() {
			retryBackoffs[0] = 0
		}()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		calls := 0
		err := handleWithRetry(
			ctx,
			func(context.Context, broker.Message) error {
				calls++
				return errors.New("failure")
			},
			broker.Message{},
			nil,
		)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf(
				"handleWithRetry() error = %v, want context.Canceled",
				err,
			)
		}
		if calls != 1 {
			t.Errorf("handler calls = %d, want 1", calls)
		}
	})
}
