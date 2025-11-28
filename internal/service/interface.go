package service

import (
	"context"
	chatdomain "ichat/internal/domain/chat"
)

type ChatService interface {
	Connect(ctx context.Context) error
	Close() error
	SendMessage(msg chatdomain.Message) error
	ReceiveMessages(chatID int64) (<-chan *chatdomain.Message, error)
	GetChats() ([]*chatdomain.Chat, error)
	CreateChat(chatType int, name string) error
	FetchChats()
}

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}
