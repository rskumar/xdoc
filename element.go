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

func NewElementForType(typ string) (Element, error) {
	nodeInfo, ok := schema[typ]
	if !ok {
		return nil, errors.Wrapf(ErrInvalidElementType, "element type '%s' not registered")
	}
	if nodeInfo.Ctor == nil {
		return nil, errors.Newf("element constructor for '%s' not registered", nodeInfo.Type)
	}

	elem := nodeInfo.Ctor()
	return elem, nil
}
