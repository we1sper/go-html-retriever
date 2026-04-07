package htmlretriever

import (
	"strings"

	"golang.org/x/net/html"
)

func RetrieveAttrs(node *html.Node, visitor func(k, v string) bool) {
	if node == nil || node.Type != html.ElementNode {
		return
	}

	for _, attr := range node.Attr {
		if !visitor(attr.Key, attr.Val) {
			return
		}
	}
}

func GetAttr(node *html.Node, attr string) (value string, present bool) {
	RetrieveAttrs(node, func(k, v string) bool {
		if k == attr {
			value, present = v, true
		}
		return !present
	})
	return
}

func GetClasses(node *html.Node) []string {
	_class, ok := GetAttr(node, "class")
	if ok {
		return strings.Split(_class, " ")
	}
	return []string{}
}

func IsTag(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func HasId(node *html.Node, id string) bool {
	_id, ok := GetAttr(node, "id")
	return ok && _id == id
}

func HasAttr(node *html.Node, attr string) bool {
	_, ok := GetAttr(node, attr)
	return ok
}

func HasClasses(node *html.Node, classes ...string) bool {
	_classes := GetClasses(node)
	lookup := make(map[string]struct{}, len(_classes))
	for _, class := range _classes {
		lookup[class] = struct{}{}
	}
	for _, class := range classes {
		if _, ok := lookup[class]; !ok {
			return false
		}
	}
	return true
}

func DeepFirstSearch(node *html.Node, visitor func(n *html.Node) bool) {
	if node != nil {
		_ = deepFirstSearch(node, visitor)
	}
}

func deepFirstSearch(node *html.Node, visitor func(n *html.Node) bool) bool {
	if !visitor(node) {
		return false
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if !deepFirstSearch(child, visitor) {
			return false
		}
	}
	return true
}
