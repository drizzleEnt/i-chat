package service

import (
	"context"
	chatdomain "ichat/internal/domain/chat"
)

type ChatService interface {
	Connect(ctx context.Context) ( error)
	Close() error
	SendMessage(msg chatdomain.Message) error
	ReceiveMessages(chatID string) (<-chan *chatdomain.Message, error)
	GetChats() ([]*chatdomain.Chat, error)
	CreateChat(name string) error
}

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}
