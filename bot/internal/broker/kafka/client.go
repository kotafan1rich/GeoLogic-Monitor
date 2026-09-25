package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/broker"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	client *kgo.Client
}

var retryBackoffs = [...]time.Duration{
	time.Second,
	2 * time.Second,
	4 * time.Second,
}

func New(brokers []string, topic string, groupID string) (*Client, error) {
	const op = "broker.kafka.New"

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
		client: client,
	}, nil
}

func (c *Client) Consume(ctx context.Context, handler broker.Handler) error {
	const op = "broker.kafka.Client.Consume"

	if handler == nil {
		return fmt.Errorf("%s: handler is required", op)
	}

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

		if err := handleWithRetry(ctx, handler, message); err != nil {
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
		c.client.AllowRebalance()
	}
}

func handleWithRetry(
	ctx context.Context,
	handler broker.Handler,
	message broker.Message,
) error {
	err := handler(ctx, message)
	if err != nil {
		return err
	}

	for _, backoff := range retryBackoffs {
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
