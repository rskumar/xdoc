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
	//getProps() map[string]string
}

type TextValuer interface {
	GetText() string
}

// Types defined

// Page node
type Page struct {
	Version  int    `json:"version,string" xml:"version,attr"`
	Lang     string `json:"lang" xml:"lang,attr"`
	Children `json:"children" xml:",any"`
}

func (p *Page) GetType() string {
	return TypePage
}

func (p *Page) CanContain(e Element) bool {
	if isNilValue(e) {
		return false
	}
	return IsValidParent(p.GetType(), e.GetType())
}

/*type pageStruct struct {
	Type     string    `json:"type"`
	Version  string    `json:"version"`
	Lang     string    `json:"lang"`
	Children []Element `json:"children"`
}

func (p *Page) MarshalJSON() ([]byte, error) {
	v := pageStruct{
		p.GetType(),
		strconv.Itoa(p.Version),
		p.Lang,
		p.node.ChildElements(),
	}
	return json.Marshal(v)
}

func (p *Page) UnmarshalJSON(data []byte) error {
	v := &pageStruct{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if v.Type != TypePage {
		return ErrInvalidElementType
	}
	ver, err := strconv.Atoi(v.Version)
	if err != nil {
		return err
	}
	p.Version = ver
	p.Lang = v.Lang
	err = p.AppendChildElement(v.Children...)
	if err != nil {
		return err
	}
	return nil
}*/

func NewPage() *Page {
	v := &Page{}
	return v
}

// Title node
type Title struct {
	Children `json:"children"`
}

func (t *Title) GetType() string {
	return TypeTitle
}

func (t *Title) CanContain(e Element) bool {
	if isNilValue(e) {
		return false
	}
	return IsValidParent(t.GetType(), e.GetType())
}

/*type titleStruct struct {
	Type     string    `json:"type"`
	Children []Element `json:"children"`
}

func (t *Title) MarshalJSON() ([]byte, error) {
	v := titleStruct{
		t.GetType(),
		t.node.ChildElements(),
	}
	return json.Marshal(v)
}

func (t *Title) UnmarshalJSON(data []byte) error {
	v := &titleStruct{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}


	p.Version = ver
	p.Lang = v.Lang
	err = p.AppendChildElement(v.Children...)
	if err != nil {
		return err
	}
	return nil
}*/

func NewTitle() *Title {
	v := &Title{}
	return v
}

// Paragraph node
type Paragraph struct {
	Children `json:"children" xml:",any"`
}

func (p *Paragraph) GetType() string {
	return TypePara
}

func (p *Paragraph) CanContain(e Element) bool {
	if isNilValue(e) {
		return false
	}
	return IsValidParent(p.GetType(), e.GetType())
}

func NewParagraph() *Paragraph {
	v := &Paragraph{}
	return v
}

type Text struct {
	Text          string `json:"text" xml:",chardata"`
	Code          bool   `json:"code,omitempty" xml:"code,attr,omitempty"`
	Bold          bool   `json:"bold,omitempty" xml:"bold,attr,omitempty"`
	Italics       bool   `json:"italics,omitempty" xml:"italics,attr,omitempty"`
	Underline     bool   `json:"underline,omitempty" xml:"underline,attr,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty" xml:"strikethrough,attr,omitempty"`
}

func (t *Text) GetType() string {
	return TypeText
}

func (t *Text) CanContain(e Element) bool {
	if isNilValue(e) {
		return false
	}
	return IsValidParent(t.GetType(), e.GetType())
}

func NewText() *Text {
	v := &Text{}
	return v
}

type Link struct {
	Href  string
	Title string
	Children
}

func (Link) GetType() string {
	return TypeLink
}

func (l *Link) CanContain(e Element) bool {
	if isNilValue(e) {
		return false
	}
	return IsValidParent(l.GetType(), e.GetType())
}

func NewLink() *Link {
	return &Link{}
}
