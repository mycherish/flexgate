package gateway

import (
	"flexgate/internal/config"
	"flexgate/internal/limiter"
	"flexgate/internal/middleware"
	"flexgate/internal/proxy"
	"log"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	ginEngine *gin.Engine
	// Key: PathPrefix (匹配前缀), Value: 预初始化的代理实例
	proxyPool map[string]*httputil.ReverseProxy
	config    *config.Config
}

func NewEngine(cfg *config.Config) *Engine {
	// 生产环境下建议使用 ReleaseMode 减少日志开销
	gin.SetMode(gin.ReleaseMode)
	// 创建干净的 Gin 实例
	r := gin.New()
	r.RedirectTrailingSlash = false
	// 3. 注册核心中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	// 优先使用我们手写的 Logger，面试时好聊逻辑
	r.Use(middleware.Logger())

	e := &Engine{
		ginEngine: r,
		proxyPool: make(map[string]*httputil.ReverseProxy),
		config:    cfg,
	}

	e.setupRoutes()
	return e
}

// setupRoutes 处理全量请求转发
func (e *Engine) setupRoutes() {
	for _, route := range e.config.Routes {
		// --- 1. 初始化反向代理实例 ---
		p, err := proxy.NewReverseProxy(route.Upstream, route.StripPrefix)
		if err != nil {
			log.Printf("[Engine] ⚠️ 路由 [%s] 初始化失败: %v", route.PathPrefix, err)
			continue
		}
		e.proxyPool[route.PathPrefix] = p

		// --- 2. 初始化该路由专属的限流器 ---
		rate := route.RateLimit.Rate
		cap := route.RateLimit.Capacity
		// 如果没配或者配错了，给个宽泛的默认值
		if rate <= 0 {
			rate, cap = 1000, 1000
		}
		tb := limiter.NewTokenBucket(rate, cap)

		// --- 3. 核心：创建路由组并绑定专属中间件 ---
		// 这样每个 prefix 都有自己独立的“安检闸机”
		group := e.ginEngine.Group(route.PathPrefix)
		{
			// 挂载限流中间件（只对当前 group 生效）
			group.Use(middleware.RateLimitMiddleware(tb))

			// 实际转发逻辑
			// 1. 匹配 /api/users 这种不带斜杠的根路径
			group.Any("", func(c *gin.Context) {
				p.ServeHTTP(c.Writer, c.Request)
			})

			// 2. 匹配 /api/users/ 以及后续所有子路径
			group.Any("/*any", func(c *gin.Context) {
				p.ServeHTTP(c.Writer, c.Request)
			})
		}

		log.Printf("[Engine] ✅ 路由映射成功: %s -> %s (剥离前缀: %s)[QPS: %d]",
			route.PathPrefix, route.Upstream, route.StripPrefix, rate)
	}
}

// Run 启动 HTTP 服务
func (e *Engine) Run(addr string) error {
	// addr 格式应为 ":8080"
	return e.ginEngine.Run(addr)
}
