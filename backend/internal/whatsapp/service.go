package whatsapp

import (
	"log"
)

type whatsappService struct {
	Repository whatsappRepository
}

type WhatsappConfigRequest struct {
	WhatsappToken  string `json:"whatsapp_token" binding:"required"`
	WhatsappWabaID string `json:"whatsapp_waba_id" binding:"required"`
}

type WhatsappConfigResponse struct {
	WhatsappWabaID      string                `json:"whatsapp_waba_id"`
	WhatsappPhoneNumber []WhatsappPhoneNumber `json:"whatsapp_phone_number"`
}

type WhatsappPhoneNumber struct {
	VerifiedName       string `bson:"verified_name" json:"verified_name"`
	DisplayPhoneNumber string `bson:"display_phone_number" json:"display_phone_number"`
	ID                 string `bson:"id" json:"id"`
	QualityRating      string `bson:"quality_rating" json:"quality_rating"`
}

func NewWhatsappService(r *whatsappRepository) *whatsappService {
	return &whatsappService{
		Repository: *r,
	}
}

func (s *whatsappService) ConfigWhatsapp(wcr WhatsappConfigRequest) (WhatsappConfigResponse, error) {
	client, err := Get("phone_numbers", nil, WhatsappClientParam{
		WhatsappBusinesID: wcr.WhatsappWabaID,
		WhatsappToken:     wcr.WhatsappToken,
	})
	if err != nil {
		log.Println("[WhatsappService.ConfigWhatsapp] Failed to get phone numbers", err)
		return WhatsappConfigResponse{}, err
	}

	phoneNumbers := []WhatsappPhoneNumber{}
	for _, pn := range client.PhoneNumberResponse.Data {
		phoneNumbers = append(phoneNumbers, WhatsappPhoneNumber{
			VerifiedName:       pn.VerifiedName,
			DisplayPhoneNumber: pn.DisplayPhoneNumber,
			ID:                 pn.ID,
			QualityRating:      pn.QualityRating,
		})
	}

	update, err := s.Repository.FindAndUpdate(WhatsappConfig{
		WhatsappBusinesID:   wcr.WhatsappWabaID,
		WhatsappToken:       wcr.WhatsappToken,
		WhatsappPhoneNumber: phoneNumbers,
	})
	if err != nil {
		log.Println("[WhatsappService.ConfigWhatsapp] Failed to update whatsapp config", err)
		return WhatsappConfigResponse{}, err
	}

	return WhatsappConfigResponse{
		WhatsappWabaID:      update.WhatsappBusinesID,
		WhatsappPhoneNumber: update.WhatsappPhoneNumber,
	}, nil
}
