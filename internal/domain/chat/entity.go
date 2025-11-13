package chatdomain

type Chat struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	SenderID string `json:"sender"`
	Content  string `json:"content"`
	ChatID   string `json:"chat_id"`
}

type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
