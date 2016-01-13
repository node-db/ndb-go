package ndb

import (
	"strings"
)

type Node struct {
	name string
	values map[string][]string
	children []*Node
}

func (node *Node) SetName(name string) {
	node.name = name
}

func (node *Node) GetName() string {
	return node.name
}

func (node *Node) GetValues() map[string][]string {
	return node.values
}

func (node *Node) FindChildByName(name string) []*Node {
	if name == "" {
		return nil
	}
	children := node.children
	if children != nil {
		matchedChildren := []*Node{}
		for _, child := range children {
			if child != nil && child.name == name {
				matchedChildren = append(matchedChildren, child)
			}
		}
		return matchedChildren
	}
	return nil
}

func (node *Node) AddChild(child *Node) {
	node.children = append(node.children, child)
}

func (node *Node) AddChildren(children []*Node) {
	for _, child := range children {
		node.AddChild(child)
	}
}

func (node *Node) GetChileren() []*Node {
	return node.children
}

func (node *Node) GetValue(key string) []string {
	value := node.values
	return value[key]
}

func (node *Node) DeleteValue(key string) {
	delete(node.values, key)
}

func (node *Node) ClearValue() {
	for key, _ := range node.values {
		delete(node.values, key)
	}
}

func (node *Node) GetValueString(key string) string {
	value := node.GetValue(key)
	return strings.Join(value, ",")
}

func (node *Node) SetValue(key string, value []string) {
	nodeValue := node.values
	if (nodeValue == nil) {
		nodeValue = make(map[string][]string)
	}
	nodeValue[key] = value
	node.values = nodeValue
}
