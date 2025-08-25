package proxy

import (
	"bytes"
	"fmt"

	"proxy-enhancer-ultra/internal/models"
	"proxy-enhancer-ultra/pkg/logger"

	"golang.org/x/net/html"
	"gorm.io/gorm"
)

// HTMLInjector HTML内容注入器
type HTMLInjector struct {
	db             *gorm.DB
	logger         logger.Logger
	htmlParser     *HTMLParser
	assetsProvider *AssetsProvider
}

// NewHTMLInjector 创建新的HTML注入器
func NewHTMLInjector(db *gorm.DB, logger logger.Logger) *HTMLInjector {
	return &HTMLInjector{
		db:             db,
		logger:         logger,
		htmlParser:     NewHTMLParser(),
		assetsProvider: NewAssetsProvider(),
	}
}

// ProcessHTML 处理HTML内容
func (hi *HTMLInjector) ProcessHTML(body []byte, proxyConfig *models.ProxyConfig, urlRewriter *URLRewriter) ([]byte, error) {
	// 解析HTML
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return body, err
	}

	// 注入自定义内容
	hi.injectCustomContent(doc, proxyConfig)

	// 修改链接和资源路径
	urlRewriter.RewriteURLs(doc, proxyConfig)

	// 将修改后的HTML转换回字节
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return body, err
	}

	return buf.Bytes(), nil
}

// injectCustomContent 注入自定义内容
func (hi *HTMLInjector) injectCustomContent(doc *html.Node, proxyConfig *models.ProxyConfig) {
	// 查找head标签
	headNode := hi.htmlParser.FindNode(doc, "head")
	if headNode == nil {
		return
	}

	// 注入基础样式和脚本
	hi.injectBaseAssets(headNode, proxyConfig)

	// 查找body标签
	bodyNode := hi.htmlParser.FindNode(doc, "body")
	if bodyNode == nil {
		return
	}

	// 注入弹窗和交互功能
	hi.injectInteractiveElements(bodyNode, proxyConfig)
}

// injectBaseAssets 注入基础资源
func (hi *HTMLInjector) injectBaseAssets(headNode *html.Node, proxyConfig *models.ProxyConfig) {
	// 注入CSS样式
	styleNode := &html.Node{
		Type: html.ElementNode,
		Data: "style",
		Attr: []html.Attribute{{Key: "type", Val: "text/css"}},
	}
	styleContent := &html.Node{
		Type: html.TextNode,
		Data: hi.assetsProvider.GetCSS(),
	}
	styleNode.AppendChild(styleContent)
	headNode.AppendChild(styleNode)

	// 注入JavaScript
	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{{Key: "type", Val: "text/javascript"}},
	}
	scriptContent := &html.Node{
		Type: html.TextNode,
		Data: hi.assetsProvider.GetJavaScript(proxyConfig),
	}
	scriptNode.AppendChild(scriptContent)
	headNode.AppendChild(scriptNode)
}

// injectInteractiveElements 注入交互元素
func (hi *HTMLInjector) injectInteractiveElements(bodyNode *html.Node, proxyConfig *models.ProxyConfig) {
	// 如果没有数据库连接，跳过弹窗注入
	if hi.db == nil {
		return
	}

	// 获取该代理配置的弹窗规则
	var popups []models.Popup
	hi.db.Where("proxy_config_id = ? AND enabled = ?", proxyConfig.ID, true).Find(&popups)

	for _, popup := range popups {
		// 创建弹窗容器
		popupContainer := hi.createPopupElement(&popup)
		bodyNode.AppendChild(popupContainer)
	}
}

// createPopupElement 创建弹窗元素
func (hi *HTMLInjector) createPopupElement(popup *models.Popup) *html.Node {
	// 创建弹窗容器
	container := &html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			{Key: "id", Val: fmt.Sprintf("popup-%d", popup.ID)},
			{Key: "class", Val: "proxy-popup"},
			{Key: "style", Val: "display: none;"},
		},
	}

	// 添加弹窗内容
	content := &html.Node{
		Type: html.RawNode,
		Data: popup.Content,
	}
	container.AppendChild(content)

	return container
}
