package gateway

import (
	"flexgate/internal/config"
	"flexgate/internal/proxy"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	ginEngine *gin.Engine
	// Key: PathPrefix (匹配前缀), Value: 预初始化的代理实例
	proxyPool map[string]*httputil.ReverseProxy
}

func NewEngine(cfg *config.Config) *Engine {
	// 生产环境下建议使用 ReleaseMode 减少日志开销
	gin.SetMode(gin.ReleaseMode)

	e := &Engine{
		ginEngine: gin.New(), // 使用 New 配合自定义中间件更专业
		proxyPool: make(map[string]*httputil.ReverseProxy),
	}

	// 注册基础中间件
	e.ginEngine.Use(gin.Logger())
	e.ginEngine.Use(gin.Recovery())

	// 核心：启动时完成池化，避免请求时重复创建对象
	for _, route := range cfg.Routes {
		p, err := proxy.NewReverseProxy(route.Upstream, route.StripPrefix)
		if err != nil {
			log.Printf("[Engine] 路由 %s 代理初始化失败: %v", route.PathPrefix, err)
			continue
		}
		e.proxyPool[route.PathPrefix] = p
		log.Printf("[Engine] 路由加载成功: %s -> %s (Strip: %s)",
			route.PathPrefix, route.Upstream, route.StripPrefix)
	}

	e.setupRoutes()
	return e
}

func (e *Engine) setupRoutes() {
	// 匹配所有路径
	e.ginEngine.Any("/*path", func(c *gin.Context) {
		fullPath := c.Request.URL.Path

		// 1. 寻找匹配的路由（当前为线性匹配，简历上可写为后期优化点）
		var targetProxy *httputil.ReverseProxy
		for prefix, p := range e.proxyPool {
			if strings.HasPrefix(fullPath, prefix) {
				targetProxy = p
				break
			}
		}

		// 2. 转发或报错
		if targetProxy != nil {
			// 直接使用缓存的 Proxy，复用底层的连接池
			targetProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "FlexGate: Route not found",
			"path":    fullPath,
			"message": "请检查 gateway.yaml 中的 path_prefix 配置",
		})
	})
}

func (e *Engine) Run(addr string) error {
	return e.ginEngine.Run(addr)
}
