package broadcast

import (
	"context"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
)

type broadcastHandler struct {
	Ctx     *context.Context
	Service broadcastService
}

func NewBroadcastHandler(ctx *context.Context, service broadcastService) *broadcastHandler {
	return &broadcastHandler{
		Ctx:     ctx,
		Service: service,
	}
}

func (h *broadcastHandler) CreateBroadcast(c *gin.Context) {
	body := CreateBroadcastRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
			Error:   pkg.ParseValidationMessage(err),
		})
		return
	}

	switch body.BroadcastChannel {
	case "whatsapp":
		res, err := h.Service.CreateWhatsappBroadcast(body)
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
		return
	default:
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: "Invalid broadcast channel",
			Status:  false,
		})
		return
	}
}
