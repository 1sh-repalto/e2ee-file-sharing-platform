package handler

import (
	"net/http"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShareHandler struct {
	shareUsecase usecase.ShareUsecase
}

func NewShareHandler(shareUsecase usecase.ShareUsecase) *ShareHandler {
	return &ShareHandler{shareUsecase: shareUsecase}
}

func (h *ShareHandler) ShareFile(c *gin.Context) {
	var req struct {
		FileID			string	`json:"file_id" binding:"required"`
		RecipientID		string	`json:"recipient_id" binding:"required"`
		WrappedKey		[]byte	`json:"wrapped_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileID, _ := uuid.Parse(req.FileID)
	recipientID, _ := uuid.Parse(req.RecipientID)
	uid, _ := c.Get("userID")
	ownerID := uid.(uuid.UUID)

	if err := h.shareUsecase.ShareFile(c.Request.Context(), fileID, ownerID, recipientID, req.WrappedKey); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file shared"})
}

func (h *ShareHandler) ListShares(c *gin.Context) {
	uid, _ := c.Get("userID")
	recipientID := uid.(uuid.UUID)

	shares, err := h.shareUsecase.GetSharesForRecipient(c.Request.Context(), recipientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shares)
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	uid, _ := c.Get("userID")
	recipientID := uid.(uuid.UUID)

	share, err := h.shareUsecase.GetShare(c.Request.Context(), fileID, recipientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
		return
	}
	c.JSON(http.StatusOK, share)
}

func (h *ShareHandler) Unshare(c * gin.Context) {
	shareID, err := uuid.Parse(c.Param("share_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share id"})
		return
	}

	uid, _ := c.Get("userID")
	ownerID := uid.(uuid.UUID)

	if err := h.shareUsecase.Unshare(c.Request.Context(), shareID, ownerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "share revoked"})
}