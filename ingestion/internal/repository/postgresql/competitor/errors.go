package competitor

import (
	"errors"
)

var (
	ErrInsertData   = errors.New("failed to insert data")
	ErrUpdateStatus = errors.New("failed to update status")
)
