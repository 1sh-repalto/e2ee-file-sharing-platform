package handler

import (
	"net/http"

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
	var req struct {
		Filename     string `json:"filename" binding:"required"`
		MimeType     string `json:"mime_type" binding:"required"`
		Size         int64  `json:"size" binding:"required"`
		IV           []byte `json:"iv" binding:"required"`
		EncryptedKey []byte `json:"encrypted_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ownerID, ok := uid.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	file := domain.File{
		OwnerID:      ownerID,
		Filename:     req.Filename,
		MimeType:     req.MimeType,
		Size:         req.Size,
		IV:           req.IV,
		EncryptedKey: req.EncryptedKey,
	}

	if err := h.fileUsecase.Upload(c.Request.Context(), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file uploaded", "id": file.ID})

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
