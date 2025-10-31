package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config := GetConfig()

	minioClient, err := InitMinIO(&config)
	if err != nil {
		log.Fatalf("Falha ao inicializar o cliente MinIO: %v", err)
	}

	err = EnsureBucketExists(minioClient, config.MinioBucket)
	if err != nil {
		log.Fatalf("Falha ao garantir a existência do bucket: %v", err)
	}

	h := &Handler{
		minioClient: minioClient,
		config:      &config,
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	loginRouter := router.Group("admin")
	{
		loginRouter.POST("/login", h.UserLoginCodeProcess)
		loginRouter.POST("/validate", h.ValidateLoginCode)
	}

	UpDelAuthRouter := router.Group("/")
	UpDelAuthRouter.Use(AuthMiddleware(config.UploadSecret))
	{
		UpDelAuthRouter.POST("/upload", h.UploadFile)
		UpDelAuthRouter.DELETE("/files/:filename", h.DeleteFile)
	}

	downAuthRouter := router.Group("/")
	{
		downAuthRouter.GET("/files/:filename", h.GetFile)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "6868"
	}

	log.Printf("Servidor de armazenamento iniciado na porta :%s", port)
	log.Printf("Bucket para uploads: '%s'", config.MinioBucket)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Falha ao iniciar o servidor Gin: %v", err)
	}
}
