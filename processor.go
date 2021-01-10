package xdoc

import (
	"encoding/json"
	"encoding/xml"
)

func FromJSON(data []byte) (n *Node, err error) {
	n = NewNode()
	err = json.Unmarshal(data, n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func ToJSON(n *Node) ([]byte, error) {
	return json.Marshal(n)
}

func FromXML(data []byte) (n *Node, err error) {
	n = NewNode()
	err = xml.Unmarshal(data, n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func ToXML(n *Node) ([]byte, error) {
	return xml.Marshal(n)
}
