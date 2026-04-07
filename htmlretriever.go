package htmlretriever

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Node interface {
	GetElementById(id string) *EnhancedNode
	GetElementsByTagName(tag string) *EnhancedNodes
	GetElementsByClassName(classes ...string) *EnhancedNodes
}

type UnsafeNode struct {
	_node *EnhancedNode
}

func (u *UnsafeNode) GetAttr(attr string) string {
	value, _ := u._node.GetAttr(attr)
	return value
}

func (u *UnsafeNode) GetText() string {
	text, _ := u._node.GetText()
	return text
}

type EnhancedNode struct {
	_node *html.Node
}

func NewEnhancedNode(node *html.Node) *EnhancedNode {
	return &EnhancedNode{_node: node}
}

func From(reader io.Reader) (*EnhancedNode, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	return NewEnhancedNode(node), nil
}

func (n *EnhancedNode) GetElementById(id string) *EnhancedNode {
	next := &EnhancedNode{}

	if n._node != nil {
		DeepFirstSearch(n._node, func(node *html.Node) bool {
			if HasId(node, id) {
				next._node = node
			}
			return next._node == nil
		})
	}

	return next
}

func (n *EnhancedNode) GetElementsByTagName(tag string) *EnhancedNodes {
	next := &EnhancedNodes{_nodes: make([]*EnhancedNode, 0)}

	DeepFirstSearch(n._node, func(node *html.Node) bool {
		if node != n._node && IsTag(node, tag) {
			next._nodes = append(next._nodes, NewEnhancedNode(node))
		}
		return true
	})

	return next
}

func (n *EnhancedNode) GetElementsByClassName(classes ...string) *EnhancedNodes {
	next := &EnhancedNodes{_nodes: make([]*EnhancedNode, 0)}

	DeepFirstSearch(n._node, func(node *html.Node) bool {
		if node != n._node && HasClasses(node, classes...) {
			next._nodes = append(next._nodes, NewEnhancedNode(node))
		}
		return true
	})

	return next
}

func (n *EnhancedNode) GetText() (text string, present bool) {
	if n._node != nil {
		for child := n._node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				text, present = strings.TrimSpace(child.Data), true
				break
			}
		}
	}
	return
}

func (n *EnhancedNode) GetTexts() []string {
	texts := make([]string, 0)

	if n._node != nil {
		for child := n._node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				texts = append(texts, strings.TrimSpace(child.Data))
			}
		}
	}

	return texts
}

func (n *EnhancedNode) ExtractTexts() []string {
	texts := make([]string, 0)

	DeepFirstSearch(n._node, func(node *html.Node) bool {
		if node.Type == html.TextNode {
			texts = append(texts, strings.TrimSpace(node.Data))
		}
		return true
	})

	return texts
}

func (n *EnhancedNode) RetrieveAttrs(processor func(k, v string) bool) {
	RetrieveAttrs(n._node, processor)
}

func (n *EnhancedNode) GetAttr(attr string) (value string, present bool) {
	return GetAttr(n._node, attr)
}

func (n *EnhancedNode) IsTag(tag string) bool {
	return IsTag(n._node, tag)
}

func (n *EnhancedNode) HasId(id string) bool {
	return HasId(n._node, id)
}

func (n *EnhancedNode) HasAttr(attr string) bool {
	return HasAttr(n._node, attr)
}

func (n *EnhancedNode) HasClasses(classes ...string) bool {
	return HasClasses(n._node, classes...)
}

func (n *EnhancedNode) Present() bool {
	return n._node != nil
}

func (n *EnhancedNode) Unsafe() *UnsafeNode {
	return &UnsafeNode{_node: n}
}

func (n *EnhancedNode) Copy() *EnhancedNode {
	return NewEnhancedNode(n._node)
}

func (n *EnhancedNode) Raw() *html.Node {
	return n._node
}

type EnhancedNodes struct {
	_nodes []*EnhancedNode
}

func (n *EnhancedNodes) GetElementById(id string) *EnhancedNode {
	next := &EnhancedNode{}

	for _, node := range n._nodes {
		candidate := node.GetElementById(id)
		if candidate.Present() {
			next._node = candidate.Raw()
		}
	}

	return next
}

func (n *EnhancedNodes) GetElementsByTagName(tag string) *EnhancedNodes {
	next := &EnhancedNodes{_nodes: make([]*EnhancedNode, 0)}

	for _, node := range n._nodes {
		node.GetElementsByTagName(tag).ForEach(func(n *EnhancedNode) {
			next._nodes = append(next._nodes, n.Copy())
		})
	}

	return next
}

func (n *EnhancedNodes) GetElementsByClassName(classes ...string) *EnhancedNodes {
	next := &EnhancedNodes{_nodes: make([]*EnhancedNode, 0)}

	for _, node := range n._nodes {
		node.GetElementsByClassName(classes...).ForEach(func(n *EnhancedNode) {
			next._nodes = append(next._nodes, n.Copy())
		})
	}

	return next
}

func (n *EnhancedNodes) ForEach(visitor func(n *EnhancedNode)) {
	for _, node := range n._nodes {
		visitor(node)
	}
}

func (n *EnhancedNodes) Len() int {
	return len(n._nodes)
}

func (n *EnhancedNodes) First() *EnhancedNode {
	return n.At(0)
}

func (n *EnhancedNodes) Last() *EnhancedNode {
	return n.At(-1)
}

func (n *EnhancedNodes) At(index int) *EnhancedNode {
	next := &EnhancedNode{}

	if index < 0 {
		index += n.Len()
	}

	if index >= 0 && index < n.Len() {
		next._node = n._nodes[index]._node
	}

	return next
}
