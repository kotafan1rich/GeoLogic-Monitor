package infra

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const defaultBurst = 1

type RateLimit struct {
	RPS   float64
	Burst int
}

func (r RateLimit) Enabled() bool {
	return r.RPS > 0
}

type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	burst    float64
	tokens   float64
	last     time.Time
}

func NewLimiter(rl RateLimit) *Limiter {
	if !rl.Enabled() {
		return &Limiter{}
	}

	burst := rl.Burst
	if burst <= 0 {
		burst = defaultBurst
	}

	return &Limiter{
		interval: time.Duration(float64(time.Second) / rl.RPS),
		burst:    float64(burst),
		tokens:   float64(burst),
		last:     time.Now(),
	}
}

func (l *Limiter) Wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	delay := l.reserve()
	if delay <= 0 {
		return nil
	}

	t := time.NewTimer(delay)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// reserve забирает один токен и сообщает, сколько нужно подождать до его
// фактического появления. Токены уходят в минус, поэтому очередь из запросов
// равномерно размазывается по времени вместо всплеска.
func (l *Limiter) reserve() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.interval <= 0 {
		return 0
	}

	now := time.Now()

	if elapsed := now.Sub(l.last); elapsed > 0 {
		l.tokens += float64(elapsed) / float64(l.interval)
		if l.tokens > l.burst {
			l.tokens = l.burst
		}
	}

	l.last = now
	l.tokens--

	if l.tokens >= 0 {
		return 0
	}

	return time.Duration(-l.tokens * float64(l.interval))
}

type rateLimitTransport struct {
	limiter *Limiter
	next    http.RoundTripper
}

func (t *rateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(r.Context()); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRateLimit, err)
	}

	return t.next.RoundTrip(r)
}
