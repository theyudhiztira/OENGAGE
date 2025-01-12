package whatsapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type CloudApiMessageResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

type WhatsappTemplate struct {
	Name                  string      `json:"name" binding:"required"`
	ParameterFormat       string      `json:"parameter_format"`
	Components            []Component `json:"components" binding:"required"`
	Language              string      `json:"language" binding:"required"`
	Status                string      `json:"status" binding:"required"`
	Category              string      `json:"category" binding:"required"`
	ID                    string      `json:"id,omitempty"`
	MessageSendTTLSeconds int         `json:"message_send_ttl_seconds,omitempty"`
	AllowCategoryChange   string      `json:"allow_category_change,omitempty"`
}

type Component struct {
	Type                      string   `json:"type,omitempty" binding:"oneof=HEADER BODY FOOTER BUTTON CAROUSEL"`
	Format                    string   `json:"format,omitempty"`
	Text                      string   `json:"text,omitempty"`
	Example                   *Example `json:"example,omitempty"`
	Buttons                   []Button `json:"buttons,omitempty"`
	Cards                     []Card   `json:"cards,omitempty"`
	AddSecurityRecommendation bool     `json:"add_security_recommendation,omitempty"`
	CodeExpirationMinutes     int      `json:"code_expiration_minutes,omitempty"`
}

type Example struct {
	HeaderHandle []string   `json:"header_handle,omitempty"`
	BodyText     [][]string `json:"body_text,omitempty"`
}

type Button struct {
	Type        string   `json:"type"`
	Text        string   `json:"text"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	URL         string   `json:"url,omitempty"`
	Example     []string `json:"example,omitempty"`
}

type Card struct {
	Components []Component `json:"components"`
}

type MetaTemplateResponse struct {
	Data   []WhatsappTemplate `json:"data"`
	Paging Paging             `json:"paging"`
}

type Paging struct {
	Cursors Cursor `json:"cursors"`
	Next    string `json:"next,omitempty"`
	Before  string `json:"before,omitempty"`
}

type Cursor struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type WhatsappClientResponse struct {
	CapiSuccess CloudApiMessageResponse
	CapiError   struct {
		Error CloudApiErrorResponse `json:"error"`
	}
	CreateTemplate      CreateTemplateResponse
	GetTemplate         MetaTemplateResponse
	PhoneNumberResponse WhatsappPhoneNumberResponse
}

type WhatsappPhoneNumberResponse struct {
	Data []WhatsappPhoneNumber `json:"data"`
}

type CreateTemplateResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

type CloudApiErrorResponse struct {
	Message     string `json:"message,omitempty"`
	Type        string `json:"type,omitempty"`
	Code        int    `json:"code,omitempty"`
	FBT         string `json:"fbtrace_id,omitempty"`
	SubCode     int    `json:"error_subcode,omitempty"`
	IsTransient bool   `json:"is_transient,omitempty"`
	ErrorData   struct {
		MessagingProduct string `json:"messaging_product,omitempty"`
		Details          string `json:"details,omitempty"`
	} `json:"error_data,omitempty"`
	ErrorUserTitle string `json:"error_user_title,omitempty"`
	ErrorUserMsg   string `json:"error_user_msg,omitempty"`
}

type WhatsappClientParam struct {
	WhatsappBusinesID     string `json:"whatsapp_business_id"`
	WhatsappPhoneNumberID string `json:"whatsapp_phone_number"`
	WhatsappToken         string `json:"whatsapp_token"`
}

func Post(path string, body any, config WhatsappClientParam) (WhatsappClientResponse, error) {
	url := buildURL(config, path)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("[Whatsapp.Post] JSON Marshal failed", err)
		return WhatsappClientResponse{}, err
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(strings.NewReader(string(jsonBody))))
	if err != nil {
		return WhatsappClientResponse{}, err
	}

	setHeaders(req, config.WhatsappToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WhatsappClientResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return WhatsappClientResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return WhatsappClientResponse{}, HandleError(respBody)
	}

	log.Println(string(respBody))

	return WhatsappClientResponse{}, nil
}

func Get(path string, query map[string]string, config WhatsappClientParam) (WhatsappClientResponse, error) {
	url := buildURL(config, path)

	url, err := ParseQueryParam(url, query)
	if err != nil {
		log.Println("[Whatsapp.Get] Failed to parse url", err)
		return WhatsappClientResponse{}, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[Whatsapp.Get] Failed to create new request", err)
		return WhatsappClientResponse{}, err
	}

	setHeaders(req, config.WhatsappToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[Whatsapp.Get] Failed to call CAPI", err)
		return WhatsappClientResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Whatsapp.Get] Failed to read response body", err)
		return WhatsappClientResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return WhatsappClientResponse{}, HandleError(respBody)
	}

	result, err := ParseResponse("GET", path, respBody)
	if err != nil {
		log.Println("[Whatsapp.Get] Failed to parse response", err)
		return WhatsappClientResponse{}, err
	}

	return result, nil
}

func HandleError(resp []byte) error {
	var mBody WhatsappClientResponse
	if err := json.Unmarshal(resp, &mBody.CapiError); err != nil {
		log.Println("[Whatsapp.HandleError] JSON Unmarshal failed ", err)
		return err
	}

	if mBody.CapiError.Error.ErrorData.Details != "" {
		return errors.New(mBody.CapiError.Error.ErrorData.Details)
	} else if mBody.CapiError.Error.ErrorUserMsg != "" {
		return errors.New(mBody.CapiError.Error.ErrorUserMsg)
	}

	return errors.New(mBody.CapiError.Error.Message)
}

func ParseQueryParam(urlString string, params map[string]string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	query := u.Query()

	for key, value := range params {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

func ParseResponse(method string, path string, resp []byte) (WhatsappClientResponse, error) {
	var result WhatsappClientResponse
	switch path {
	case "message_templates":
		if method == "GET" {
			if err := json.Unmarshal(resp, &result.GetTemplate); err != nil {
				log.Println("[Whatsapp.ParseResponse] JSON Unmarshal failed", err)
				return WhatsappClientResponse{}, err
			}
		} else if method == "POST" {
			if err := json.Unmarshal(resp, &result.CreateTemplate); err != nil {
				log.Println("[Whatsapp.ParseResponse] JSON Unmarshal failed", err)
				return WhatsappClientResponse{}, err
			}
		}
	case "messages":
		if err := json.Unmarshal(resp, &result.CapiSuccess); err != nil {
			log.Println("[Whatsapp.ParseResponse] JSON Unmarshal failed", err)
			return WhatsappClientResponse{}, err
		}
	case "phone_numbers":
		if err := json.Unmarshal(resp, &result.PhoneNumberResponse); err != nil {
			log.Println("[Whatsapp.ParseResponse] JSON Unmarshal failed", err)
			return WhatsappClientResponse{}, err
		}
	}

	return result, nil
}

func buildURL(config WhatsappClientParam, path string) string {
	if config.WhatsappPhoneNumberID != "" {
		return fmt.Sprintf("https://graph.facebook.com/v21.0/%s/%s", config.WhatsappPhoneNumberID, path)
	}
	return fmt.Sprintf("https://graph.facebook.com/v21.0/%s/%s", config.WhatsappBusinesID, path)
}

func setHeaders(req *http.Request, token string) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
}
