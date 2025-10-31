package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sanisamoj/image_repo/cache"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	minioClient *minio.Client
	config      *Config
}

type UploadResponse struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

func (u *Handler) UserLoginCodeProcess(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configEmail := GetConfig().Email
	if req.Email != configEmail {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	code := gen6DigCod()
	key := fmt.Sprintf("login-code-%s", req.Email)
	if err := cache.Set(key, code, 5*time.Minute); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	go SendLoginCodeEmail(req.Email, "superadmin", req.Email, code, 5)

	c.JSON(200, gin.H{"message": "code sent"})
}

func (u *Handler) ValidateLoginCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fmt.Sprintf("login-code-%s", req.Email)
	val, err := cache.Get(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}

	if val == req.Code {
		secretKey := GetConfig().UploadSecret
		jwt := NewJWTGenerator(secretKey)
		tokenStr, err := jwt.GenerateLoginToken(req.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(200, gin.H{"token": tokenStr})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
}

func (h *Handler) UploadFile(c *gin.Context) {
	bucketName := "uploads"
	err := EnsureBucketExists(h.minioClient, bucketName)
	if err != nil {
		log.Printf("Erro ao garantir a existência do bucket '%s': %v", bucketName, err)
		detailedError := fmt.Sprintf("Não foi possível processar o bucket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": detailedError})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O campo 'file' é obrigatório"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	objectName := uuid.NewString() + ext
	contentType := header.Header.Get("Content-Type")
	fileSize := header.Size

	_, err = h.minioClient.PutObject(c.Request.Context(), bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		log.Printf("Erro ao fazer upload para o MinIO: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível salvar o arquivo"})
		return
	}

	host := c.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}

	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}

	//fileURL := fmt.Sprintf("%s://%s/files/%s", scheme, c.Request.Host, objectName)
	var selfHost *string
	if h.config.SelfHost != nil {
		selfHost = h.config.SelfHost
	} else {
		selfHost = &host
	}

	fileURL := fmt.Sprintf("%s://%s/files/%s", scheme, *selfHost, objectName)

	response := UploadResponse{
		FileName: objectName,
		URL:      fileURL,
		Size:     fileSize,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetFile(c *gin.Context) {
	bucketName := "uploads"
	objectName := c.Param("filename")

	object, err := h.minioClient.GetObject(c.Request.Context(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
		return
	}

	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
	c.Header("Content-Disposition", "inline; filename=\""+objectName+"\"")
	c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, nil)
}

func (h *Handler) DeleteFile(c *gin.Context) {
	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orgIDStr, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (orgID type)"})
		return
	}

	bucketName := fmt.Sprintf("orgid-%s", orgIDStr)
	objectName := c.Param("filename")

	err := h.minioClient.RemoveObject(c.Request.Context(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível remover o arquivo"})
		return
	}

	c.Status(http.StatusNoContent)
}
