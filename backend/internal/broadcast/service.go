package broadcast

type broadcastService struct {
	Repostory broadcastRepository
}

func NewBroadcastService(repo *broadcastRepository) *broadcastService {
	return &broadcastService{
		Repostory: *repo,
	}
}

func (s *broadcastService) CreateWhatsappBroadcast(req CreateBroadcastRequest) (CreateBroadcastResponse, error) {
	return CreateBroadcastResponse{
		BroadcastID:            "123123123",
		CreateBroadcastRequest: req,
	}, nil
}
