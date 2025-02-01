package contacts

type contactsService struct {
	Repository contactsRepository
}

func NewContactsService(repo *contactsRepository) *contactsService {
	return &contactsService{
		Repository: *repo,
	}
}

func (s *contactsService) CreateContacts(data []CreateContactRequest) ([]ContactCard, error) {
	contacts := []ContactCard{}
	for _, contact := range data {
		contacts = append(contacts, ContactCard{
			Name:        contact.Name,
			Phone:       contact.Phone,
			Email:       contact.Email,
			Address:     contact.Address,
			Tags:        contact.Tags,
			HasWhatsapp: contact.HasWhatsapp,
			HasTelegram: contact.HasTelegram,
		})
	}

	createContacts, createContactsError := s.Repository.CreateContact(contacts)
	if createContactsError != nil {
		return []ContactCard{}, createContactsError
	}

	return createContacts, nil
}
