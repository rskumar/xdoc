package xdoc

// Title node
type Title struct {
	Children `json:"children" xml:",any"`
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

func (t *Title) WithText(s string) *Title {
	tNode := NewElementNode(NewText().WithText(s))
	t.append(tNode)
	return t
}

func NewTitle() *Title {
	v := &Title{}
	return v
}

func titleCtor() Element {
	return NewTitle()
}
