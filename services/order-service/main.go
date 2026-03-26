package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/orders", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "order-service",
			"data":    []int{101, 102},
		})
	})

	r.Run(":9002")
}
