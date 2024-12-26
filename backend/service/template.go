package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"theyudhiztira/oengage-backend/config"
	"theyudhiztira/oengage-backend/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateService interface {
	GetTemplate(c *gin.Context)
	FetchFacebookTemplate(c *gin.Context)
}

type templateService struct {
	DB  mongo.Database
	Ctx *context.Context
}

func NewTemplateService(db mongo.Database, ctx *context.Context) *templateService {
	return &templateService{DB: db, Ctx: ctx}
}

func (r *templateService) GetTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (r *templateService) FetchFacebookTemplate(c *gin.Context) {
	url, err := r.parseTemplateQueryUrl(c.Request.URL.Query())
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Env("META_TOKEN")))

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("Failed with status code %d", resp.StatusCode)})
		return
	}

	var jsonResponse dto.MetaTemplateResponse
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response"})
		return
	}

	response := dto.TemplateResponseDto{
		Data:       jsonResponse.Data,
		Pagination: jsonResponse.Paging.Cursors,
	}

	c.JSON(http.StatusOK, response)
}

func (r *templateService) parseTemplateQueryUrl(params map[string][]string) (string, string) {
	baseUrl := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/message_templates", config.Env("WABA_ID"))

	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err.Error()
	}

	query := u.Query()
	for key, value := range params {
		if key == "page" {
			query.Add("after", value[0])
		} else {
			query.Add(key, value[0])
		}
	}
	u.RawQuery = query.Encode()

	return u.String(), ""
}
