package chatsrv

import "ichat/internal/service"

var _ service.ChatService = (*connAdapter)(nil)

func NewConnAdapter() service.ChatService {
	return &connAdapter{}
}

type connAdapter struct {
}

// Close implements service.ChatService.
func (c *connAdapter) Close() error {
	panic("unimplemented")
}

// Connect implements service.ChatService.
func (c *connAdapter) Connect() error {
	panic("unimplemented")
}

// ReceiveMessages implements service.ChatService.
func (c *connAdapter) ReceiveMessages(chatID string) ([]string, error) {
	panic("unimplemented")
}

// SendMessage implements service.ChatService.
func (c *connAdapter) SendMessage(chatID string, message string) error {
	panic("unimplemented")
}
