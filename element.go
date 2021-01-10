package xdoc

import (
	"github.com/cockroachdb/errors"
)

var (
	ErrInvalidElementType = errors.New("invalid element type")
)

type Element interface {
	GetType() string
	CanContain(Element) bool
	// AllowedChild() []string	// should be simple fn than method, so not all element has to implement it
}

type TextValuer interface {
	GetText() string
}

/*type Validator interface {
	Validate(context.Context) error
}

type Normalizer interface {
	Normalize()
}*/

// Types defined
