package template

import (
	"fmt"
	"log"
	"reflect"
	"theyudhiztira/oengage-backend/internal/whatsapp"
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

	query := buildQueryParams(q)

	res, err := whatsapp.Get("message_templates", query, whatsapp.WhatsappConfig{
		WhatsappBusinesID: sc.WhatsappWabaID,
		WhatsappToken:     sc.WhatsappToken,
	})
	if err != nil {
		log.Println("[TemplateService.GetTemplate] Failed to get template", err)
		return nil, err
	}

	return res.GetTemplate, nil
}

func (s *templateService) CreateTemplate(payload WhatsappTemplate) (MetaCreateTemplateResponse, error) {
	sc, err := s.Repository.GetWhatsappConfig()
	if err != nil {
		return MetaCreateTemplateResponse{}, err
	}

	res, err := whatsapp.Post("message_templates", payload, whatsapp.WhatsappConfig{
		WhatsappBusinesID: sc.WhatsappWabaID,
		WhatsappToken:     sc.WhatsappToken,
	})
	if err != nil {
		log.Println("[TemplateService.CreateTemplate] Failed to create template", err)
		return MetaCreateTemplateResponse{}, err
	}

	return MetaCreateTemplateResponse{
		ID:       res.CreateTemplate.ID,
		Status:   res.CreateTemplate.Status,
		Category: res.CreateTemplate.Category,
	}, nil
}

func buildQueryParams(q TemplateQueryParam) map[string]string {
	query := make(map[string]string)
	v := reflect.ValueOf(q)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		key := field.Tag.Get("form")
		if key == "" {
			key = field.Name
		}
		value := v.Field(i).Interface()
		if value != "" {
			query[key] = fmt.Sprintf("%v", value)
		}
	}
	return query
}
