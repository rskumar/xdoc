package xdoc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/cockroachdb/errors"
	"strings"
)

var (
	MaxIteration      = 1000
	MaxRecursionDepth = 1000

	SkipLevel              = errors.New("skip level")
	StopWalk               = errors.New("stop walk")
	ErrMaxIter             = errors.New("max iteration reached")
	ErrMaxRecursion        = errors.New("max recursion depth reached")
	ErrInconsistentTree    = errors.New("inconsistent tree")
	ErrInvalidNode         = errors.New("invalid node")
	ErrNodeAlreadyAttached = errors.New("node already attached")
	ErrNodeNotFound        = errors.New("node not found")
	ErrParentNodeNotFound  = errors.New("parent node not found")
	ErrChildNotAllowed     = errors.New("child not allowed")
)

/*type NodeProps map[string]string

type NodeAccessorFn func() *Node

type NodeAccessor interface {
	bindAccessor(NodeAccessorFn)
}*/

type ChildContainer interface {
	append(*Node)
	GetChildren() []*Node
}

type Children []*Node

func (c *Children) append(n *Node) {
	*c = append(*c, n)
}

func (c *Children) GetChildren() []*Node {
	return *c
}

type Node struct {
	parent *Node
	Element
}

func (n Node) IsRoot() bool {
	return n.parent == nil
}

func (n Node) HasParent() bool {
	return n.parent != nil
}

func (n Node) Parent() *Node {
	return n.parent
}

func (n Node) NumChild() int {
	if n.Element != nil {
		// if element has child container embedded
		if container, ok := n.Element.(ChildContainer); ok {
			return len(container.GetChildren())
		}
	}
	return 0
}

func (n Node) FirstChild() *Node {
	if n.Element != nil {
		// if element has child container embedded
		if container, ok := n.Element.(ChildContainer); ok {
			children := container.GetChildren()
			if len(children) > 0 {
				return children[0]
			}
		}
	}
	return nil
}

func (n Node) LastChild() *Node {
	if n.Element != nil {
		// if element has child container embedded
		if container, ok := n.Element.(ChildContainer); ok {
			children := container.GetChildren()
			if len(children) > 0 {
				return children[len(children)-1]
			}
		}
	}
	return nil
}

func (n *Node) AppendChild(child ...*Node) error {
	if len(child) == 0 {
		return nil
	}

	for _, c := range child {
		err := n.appendSingleChild(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) appendSingleChild(c *Node) error {
	if n.Element == nil {
		return ErrInvalidNode
	}

	if !n.Element.CanContain(c.Element) {
		return ErrChildNotAllowed
	}

	// if element has child container embedded
	if container, ok := n.Element.(ChildContainer); ok {
		container.append(c)
	} else {
		return ErrChildNotAllowed
	}
	return nil
}

func (n *Node) unlink() {
	n.parent = nil
}

func (n *Node) GetElement() Element {
	if n.Element != nil {
		return n.Element
	}
	return nil
}

func (n *Node) IsLeaf() bool {
	if n.Element != nil {
		// if element has child container embedded
		if _, ok := n.Element.(ChildContainer); ok {
			return false
		}
	}
	return true
}

type typedNode struct {
	Type string `json:"type,omitempty"`
	Element
}

/*type nodeWithoutChildren struct {
	Type string `json:"type,omitempty"`
	Element
}*/

func (n *Node) MarshalJSON() ([]byte, error) {
	if n.Element == nil {
		return json.Marshal(nil)
	}

	var jElem []byte
	var err error
	if em, ok := n.Element.(json.Marshaler); ok {
		jElem, err = em.MarshalJSON()
		if err != nil {
			return nil, err
		}
	} else {
		jElem, err = json.Marshal(n.Element)
		if err != nil {
			return nil, err
		}
		if !InStrList(n.Element.GetType(), []string{"", "text"}) {
			var tJson []byte
			tJson, err = json.Marshal(&struct {
				Type string `json:"type,omitempty"`
			}{n.Element.GetType()})
			if err != nil {
				return nil, err
			}
			buf := &bytes.Buffer{}
			_, err = buf.Write(tJson[:len(tJson)-1])
			if err != nil {
				return nil, err
			}
			_, err = buf.Write([]byte(", "))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(jElem[1:])
			if err != nil {
				return nil, err
			}
			jElem = buf.Bytes()
		}
	}

	//jElem, err := json.Marshal(n.Element)
	/*v := make(map[string]interface{})
	err := mapstructure.Decode(n.Element, &v)
	if err != nil {
		return nil, err
	}
	typ := n.Element.GetType()
	if typ != "text" {
		v["type"] = typ
	}*/
	//v := typedNode{Element: n.Element, Type: n.Element.GetType()}
	//return json.Marshal(v)
	return jElem, nil
}

func (n *Node) UnmarshalJSON(data []byte) error {
	var err error
	var vt = &struct {
		Type string  `json:"type"`
		Text *string `json:"text"`
	}{}
	err = json.Unmarshal(data, vt)
	var nodeTyp = vt.Type
	if strings.TrimSpace(nodeTyp) == "" && vt.Text != nil {
		// its text type
		nodeTyp = TypeText
	}

	var elem Element
	switch nodeTyp {
	case TypePage:
		elem = NewPage()
	case TypeTitle:
		elem = NewTitle()
	case TypePara:
		elem = NewTitle()
	case TypeText:
		elem = NewText()
	default:
	}
	err = json.Unmarshal(data, elem)
	if err != nil {
		return err
	}
	n.Element = elem
	// If this element contains children, set all's parent to self
	if cc, ok := n.Element.(ChildContainer); ok {
		c := (Children)(cc.GetChildren())
		for _, cNode := range c {
			cNode.parent = n // I am your parent
		}
	}
	return nil
}

func (n *Node) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Is element handling its own XMl marshalling?
	if em, ok := n.Element.(xml.Marshaler); ok {
		err := em.MarshalXML(e, start)
		if err != nil {
			return err
		}
	} else {
		st := xml.StartElement{
			Name: xml.Name{
				Space: "",
				Local: n.Element.GetType(),
			},
			Attr: nil,
		}
		err := e.EncodeElement(n.Element, st)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fmt.Println("Start: ", start.Name.Local)
	nodeTyp := start.Name.Local
	var elem Element

	switch nodeTyp {
	case TypePage:
		elem = NewPage()
	case TypeTitle:
		elem = NewTitle()
	case TypePara:
		elem = NewTitle()
	case TypeText:
		elem = NewText()
	default:
	}
	err := d.DecodeElement(elem, &start)
	if err != nil {
		return err
	}

	// If this element contains children, set all's parent to self
	if cc, ok := elem.(ChildContainer); ok {
		c := (Children)(cc.GetChildren())
		for _, cNode := range c {
			cNode.parent = n // I am your parent
		}
	}

	n.Element = elem
	return nil
}

func NewNode() *Node {
	n := &Node{}
	return n
}

func NewElementNode(el Element) *Node {
	if isNilValue(el) {
		panic("nil element unsupported")
	}
	n := &Node{
		Element: el,
	}
	return n
}

/*type WalkFn func(node *Node, depth int, index int) (err error)

func (n *Node) walkInner(start *Node, fn WalkFn, depth int) (err error) {
	if depth > MaxRecursionDepth {
		return ErrMaxRecursion
	}
	currNode := n
	iterCount := 0

LOOP:
	for {
		//fmt.Println("For/start: Iter/depth", iterCount, depth)
		if iterCount > MaxIteration {
			return ErrMaxIter
		}

		if currNode == nil {
			//fmt.Println("For/break-loop: nil node found")
			break LOOP
		}

		// call fn
		err = fn(n, depth, iterCount)
		if err != nil {
			if err == SkipLevel {
				break LOOP
			} else if err == StopWalk {
				break LOOP
			} else {
				return err
			}
		}
		iterCount++
		err = n.walkInner(currNode.firstChild, fn, depth+1)
		if err == StopWalk {
			return err
		}

		currNode = currNode.next
	}
	return err
}*/

/*func (n *Node) Walk(fn WalkFn) error {
	err := n.walkInner(n, fn, 0)
	if err != SkipLevel && err != StopWalk {
		return err
	}
	return nil
}*/

/*func (n *Node) NextSibling() *Node {
	if n.parent == nil {
		return nil
	}
	foundSelf := false
	for _, c := range n.parent.Children {
		if foundSelf {
			return c
		} else if c == n {
			foundSelf = true
		}
	}
	return nil
}

func (n *Node) PrevSibling() *Node {
	if n.parent == nil {
		return nil
	}
	var lastNode *Node = nil
	for _, c := range n.parent.Children {
		if c == n {
			return lastNode
		} else {
			lastNode = c
		}
	}
	return nil
}

func (n *Node) LastSibling() *Node {
	return nil // todo: implement
}

func (n *Node) AddSiblingAfter(sibling *Node) (*Node, error) {
	return nil, nil // todo: implement
}

func (n *Node) Contains(v *Node) (found bool, at int) {
	for idx, node := range n.Children {
		if node == v {
			return true, idx
		}
	}
	return false, -1
}

type WalkChildrenFn func(n *Node, index int) error // fn should return StopWalk error to stop iteration
// TODO: walk and walkchildren can use same WalkChildrenFn, where depth in case of walkChildren will always be 0
func (n *Node) WalkChildren(fn WalkChildrenFn) error {
	var err error
	var nextNode = n.FirstChild()
	index := 0
	for !isNilValue(nextNode) {
		if index > MaxIteration {
			panic(ErrMaxIter)
		}
		e := fn(nextNode, index)
		if e == StopWalk {
			break
		} else if e != nil {
			err = e
			break
		}
		index++
		nextNode = nextNode.NextSibling()
	}
	return err
}

func (n *Node) RemoveChild(child *Node) error {
	found, at := n.Contains(child)
	if !found {
		return ErrNodeNotFound
	}
	copy(n.Children[at:], n.Children[at+1:])
	n.Children[len(n.Children)-1] = nil
	n.Children = n.Children[:len(n.Children)-1]
	return nil
}

func (n *Node) Remove() error {
	if n.parent == nil {
		return nil
	} else {
		return n.parent.RemoveChild(n)
	}
}
*/
