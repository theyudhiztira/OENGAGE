package template

import (
	"fmt"
	"net/url"
	"reflect"
)

func NewTemplateService(r *templateRepository) *templateService {
	return &templateService{
		Repository: *r,
	}
}

func (s *templateService) GetTemplate(q TemplateQueryParam) (interface{}, error) {
	sc, err := s.Repository.GetWhatsappConfig()
	if err != nil {
		return nil, err
	}

	urlStr, parseErr := parseTemplateQueryURL(q, sc.WhatsappWabaID)
	if parseErr != nil {
		return nil, parseErr
	}

	res, err := s.Repository.GetWhatsappTemplate(urlStr, sc.WhatsappToken)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseTemplateQueryURL(params TemplateQueryParam, wabaID string) (string, error) {
	baseURL := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/message_templates", wabaID)
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	query := u.Query()

	v := reflect.ValueOf(params)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		key := field.Tag.Get("form")
		if key == "" {
			key = field.Name
		}
		value := v.Field(i).Interface()
		if value != "" {
			query.Set(key, fmt.Sprintf("%v", value))
		}
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}
