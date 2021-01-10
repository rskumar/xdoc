package xdoc

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

func (p *Page) WithVersion(v int) *Page {
	p.Version = v
	return p
}

func (p *Page) WithLang(l string) *Page {
	p.Lang = l
	return p
}

func NewPage() *Page {
	v := &Page{}
	return v
}
