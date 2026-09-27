package producer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v7"
	"github.com/twmb/franz-go/pkg/kgo"
)

const service = "ingestion-service"

type Producer struct {
	client         *kgo.Client
	log            *slog.Logger
	flushTimeout   time.Duration
	newBackOff     func() backoff.BackOff
	attemptTimeout time.Duration
	maxRetries     uint
}

func New(
	address []string,
	log *slog.Logger,
	flushTimeout time.Duration,
	attemptTimeout time.Duration,
	maxRetries uint,
	batchSize int32,
	bufferMaxMsg int,
	bufferMaxBytes int,
) (*Producer, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	client, err := kgo.NewClient([]kgo.Opt{
		kgo.SeedBrokers(address...),
		kgo.AllowAutoTopicCreation(),
		kgo.ProducerBatchMaxBytes(batchSize),
		kgo.MaxBufferedRecords(bufferMaxMsg),
		kgo.MaxBufferedBytes(bufferMaxBytes),
	}...)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client:         client,
		log:            log,
		flushTimeout:   flushTimeout,
		attemptTimeout: attemptTimeout,
		newBackOff:     func() backoff.BackOff { return backoff.NewExponentialBackOff() },
		maxRetries:     maxRetries,
	}, nil
}

func MustNew(
	address []string,
	log *slog.Logger,
	flushTimeout time.Duration,
	attemptTimeout time.Duration,
	maxRetries uint,
	batchSize int32,
	bufferMaxMsg int,
	bufferMaxBytes int,
) *Producer {
	const op = "producer.MustNewProducer"

	p, err := New(
		address,
		log,
		flushTimeout,
		attemptTimeout,
		maxRetries,
		batchSize,
		bufferMaxMsg,
		bufferMaxBytes,
	)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to init producer: %v", op, err))
	}

	return p
}

func (p *Producer) Produce(ctx context.Context, msg, topic string) error {
	const op = "producer.Producer.Produce"

	_, err := backoff.Retry(
		ctx,
		func() (struct{}, error) {
			return p.attempt(ctx, msg, topic)
		},
		backoff.WithBackOff(p.newBackOff()),
		backoff.WithMaxElapsedTime(0),
		backoff.WithMaxTries(p.maxRetries),
		backoff.WithNotify(func(err error, d time.Duration) {
			p.log.WarnContext(ctx, "failed to send message",
				slog.Any("error", err), slog.Duration("next_retry_in", d))
		}),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Producer) attempt(ctx context.Context, msg, topic string) (struct{}, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, p.attemptTimeout)
	defer cancel()

	err := p.client.ProduceSync(attemptCtx, &kgo.Record{
		Topic: topic,
		Value: []byte(msg),
		Key:   nil,
		Headers: []kgo.RecordHeader{
			{Key: "service", Value: []byte(service)},
		},
	}).FirstErr()
	if err == nil {
		return struct{}{}, nil
	}

	p.log.Info("err", slog.Any("error", err))

	if ctx.Err() != nil {
		return struct{}{}, backoff.Permanent(fmt.Errorf("%w: %v", ErrDeliveryMsg, err))
	}
	if isRetriable(err) {
		return struct{}{}, fmt.Errorf("%w: %v", ErrDeliveryMsg, err)
	}
	return struct{}{}, backoff.Permanent(fmt.Errorf("%w: %v", ErrDeliveryMsg, err))
}

func (p *Producer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.flushTimeout)
	defer cancel()

	err := p.client.Flush(ctx)
	if err != nil {
		p.log.Warn("failed to flush producer", slog.Any("error", err))
		return err
	}
	p.client.Close()

	return nil
}
