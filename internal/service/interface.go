package service

import chatdomain "ichat/internal/domain/chat"

type ChatService interface {
	Connect() error
	Close() error
	SendMessage(msg chatdomain.Message) error
	ReceiveMessages(chatID string) (<-chan *chatdomain.Message, error)
	GetChats() ([]*chatdomain.Chat, error)
}

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}
