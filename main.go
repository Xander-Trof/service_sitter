package main

import (
	"github.com/gin-gonic/gin"
)



type statusResponse struct {
	Status string `json:"status"`
	ServiceTag string `json:"service_tag"`
	Uptime string `json:"uptime"`
}

type reloadData struct {
	Done bool `json:"done"`
}


func main() {
	router := gin.Default()
	router.GET("/", getHello)
	router.GET("/status/:serviceName", getStatus)
	router.GET("/logs/:serviceName", getLogs)
	router.POST("/reload/:serviceName", reloadService)

	router.Run("localhost:8080")
}

func getHello(c *gin.Context) {

	c.String(200, "Hello World")
}

func getStatus(c *gin.Context) {
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}

func getLogs(c *gin.Context) {
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}

func reloadService(c *gin.Context) {
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}
