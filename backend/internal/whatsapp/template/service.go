package template

func NewTermplateService(r *templateRepository) *templateService {
	return &templateService{
		Repository: r,
	}
}
