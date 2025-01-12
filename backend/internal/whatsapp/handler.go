package whatsapp

import (
	"context"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
)

type whatsappHandler struct {
	Ctx     *context.Context
	Service whatsappService
}

func NewWhatsappHandler(ctx *context.Context, service whatsappService) *whatsappHandler {
	return &whatsappHandler{
		Ctx:     ctx,
		Service: service,
	}
}

func (h *whatsappHandler) ConfigHandler(c *gin.Context) {
	body := WhatsappConfigRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
		})
		return
	}

	res, err := h.Service.ConfigWhatsapp(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ApiResponse{
			Message: pkg.InternalServerError,
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusOK, pkg.ApiResponse{
		Status: true,
		Data:   res,
	})
}
