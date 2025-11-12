package chatdomain

type Chat struct {
	ID       string
	Name     string
	Members  []Member
	Messages []Message
}

type Message struct {
	SenderID  string
	Content   string
	Timestamp int64
}

type Member struct {
	ID   string
	Name string
}
