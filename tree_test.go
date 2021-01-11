package xdoc

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TreeSimple2(t *testing.T) {
	elPage := NewElementNode(NewPage().WithVersion(1).WithLang("en")).
		WithChild(NewElementNode(NewTitle().WithText("Sample document"))).
		WithChild(
			NewElementNode(NewParagraph()).WithChild(
				NewElementNode(NewText().WithText("Paragraph text fragment 1")),
			).WithChild(
				NewElementNode(NewText().WithText("Paragraph text fragment 2").WithBold(true)),
			),
		)

	t.Log(spew.Sdump(elPage))
	j, err := json.MarshalIndent(elPage, "", "  ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(j))
}

func Test_ValidXML_Short(t *testing.T) {
	xData := `
<?xml version="1.0" encoding="UTF-8"?>
<page version="1" lang="en-us">
	<title>Some title text</title>
	<p>
		<text>Sample text fragment(1) in para </text>
		<text bold="true" italics="false" underline="true">with styled part.</text>
		<link href="#">
			<text>This is link content</text>
		</link>
	</p>
</page>
`
	n, err := FromXML([]byte(xData))
	assert.Nil(t, err, "error should be nil for valid xml data")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, n, "non nil *Node should be returned")

	// root node's parent should be nil
	assert.Nil(t, n.parent, "top node parent should be nil")
	assert.False(t, n.HasParent(), "root node should not have parent")

	assert.Equal(t, 2, n.NumChild(), "children count should be 2")

	// first node should have page element type
	assert.Equal(t, TypePage, n.Element.GetType(), "root element type should be page")

	// root node's child should be title and p
	rootChildTypes := []string{
		// NOTE: order of types MUST be as defined in XML
		TypeTitle, TypePara,
	}

	if rc, found := n.Element.(ChildContainer); found {
		for idx, c := range rc.GetChildren() {
			assert.Equalf(t, rootChildTypes[idx], c.Element.GetType(), "page's children element type not in order")
		}
	} else {
		t.Error("page should contain children")
	}

}
