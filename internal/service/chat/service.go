package chatsrv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	chatdomain "ichat/internal/domain/chat"
	"ichat/internal/service"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var _ service.ChatService = (*connAdapter)(nil)

func NewConnAdapter() service.ChatService {
	c := connAdapter{}
	c.baseURL = url.URL{
		Scheme: "http",
		Host:   "0.0.0.0:8181",
	}
	c.client = &http.Client{}
	return &c
}

type connAdapter struct {
	baseURL url.URL
	client  *http.Client
	ws      *websocket.Conn
}

// FetchChats implements service.ChatService.
func (c *connAdapter) FetchChats() {

}

// CreateChat implements service.ChatService.
func (c *connAdapter) CreateChat(chatType int, name string) error {
	url := c.baseURL
	url.Path = "api/v1/chats"

	jsonData := map[string]interface{}{
		"name": name,
		"type": chatType,
	}

	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(jsonDataBytes))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed creating chat, status code: %d", resp.StatusCode)
	}

	return nil
}

// GetChats implements service.ChatService.
func (c *connAdapter) GetChats() ([]*chatdomain.Chat, error) {
	url := c.baseURL
	url.Path = "/api/v1/chats"

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chats []*chatdomain.Chat
	err = json.Unmarshal(body, &chats)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

// Close implements service.ChatService.
func (c *connAdapter) Close() error {
	fmt.Println("close")
	if c.ws != nil {
		return c.ws.Close()
	}
	return nil
}

// Connect implements service.ChatService.
func (c *connAdapter) Connect(ctx context.Context) error {
	fmt.Println("connect")
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "0.0.0.0:8181",
		Path:   "/ws",
	}

	var ws *websocket.Conn
	var connErr error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//initial connection
		ws, connErr = websocket.Dial(wsURL.String(), "", "http://0.0.0.0:8181")
		if connErr == nil {
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("try connect")
				ws, connErr = websocket.Dial(wsURL.String(), "", "http://0.0.0.0:8181")
				if connErr == nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
	c.ws = ws
	return nil
}

// ReceiveMessages implements service.ChatService.
func (c *connAdapter) ReceiveMessages(chatID int64) (<-chan *chatdomain.Message, error) {
	msgCh := make(chan *chatdomain.Message)
	go func() {
		defer close(msgCh)
		for {
			var msg chatdomain.Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err != nil {
				fmt.Printf("err.Error(): %v\n", err.Error())
				// Optionally log or handle error here
				return
			}
			fmt.Printf("get msg %+v\n", msg)
			msgCh <- &msg
		}
	}()
	return msgCh, nil
}

// SendMessage implements service.ChatService.
func (c *connAdapter) SendMessage(msg chatdomain.Message) error {
	err := websocket.JSON.Send(c.ws, msg)
	return err
}
