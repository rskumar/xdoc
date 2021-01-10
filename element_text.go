package xdoc

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

func (t *Text) WithText(txt string) *Text {
	t.Text = txt
	return t
}

func (t *Text) WithBold(bold bool) *Text {
	t.Bold = bold
	return t
}

func NewText() *Text {
	v := &Text{}
	return v
}
