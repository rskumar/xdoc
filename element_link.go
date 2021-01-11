package xdoc

type Link struct {
	Href     string `json:"href,omitempty" xml:"href,attr,omitempty"`
	Title    string `json:"title,omitempty" xml:"title,attr,omitempty"`
	Children `json:"children" xml:",any"`
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

func linkCtor() Element {
	return NewLink()
}
