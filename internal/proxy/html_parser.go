package proxy

import (
	"golang.org/x/net/html"
)

// HTMLParser HTML解析器
type HTMLParser struct{}

// NewHTMLParser 创建新的HTML解析器
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// FindNode 查找指定标签的节点
func (hp *HTMLParser) FindNode(doc *html.Node, tagName string) *html.Node {
	var result *html.Node
	hp.WalkNodes(doc, func(node *html.Node) {
		if result == nil && node.Type == html.ElementNode && node.Data == tagName {
			result = node
		}
	})
	return result
}

// WalkNodes 遍历所有节点
func (hp *HTMLParser) WalkNodes(node *html.Node, fn func(*html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		hp.WalkNodes(child, fn)
	}
}
