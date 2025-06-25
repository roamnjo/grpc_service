package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/roamnjo/grpc_service/pkg/hash"
	"github.com/roamnjo/grpc_service/pkg/token"
)

type Handler struct {
	repo Repository
	log  *slog.Logger
}

type InRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpRequest struct {
	Name string `json:"name"`
	InRequest
}

func NewHandler(repo Repository, log *slog.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) SignUp(c *gin.Context) {
	var upreq UpRequest

	err := c.BindJSON(&upreq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid request"})
		return
	}

	err = h.repo.FindEmail(context.Background(), upreq.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	err = h.repo.FindSameName(context.Background(), upreq.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name already exists"})
		return
	}

	hashedPassword, err := hash.HashPassword(upreq.Password)
	if err != nil {
		h.log.Error("Error: hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
		return
	}

	upreq.Password = hashedPassword

	err = h.repo.CreateUser(context.Background(), &User{
		Name:     upreq.Name,
		Email:    upreq.Email,
		Password: upreq.Password,
	})

	if err != nil {
		h.log.Error("error: create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unabale create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"User created": upreq})
}

func (h *Handler) SignIn(c *gin.Context) {
	var inreq InRequest
	var upreq UpRequest

	err := c.BindJSON(&inreq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid request"})
		return
	}

	err = h.repo.FindEmail(context.Background(), inreq.Email)
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email doesn't exist"})
		return
	}

	if !hash.CheckPasswordHash(inreq.Password, upreq.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	newToken := token.GenerateNewToken()
	h.log.Info("New token is", newToken)

	c.JSON(http.StatusAccepted, gin.H{"Status": "Accepted"})
}
