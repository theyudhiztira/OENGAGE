package template

import (
	"context"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
)

func NewTemplateHandler(ctx *context.Context, s templateService) *templateHandler {
	return &templateHandler{
		Ctx:     ctx,
		Service: s,
	}
}

func (h *templateHandler) GetTemplate(c *gin.Context) {
	var q TemplateQueryParam
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
			Data:    err.Error(),
		})
	}

	res, err := h.Service.GetTemplate(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ApiResponse{
			Message: pkg.InternalServerError,
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
