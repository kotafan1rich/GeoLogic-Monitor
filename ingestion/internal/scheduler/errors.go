package scheduler

import "errors"

var (
	ErrInvalidLogger     = errors.New("logger is nil")
	ErrInitCronScheduler = errors.New("failed to init cron scheduler")
	ErrRegisterCronJob   = errors.New("failed to register cron job")
	ErrShutdownCron      = errors.New("failed to shutdown cron scheduler")
)
