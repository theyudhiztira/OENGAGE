package auth

import (
	"context"
	"log"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
)

func NewAuthHandler(ctx *context.Context, service authService) *authHandler {
	return &authHandler{
		Ctx:     ctx,
		Service: service,
	}
}

func (r *authHandler) SetupSystem(c *gin.Context) {
	var rData SystemSetupRequest
	if err := c.ShouldBindJSON(&rData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
		})
		return
	}

	result, err := r.Service.SetupSystem(rData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, pkg.ApiResponse{
			Message: pkg.InternalServerError,
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusOK, pkg.ApiResponse{
		Message: pkg.SystemSetupSuccess,
		Status:  true,
		Data:    result,
	})
}

func (r *authHandler) LoginHandler(c *gin.Context) {
	var rData LoginRequest
	if err := c.ShouldBindJSON(&rData); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: LoginBadRequest,
			Status:  false,
		})
		return
	}

	result, err := r.Service.Login(rData.Username, rData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, pkg.ApiResponse{
			Message: pkg.Unauthorized,
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
