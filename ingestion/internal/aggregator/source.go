package aggregator

import (
	"fmt"
	"net/url"
	"strings"
)

type Kind uint8

const (
	KindUnknown Kind = iota
	KindURL
	KindFile
)

func (k Kind) String() string {
	switch k {
	case KindURL:
		return "url"
	case KindFile:
		return "file"
	default:
		return "unknown"
	}
}

type Source interface {
	Kind() Kind
	Validate() error
	String() string
}

var (
	_ Source = URL{}
	_ Source = File{}
)

type URL struct {
	u *url.URL
}

func NewURL(raw string) (URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return URL{}, fmt.Errorf("%w [%s]: %v", ErrInvalidURL, raw, err)
	}

	src := URL{u: parsed}
	if err := src.Validate(); err != nil {
		return URL{}, err
	}

	return src, nil
}

func NewURLFrom(u *url.URL) URL {
	return URL{u: u}
}

func (s URL) URL() *url.URL {
	return s.u
}

func (s URL) Kind() Kind {
	return KindURL
}

func (s URL) Validate() error {
	if s.u == nil || s.u.Scheme == "" || s.u.Host == "" {
		return fmt.Errorf("%w [%s]", ErrInvalidURL, s.String())
	}

	return nil
}

func (s URL) String() string {
	if s.u == nil {
		return ""
	}

	return s.u.String()
}

type File struct {
	path string
}

func NewFile(path string) File {
	return File{path: strings.TrimSpace(path)}
}

func (s File) Path() string {
	return s.path
}

func (s File) Kind() Kind {
	return KindFile
}

func (s File) Validate() error {
	if s.path == "" {
		return ErrInvalidFilePath
	}

	return nil
}

func (s File) String() string {
	return s.path
}
