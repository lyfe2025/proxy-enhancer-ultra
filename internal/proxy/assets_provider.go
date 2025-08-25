package proxy

import (
	"fmt"

	"proxy-enhancer-ultra/internal/models"
)

// AssetsProvider 静态资源提供器
type AssetsProvider struct{}

// NewAssetsProvider 创建新的静态资源提供器
func NewAssetsProvider() *AssetsProvider {
	return &AssetsProvider{}
}

// GetCSS 获取注入的CSS样式
func (ap *AssetsProvider) GetCSS() string {
	return `
/* 代理增强器样式 */
.proxy-popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: white;
    border: 1px solid #ccc;
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0,0,0,0.3);
    z-index: 10000;
    max-width: 500px;
    max-height: 80vh;
    overflow-y: auto;
    padding: 20px;
}

.proxy-popup-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.5);
    z-index: 9999;
}

.proxy-popup-close {
    position: absolute;
    top: 10px;
    right: 15px;
    font-size: 24px;
    cursor: pointer;
    color: #999;
}

.proxy-popup-close:hover {
    color: #333;
}
`
}

// GetJavaScript 获取注入的JavaScript代码
func (ap *AssetsProvider) GetJavaScript(proxyConfig *models.ProxyConfig) string {
	return fmt.Sprintf(`
// 代理增强器脚本
(function() {
    'use strict';
    
    // 代理配置
    window.PROXY_CONFIG = {
        domain: '%s',
        configId: '%s'
    };
    
    // 弹窗管理器
    window.ProxyPopupManager = {
        show: function(popupId) {
            var popup = document.getElementById('popup-' + popupId);
            if (popup) {
                // 创建遮罩层
                var overlay = document.createElement('div');
                overlay.className = 'proxy-popup-overlay';
                overlay.onclick = function() {
                    ProxyPopupManager.hide(popupId);
                };
                document.body.appendChild(overlay);
                
                // 显示弹窗
                popup.style.display = 'block';
                
                // 添加关闭按钮
                var closeBtn = document.createElement('span');
                closeBtn.className = 'proxy-popup-close';
                closeBtn.innerHTML = '×';
                closeBtn.onclick = function() {
                    ProxyPopupManager.hide(popupId);
                };
                popup.appendChild(closeBtn);
            }
        },
        
        hide: function(popupId) {
            var popup = document.getElementById('popup-' + popupId);
            if (popup) {
                popup.style.display = 'none';
                
                // 移除遮罩层
                var overlay = document.querySelector('.proxy-popup-overlay');
                if (overlay) {
                    overlay.remove();
                }
                
                // 移除关闭按钮
                var closeBtn = popup.querySelector('.proxy-popup-close');
                if (closeBtn) {
                    closeBtn.remove();
                }
            }
        }
    };
    
    // 数据收集功能
    window.ProxyDataCollector = {
        collectFormData: function(formElement) {
            var formData = new FormData(formElement);
            var data = {};
            for (var pair of formData.entries()) {
                data[pair[0]] = pair[1];
            }
            return data;
        },
        
        submitData: function(data, endpoint) {
            fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            }).then(function(response) {
                return response.json();
            }).then(function(result) {
                console.log('Data submitted successfully:', result);
            }).catch(function(error) {
                console.error('Error submitting data:', error);
            });
        }
    };
    
    console.log('Proxy enhancer loaded for domain:', window.PROXY_CONFIG.domain);
})();
`, proxyConfig.ProxyDomain, proxyConfig.ID)
}
