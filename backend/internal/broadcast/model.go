package broadcast

type CreateBroadcastRequest struct {
	BroadcastName    string                   `json:"broadcast_name" binding:"required"`
	BroadcastChannel string                   `json:"broadcast_channel" binding:"required"`
	WhatsappPayload  WhatsappBroadcastPayload `json:"whatsapp_payload" binding:"omitempty"`
}

type CreateBroadcastResponse struct {
	BroadcastID string `json:"broadcast_id"`
	CreateBroadcastRequest
}

type WhatsappBroadcastPayload struct {
	TemplateName string `json:"template_name" binding:"required"`
}
