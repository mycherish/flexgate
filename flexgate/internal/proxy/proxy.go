package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

func NewReverseProxy(target string, stripPrefix string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// ⭐ path rewrite 去掉前缀
		// 1. 剥离指定的前缀 (例如从 /api/users/1 中去掉 /api)
		newPath := req.URL.Path
		if stripPrefix != "" {
			newPath = strings.TrimPrefix(req.URL.Path, stripPrefix)
		}

		// ⭐ 保留 query
		req.URL.Path = path.Clean(newPath)
		// 3. 补齐开头的斜杠（path.Clean 对于空字符串可能返回 "."）
		if !strings.HasPrefix(req.URL.Path, "/") {
			req.URL.Path = "/" + req.URL.Path
		}

		// ⭐ 可选：设置 Host（更规范）
		req.Host = targetURL.Host

		log.Printf("[Proxy] %s -> %s%s",
			req.URL.Path,
			targetURL,
			newPath,
		)
	}

	// ⭐ 错误处理（非常重要）
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[Proxy Error] %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"bad gateway"}`))
	}

	return proxy, nil
}
