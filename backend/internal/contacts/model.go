package contacts

type ContactCard struct {
	ID          string   `json:"id" binding:"ommitempty"`
	Name        string   `json:"name" binding:"required"`
	Phone       string   `json:"phone" binding:"ommitempty"`
	Email       string   `json:"email" binding:"ommitempty"`
	Address     string   `json:"address" binding:"ommitempty"`
	Tags        []string `json:"tags" binding:"ommitempty"`
	HasWhatsapp bool     `json:"has_whatsapp" binding:"ommitempty"`
	HasTelegram bool     `json:"has_telegram" binding:"ommitempty"`
}

type ContactVariable struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type CreateContactRequest struct {
	ContactCard ContactCard       `json:"contact_card" binding:"required"`
	Variables   []ContactVariable `json:"variables" binding:"ommitempty"`
}
