package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Xander-Trof/service-sitter/dockercomands"
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
	router.GET("/", getDescription)
	router.GET("/status/:serviceName", getStatus)
	router.GET("/logs/:serviceName", getLogs)
	router.POST("/reload/:serviceName", reloadService)

	router.Run("localhost:8080")
}

func getDescription(c *gin.Context) {
	// Описание ручек сервиса
	// c.String(200, "Hello World")
	containers := dockercomands.DockerPS()
	c.JSON(200, containers)
}

func getStatus(c *gin.Context) {
	// Получение статуса по имени сервиса
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}

func getLogs(c *gin.Context) {
	// Получение логов по имени сервиса
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}

func reloadService(c *gin.Context) {
	// Перезапуск сервиса по его имени
	serviceName := c.Param("serviceName")

	c.String(200, serviceName)
}
