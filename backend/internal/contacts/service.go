package contacts

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contactsService struct {
	Repository contactsRepository
}

func NewContactsService(repo *contactsRepository) *contactsService {
	return &contactsService{
		Repository: *repo,
	}
}

func (s *contactsService) CreateContacts(data []CreateContactRequest, userId string) ([]ContactCard, error) {
	contacts := []ContactCard{}
	variables := []ContactVariableDB{}
	createdBy, createdByError := primitive.ObjectIDFromHex(userId)
	if createdByError != nil {
		log.Println("[ContactsService.CreateContacts] Error creating ObjectID from hex:", createdByError)
		return []ContactCard{}, createdByError
	}

	now := time.Now().Unix()

	for _, contact := range data {
		newContactId := primitive.NewObjectID()
		newContact := ContactCard{
			ID:          newContactId,
			Name:        contact.Name,
			Phone:       contact.Phone,
			Email:       contact.Email,
			Address:     contact.Address,
			Tags:        contact.Tags,
			HasWhatsapp: contact.HasWhatsapp,
			HasTelegram: contact.HasTelegram,
			CreatedBy:   createdBy,
			CreatedAt:   now,
			Variables:   contact.Variables,
		}
		contacts = append(contacts, newContact)

		if len(contact.Variables) == 0 {
			continue
		}

		for _, variable := range contact.Variables {
			variables = append(variables, ContactVariableDB{
				ContactId: newContactId,
				ContactVariable: ContactVariable{
					Key:   variable.Key,
					Value: variable.Value,
				},
			})
		}
	}

	createContacts, createContactsError := s.Repository.CreateContact(contacts)
	if createContactsError != nil {
		return []ContactCard{}, createContactsError
	}

	_, createVariablesError := s.Repository.CreateVariable(variables)
	if createVariablesError != nil {
		return []ContactCard{}, createVariablesError
	}

	return createContacts, nil
}

// func (s *contactsService) GetContacts() ([]GetContactResponse, error) {
// 	return s.Repository.GetContacts()
// }
