package service

type ChatService interface {
	Connect() error
	Close() error
	SendMessage(chatID string, message string) error
	ReceiveMessages(chatID string) ([]string, error)
}

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}
