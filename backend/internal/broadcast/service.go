package broadcast

type broadcastService struct {
	Repostory broadcastRepository
}

func NewBroadcastService(repo *broadcastRepository) *broadcastService {
	return &broadcastService{
		Repostory: *repo,
	}
}
