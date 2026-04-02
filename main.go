package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Xander-Trof/service-sitter/dockercomands"
	"github.com/docker/docker/api/types"
)



type statusResponse struct {
	Status string `json:"status"`
	ServiceTag string `json:"service_tag"`
	Uptime string `json:"uptime"`
}

type reloadData struct {
	Done bool `json:"done"`
}

type ContainersResponse struct {
	ActiveContainers []string `json:"active_containers"`
}

type ContainerStatusResponse struct {
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
	Image       string `json:"image"`
	CreatedAt   int64  `json:"created_at"`
	Health     *types.Health `json:",omitempty"`
}

type GeneralStatusResponse struct {
	Services []ContainerStatusResponse `json:"services"`
}

func main() {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	router := gin.Default()

	// Middleware для проверки API-ключа
	router.Use(apiKeyMiddleware())

	router.GET("/", getDescription)
	router.GET("/status/:serviceName", getServiceStatus)
	router.GET("/status", getGeneralStatus)
	router.GET("/logs/:serviceName", getLogs)
	router.POST("/reload/:serviceName", reloadService)
	router.POST("/reload", reloadAllServices)

	router.Run("localhost:8080")
}

// Middleware для проверки API-ключа в заголовке
func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ключ из .env
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			log.Println("WARNING: API_KEY not set in .env")
			c.JSON(500, gin.H{"error": "server configuration error"})
			c.Abort()
			return
		}

		// Получаем ключ из заголовка запроса
		requestKey := c.GetHeader("AuthorizationKey")

		// Проверяем ключ
		if requestKey == "" {
			c.JSON(401, gin.H{"error": "missing API key in header"})
			c.Abort()
			return
		}

		if requestKey != apiKey {
			log.Printf("Unauthorized request with key: %s", requestKey)
			c.JSON(401, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		// Ключ валиден — продолжаем обработку
		c.Next()
	}
}

func getDescription(c *gin.Context) {
	apiInfo := []gin.H{
			{
				"endpoint":    "/",
				"method":      "GET",
				"description": "Возвращает список с API документацией",
			},
			{
				"endpoint":    "/status",
				"method":      "GET",
				"description": "Возвращает статусы всех контейнеров (запущенных и упавших)",
			},
			{
				"endpoint":    "/status/:serviceName",
				"method":      "GET",
				"description": "Возвращает статус указанного контейнера",
			},
			{
				"endpoint":    "/logs/:serviceName",
				"method":      "GET",
				"description": "Возвращает логи указанного контейнера",
			},
			{
				"endpoint":    "/reload/:serviceName",
				"method":      "POST",
				"description": "Перезапускает указанный контейнер",
			},
			{
				"endpoint":    "/reload",
				"method":      "POST",
				"description": "Перезапускает все контейнеры",
			},
		}
	c.JSON(200, apiInfo)
}

func getGeneralStatus(c *gin.Context) {
	// Получение общей информации о сервисах
	containers := dockercomands.DockerPS()

	services := make([]ContainerStatusResponse, 0, len(containers))
	for _, c := range containers {
		services = append(services, ContainerStatusResponse{
			ServiceName: c.Names[0][1:],
			Status:      c.Status,
			Image:       c.Image,
			CreatedAt:   c.Created,
		})
	}
	
	respData := GeneralStatusResponse{Services: services}
	c.JSON(200, respData)
}

func getServiceStatus(c *gin.Context) {
	// Получение статуса по имени сервиса
	serviceName := c.Param("serviceName")
	allContainers := dockercomands.DockerPS()
	container := dockercomands.FindContainerByName(allContainers, serviceName)
	if container == nil {
		s := fmt.Sprintf("container not found %s", serviceName)
		c.JSON(404, gin.H{"error": s})
		return
	}
	respData := ContainerStatusResponse{ServiceName: serviceName, Status: container.Status, Image: container.Image, CreatedAt: container.Created}

	c.JSON(200, respData)
}

func getLogs(c *gin.Context) {
	// Получение логов по имени сервиса
	serviceName := c.Param("serviceName")
	logsStream := dockercomands.DockerLogs(serviceName)
	if logsStream == nil {
		s := fmt.Sprintf("container %s not found or logs not available", serviceName)
		c.JSON(404, gin.H{"error": s})
		return
	}
	defer logsStream.Close()

	// Читаем логи в память для логирования
	logsBytes, err := io.ReadAll(logsStream)
	if err != nil {
		log.Printf("[ERROR reading logs for %s]: %v", serviceName, err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to read logs: %v", err)})
		return
	}

	log.Printf("[LOGS for %s] length=%d bytes\n%s", serviceName, len(logsBytes), string(logsBytes))

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+serviceName+".log")

	c.String(200, string(logsBytes))
}

func reloadService(c *gin.Context) {
	// Перезапуск сервиса по его имени
	serviceName := c.Param("serviceName")

	result, err := dockercomands.DockerRestart(serviceName)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to restart container: %v", err)})
		return
	}

	c.JSON(200, result)
}

func reloadAllServices(c *gin.Context) {
	// Перезапуск всех сервисов
	result := dockercomands.DockerRestartAll()

	c.JSON(200, result)
}