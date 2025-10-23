package ports

import (
	"github.com/gin-gonic/gin"
)

// IUserHandler
type IUserHandler interface {
	ListUsers() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
}