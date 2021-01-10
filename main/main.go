package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/rskumar/xdoc"
)

var PanicIfErr = xdoc.PanicIfErr

var data = `
<page>
	<title>This is example</title>
</page>
`

func main() {
	//xmlMarshal()
	xmlUnmarshal()
}

func xmlMarshal() {
	pageNode := getDoc()
	xb, err := xml.Marshal(pageNode)
	PanicIfErr(err)
	fmt.Println(string(xb))
}

var xData = []byte(`
<page version="1" lang="en"><title><text>This is title text</text></title><p><text>This is title text</text></p></page>
`)

func xmlUnmarshal() {
	n := xdoc.NewNode()
	err := xml.Unmarshal(xData, n)
	PanicIfErr(err)
	spew.Dump(n)
}

var jData = []byte(`{
	"type": "page",
	"version": 1,
	"lang": "en",
	"children": [
		{
			"type": "title",
			"children": [
				{
					"text": "This is title text"
				}
			]
		},
		{
			"type": "p",
			"children": [
				{
					"text": "This is title text"
				}
			]
		}
	]
}`)

func jsonUnmarshal() {
	n := xdoc.NewNode()
	err := json.Unmarshal(jData, n)
	PanicIfErr(err)
	spew.Dump(n)
}

func jsonMarshal() {

	pageNode := getDoc()

	// print json
	j, err := json.MarshalIndent(pageNode, "", "\t")
	PanicIfErr(err)
	fmt.Println(string(j))

	/*fmt.Println("\n")
	x, err := xml.Marshal(page)
	PanicIfErr(err)
	fmt.Println(string(x))*/
	spew.Dump(pageNode)
}

func getDoc() *xdoc.Node {
	page := xdoc.NewPage()
	page.Version = 1
	page.Lang = "en"
	pageNode := xdoc.NewElementNode(page)

	title := xdoc.NewTitle()
	titleNode := xdoc.NewElementNode(title)

	t1 := xdoc.NewText()
	t1.Text = "This is title text"
	t1Node := xdoc.NewElementNode(t1)
	err := titleNode.AppendChild(t1Node)
	PanicIfErr(err)

	err = pageNode.AppendChild(titleNode)
	PanicIfErr(err)

	para1 := xdoc.NewParagraph()
	para1Node := xdoc.NewElementNode(para1)

	t2 := xdoc.NewText()
	t2.Text = "This is title text"
	t2Node := xdoc.NewElementNode(t2)
	err = para1Node.AppendChild(t2Node)
	PanicIfErr(err)

	err = pageNode.AppendChild(para1Node)
	PanicIfErr(err)
	return pageNode
}
