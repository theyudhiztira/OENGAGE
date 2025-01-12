package template

import (
	"context"
	"encoding/json"
	"net/http"
	"theyudhiztira/oengage-backend/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewTemplateHandler(ctx *context.Context, s templateService) *templateHandler {
	return &templateHandler{
		Ctx:     ctx,
		Service: s,
	}
}

func (h *templateHandler) GetTemplate(c *gin.Context) {
	var q TemplateQueryParam
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
			Data:    pkg.ErrorResp(err),
		})
		return
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

func (h *templateHandler) CreateTemplate(c *gin.Context) {
	var q CreateTemplateRequest
	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, pkg.ApiResponse{
			Message: pkg.BadRequest,
			Status:  false,
			Data:    pkg.ErrorResp(err),
		})
		return
	}

	p := WhatsappTemplate{
		Name:                q.Name,
		AllowCategoryChange: q.AllowCategoryChange,
		Category:            q.Category,
		Language:            q.Language,
	}

	bodyComponents := []string{q.Body, q.Header, q.Footer, q.Button, q.Carousel}
	for _, componentString := range bodyComponents {
		if componentString != "" {
			component, err := h.ParseComponent(c, componentString)
			if err != nil {
				c.JSON(http.StatusBadRequest, pkg.ApiResponse{
					Message: pkg.BadRequest,
					Status:  false,
					Data:    pkg.ErrorResp(err),
				})
				return
			}

			if component.Type != "" {
				p.Components = append(p.Components, component)
			}
		}
	}

	if _, err := h.Service.CreateTemplate(p); err != nil {
		c.JSON(http.StatusInternalServerError, pkg.ApiResponse{
			Message: err.Error(),
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusOK, pkg.ApiResponse{
		Message: "Nice",
		Status:  true,
	})
}

func (h *templateHandler) ParseComponent(c *gin.Context, bString string) (Component, error) {
	var b Component
	if err := json.Unmarshal([]byte(bString), &b); err != nil {
		return Component{}, err
	}

	if err := validate.Struct(b); err != nil {
		return Component{}, err
	}

	return b, nil
}
