package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "user-service",
			"data":    []string{"徐卫东", "Bob"},
		})
	})
	// GET /users/:id - 获取单个用户
	r.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}
		// 模拟用户数据，实际项目中通常从数据库查询
		users := map[int]string{
			1: "徐卫东",
			2: "Bob",
		}
		if name, ok := users[id]; ok {
			c.JSON(http.StatusOK, gin.H{
				"id":   id,
				"name": name,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		}
	})

	r.Run(":9001")
}
