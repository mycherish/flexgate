package main

import (
	"fmt"
	"log"

	"flexgate/internal/config"
	"flexgate/internal/proxy"
	"flexgate/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/gateway.yaml")
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	// 2. 初始化 router
	r := router.NewRouter(cfg.Routes)

	// 3. 初始化 Gin
	engine := gin.Default()

	// 4. 核心 handler
	engine.Any("/*path", func(c *gin.Context) {
		path := c.Request.URL.Path

		route := r.Match(path)
		if route == nil {
			c.JSON(404, gin.H{"error": "route not found"})
			return
		}

		// 创建代理
		proxy, err := proxy.NewReverseProxy(route.Upstream, route.PathPrefix, route.StripPrefix)
		if err != nil {
			c.JSON(500, gin.H{"error": "proxy error"})
			return
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// 5. 启动
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Println("FlexGate running on", addr)
	engine.Run(addr)
}
