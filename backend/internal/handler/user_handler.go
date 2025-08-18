package handler

import (
	"net/http"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uc}
}

type registerRequest struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	PublicKey           []byte `json:"publicKey"`
	EncryptedPrivateKey []byte `json:"encryptedPrivateKey"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUsecase.Register(c.Request.Context(), req.Username, req.Password, req.PublicKey, req.EncryptedPrivateKey)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"createdAt": user.CreatedAt,
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, encryptedPrivateKey, err := h.userUsecase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                    user.ID,
		"username":              user.Username,
		"encrypted_private_key": encryptedPrivateKey,
	})
}
