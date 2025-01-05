package template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewTemplateHandler(s TemplateService) *templateHandler {
	return &templateHandler{
		Service: s,
	}
}

func (s *templateHandler) GetTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
