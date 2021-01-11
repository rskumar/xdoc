package xdoc

import "github.com/cockroachdb/errors"

const (
	TypeNone       = ""
	TypePage       = "page"
	TypePara       = "p"
	TypeText       = "text"
	TypeLink       = "link"
	TypeBlockquote = "blockquote"
	TypeTitle      = "title"
	TypeH1         = "h1"
	TypeH2         = "h2"
	TypeH3         = "h3"
	TypeList       = "list"
	TypeImage      = "image"
	TypePanel      = "panel"

	// Inline
)

var AllNodes = []string{TypeNone, TypePage, TypePara, TypeText, TypeLink, TypeBlockquote, TypeTitle, TypeH1, TypeH2, TypeH3, TypeList}

func IsValidNode(n string) bool {
	for _, v := range AllNodes {
		if v == n {
			return true
		}
	}
	return false
}

type ctor func() Element

func nilCtor() Element {
	panic(errors.New("Not implemented"))
}

type NodeInfo struct {
	Type         string
	AllowedChild []string
	Ctor         ctor
}

var schema = map[string]NodeInfo{
	TypePage: {
		Type: TypePage,
		AllowedChild: []string{
			TypeTitle,
			TypeH1,
			TypeH2,
			TypeH3,
			TypeBlockquote,
			TypePara,
			TypeList,
			TypeImage,
		},
		Ctor: pageCtor,
	},
	TypeTitle: {
		Type: TypeTitle,
		AllowedChild: []string{
			TypeText,
		},
		Ctor: titleCtor,
	},
	TypeH1: {
		Type: TypeH1,
		AllowedChild: []string{
			TypeText,
			TypeLink,
		},
		Ctor: nilCtor,
	},
	TypeH2: {
		Type: TypeH2,
		AllowedChild: []string{
			TypeText,
			TypeLink,
		},
		Ctor: nilCtor,
	},
	TypeH3: {
		Type: TypeH3,
		AllowedChild: []string{
			TypeText,
			TypeLink,
		},
		Ctor: nilCtor,
	},
	TypeBlockquote: {
		Type: TypeBlockquote,
		AllowedChild: []string{
			TypeText,
		},
		Ctor: nilCtor,
	},
	TypePara: {
		Type: TypePara,
		AllowedChild: []string{
			TypeText,
			TypeLink,
		},
		Ctor: paragraphCtor,
	},
	TypeLink: {
		Type: TypeLink,
		AllowedChild: []string{
			TypeText,
		},
		Ctor: linkCtor,
	},
}

// GetValidParentsFor returns slice of parents which can contain `n` as child
func GetValidParentsFor(n string) (parents []string) {
	for _, v := range schema {
		for _, child := range v.AllowedChild {
			if child == n {
				parents = append(parents, v.Type)
				break
			}
		}
	}
	return parents
}

// IsValidParent returns true if parent is valid node to contain given child
func IsValidParent(parent, child string) bool {
	for _, p := range GetValidParentsFor(child) {
		if parent == p {
			return true
		}
	}
	return false
}
