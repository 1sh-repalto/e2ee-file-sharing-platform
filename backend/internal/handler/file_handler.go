package handler

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/domain"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	fileUsecase usecase.FileUsecase
}

func NewFileHandler(fu usecase.FileUsecase) *FileHandler {
	return &FileHandler{fileUsecase: fu}
}

func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	fileContent, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer fileContent.Close()

	iv := c.PostForm("iv")
	encryptedKey := c.PostForm("encrypted_key")
	mimeType := c.PostForm("mime_type")
	size := fileHeader.Size
	filename := fileHeader.Filename

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := uid.(uuid.UUID)

	file := domain.File{
		OwnerID:      ownerID,
		Filename:     filename,
		MimeType:     mimeType,
		Size:         size,
		IV:           []byte(iv),
		EncryptedKey: []byte(encryptedKey),
	}

	if err := h.fileUsecase.Upload(c.Request.Context(), file, fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file uploaded", "id": file.ID})
}

func (h *FileHandler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	recipientID := uid.(uuid.UUID)

	content, file, wrappedKey, err := h.fileUsecase.Download(c.Request.Context(), id, recipientID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	defer content.Close()

	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	if len(wrappedKey) > 0 {
		c.Header("X-Wrapped-Key", base64.StdEncoding.EncodeToString(wrappedKey))
	}

	if len(file.IV) > 0 {
		c.Header("X-IV", base64.StdEncoding.EncodeToString(file.IV))
	}

	_, err = io.Copy(c.Writer, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stream file"})
		return
	}
}

func (h *FileHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	file, err := h.fileUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	uid, _ := c.Get("userID")
	if file.OwnerID != uid.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) ListByOwner(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ownerID := uid.(uuid.UUID)

	files, err := h.fileUsecase.ListByOwner(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := uid.(uuid.UUID)

	if err := h.fileUsecase.Delete(c.Request.Context(), id, ownerID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}