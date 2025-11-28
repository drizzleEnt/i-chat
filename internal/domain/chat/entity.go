package chatdomain

type Chat struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ActionType string

const (
	ActionSendText   ActionType = "send_text"
	ActionSendBinary ActionType = "send_binary"
	ActionJoinChat   ActionType = "join_chat"
	ActionLeaveChat  ActionType = "leave_chat"
	ActionCreateChat ActionType = "create_chat"
)

type Message struct {
	Action   string `json:"action"`
	Content  string `json:"content"`
	SenderID string `json:"sender"`
	ChatID   int64  `json:"chat_id"`
}

type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
