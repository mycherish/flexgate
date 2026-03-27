package main

import (
	"flexgate/internal/config"
	"flexgate/internal/gateway"
	"fmt"
	"log"
)

func main() {
	// 1. 加载配置 (确保路径正确)
	cfg, err := config.LoadConfig("configs/gateway.yaml")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 初始化引擎
	engine := gateway.NewEngine(cfg)

	// 3. 启动
	log.Printf("[INFO] FlexGate Gateway is starting...")
	log.Printf("[INFO] Listening and serving HTTP on :%d", cfg.Server.Port)
	if err := engine.Run(fmt.Sprintf(":%v", cfg.Server.Port)); err != nil {
		log.Fatalf("网关运行异常: %v", err)
	}
}
