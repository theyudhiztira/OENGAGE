package contacts

import (
	"context"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
)

type contactsHandler struct {
	Ctx     *context.Context
	Service contactsService
}

func NewContactsHandler(ctx *context.Context, service contactsService) *contactsHandler {
	return &contactsHandler{
		Ctx:     ctx,
		Service: service,
	}
}

func (h *contactsHandler) Create(c *gin.Context) {
	body := []CreateContactRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
			Error:   pkg.ParseValidationMessage(err),
		})
		return
	}

	res, err := h.Service.CreateContacts(body, c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ApiResponse{
			Message: pkg.InternalServerError,
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusCreated, pkg.ApiResponse{
		Data:   res,
		Status: false,
	})
	return
}
