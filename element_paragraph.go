package xdoc

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
