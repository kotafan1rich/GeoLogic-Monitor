package digitalspb

import "errors"

var (
	ErrInvalidHTTPClient = errors.New("HTTP client is nil")
	ErrInvalidURL        = errors.New("URL is nil or malformed")
	ErrInvalidURLMap     = errors.New("URLs map is empty or malformed")
	ErrBuildRequest      = errors.New("failed to build request")
	ErrCreateRequest     = errors.New("failed to create request")
	ErrReadResponse      = errors.New("failed to read response")
	ErrUnexpectedStatus  = errors.New("unexpected response status")
	ErrUnmarshalData     = errors.New("failed to unmarshal data")
	ErrInvalidFileName   = errors.New("invalid file name")
	ErrOpenFile          = errors.New("failed to open file")
	ErrDecodeData        = errors.New("failed to decode data")
)
