package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hui-cyber/BoolCore/backend/internal/api"
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = corsOrigins()
	router.Use(cors.New(config))
	// 注册我们定义的所有路由
	api.RegisterRoutes(router)

	addr := serverAddress()
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func corsOrigins() []string {
	origins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	if origins == "" {
		return []string{"http://localhost:5173", "http://localhost:4173"}
	}
	return splitAndTrim(origins)
}

func serverAddress() string {
	addr := strings.TrimSpace(os.Getenv("ADDR"))
	if addr != "" {
		return addr
	}
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func splitAndTrim(input string) []string {
	parts := strings.Split(input, ",")
	trimmed := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			trimmed = append(trimmed, value)
		}
	}
	return trimmed
}
