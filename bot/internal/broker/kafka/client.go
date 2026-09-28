package kafka

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/broker"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
)

const dlqTopicSuffix = ".dlq"

var retryBackoffs = [...]time.Duration{
	time.Second,
	2 * time.Second,
	4 * time.Second,
}

type Client struct {
	client   *kgo.Client
	producer syncProducer
	log      *logger.Logger
	topic    string
	groupID  string
}

type syncProducer interface {
	ProduceSync(ctx context.Context, records ...*kgo.Record) kgo.ProduceResults
}

func New(log *logger.Logger, brokers []string, topic string, groupID string) (*Client, error) {
	const op = "broker.kafka.New"

	if log == nil {
		return nil, fmt.Errorf("%s: logger is required", op)
	}
	if len(brokers) == 0 {
		return nil, fmt.Errorf("%s: brokers are required", op)
	}
	if topic == "" {
		return nil, fmt.Errorf("%s: topic is required", op)
	}
	if groupID == "" {
		return nil, fmt.Errorf("%s: group ID is required", op)
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup(groupID),

		kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: create client: %w", op, err)
	}

	return &Client{
		client:   client,
		producer: client,
		log:      log,
		topic:    topic,
		groupID:  groupID,
	}, nil
}

func (c *Client) Consume(ctx context.Context, handler broker.Handler) error {
	const op = "broker.kafka.Client.Consume"

	if handler == nil {
		return fmt.Errorf("%s: handler is required", op)
	}

	c.log.InfoContext(
		ctx,
		"Kafka consumer started",
		"topic", c.topic,
		"group_id", c.groupID,
		"dlq_topic", c.topic+dlqTopicSuffix,
	)

	for {
		fetches := c.client.PollRecords(ctx, 1)
		if fetches.IsClientClosed() {
			return nil
		}

		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("%s: context: %w", op, err)
		}

		if fetchErrors := fetches.Errors(); len(fetchErrors) > 0 {
			fetchErr := fetchErrors[0]
			return fmt.Errorf(
				"%s: fetch topic %q partition %d: %w",
				op,
				fetchErr.Topic,
				fetchErr.Partition,
				fetchErr.Err,
			)
		}

		records := fetches.Records()
		if len(records) == 0 {
			continue
		}

		record := records[0]

		message := broker.Message{
			Key:   record.Key,
			Value: record.Value,
		}

		handleErr := handleWithRetry(
			ctx,
			handler,
			message,
			func(attempt int, backoff time.Duration, err error) {
				c.log.WarnContext(
					ctx,
					"retrying Kafka message",
					"topic", record.Topic,
					"partition", record.Partition,
					"offset", record.Offset,
					"attempt", attempt,
					"max_attempts", 1+len(retryBackoffs),
					"backoff", backoff,
					"err", err,
				)
			},
		)
		movedToDLQ := false
		if handleErr != nil {
			c.log.WarnContext(
				ctx,
				"Kafka message processing failed",
				"topic", record.Topic,
				"partition", record.Partition,
				"offset", record.Offset,
				"attempts", 1+len(retryBackoffs),
				"err", handleErr,
			)

			if err := c.publishDLQ(ctx, record, handleErr); err != nil {
				c.client.AllowRebalance()
				return fmt.Errorf(
					"%s: handle topic %q partition %d offset %d: %w",
					op,
					record.Topic,
					record.Partition,
					record.Offset,
					err,
				)
			}
			movedToDLQ = true
			c.log.WarnContext(
				ctx,
				"Kafka message moved to DLQ",
				"topic", record.Topic,
				"dlq_topic", record.Topic+dlqTopicSuffix,
				"partition", record.Partition,
				"offset", record.Offset,
				"attempts", 1+len(retryBackoffs),
				"err", handleErr,
			)
		}

		if err := c.client.CommitRecords(ctx, record); err != nil {
			c.client.AllowRebalance()

			return fmt.Errorf(
				"%s: commit topic %q partition %d offset %d: %w",
				op,
				record.Topic,
				record.Partition,
				record.Offset,
				err,
			)
		}
		c.log.DebugContext(
			ctx,
			"Kafka offset committed",
			"topic", record.Topic,
			"partition", record.Partition,
			"offset", record.Offset,
			"moved_to_dlq", movedToDLQ,
		)
		c.client.AllowRebalance()
	}
}

func (c *Client) publishDLQ(ctx context.Context, record *kgo.Record, recordErr error) error {
	dlqRecord := &kgo.Record{
		Topic: record.Topic + dlqTopicSuffix,
		Key:   record.Key,
		Value: record.Value,
		Headers: []kgo.RecordHeader{
			{
				Key:   "x-original-topic",
				Value: []byte(record.Topic),
			},
			{
				Key:   "x-original-partition",
				Value: []byte(strconv.FormatInt(int64(record.Partition), 10)),
			},
			{
				Key:   "x-original-offset",
				Value: []byte(strconv.FormatInt(record.Offset, 10)),
			},
			{
				Key:   "x-failed-at",
				Value: []byte(time.Now().UTC().Format(time.RFC3339)),
			},
			{
				Key:   "x-attempts",
				Value: []byte(strconv.Itoa(1 + len(retryBackoffs))),
			},
			{
				Key:   "x-error",
				Value: []byte(recordErr.Error()),
			},
		},
	}

	results := c.producer.ProduceSync(ctx, dlqRecord)
	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("publish to DLQ: %w", err)
	}
	return nil
}

type retryObserver func(attempt int, backoff time.Duration, err error)

func handleWithRetry(
	ctx context.Context,
	handler broker.Handler,
	message broker.Message,
	observeRetry retryObserver,
) error {
	err := handler(ctx, message)
	if err == nil {
		return nil
	}

	for retryIndex, backoff := range retryBackoffs {
		if observeRetry != nil {
			observeRetry(retryIndex+2, backoff, err)
		}

		timer := time.NewTimer(backoff)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}

		retryErr := handler(ctx, message)
		if retryErr == nil {
			return nil
		}

		err = retryErr
	}

	return err
}

func (c *Client) Close() error {
	c.client.Close()
	return nil
}
