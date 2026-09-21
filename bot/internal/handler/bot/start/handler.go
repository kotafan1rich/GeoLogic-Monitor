package start

import "context"

type UserService interface {
	RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error
}

type handler struct {
	userService UserService
}

func NewStartHandler(userService UserService) *handler {
	return &handler{
		userService: userService,
	}
}

func (h *handler) Start() {

}
