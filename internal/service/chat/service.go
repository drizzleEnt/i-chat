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
	c := &connAdapter{
		baseURL: url.URL{
			Scheme: "http",
			Host:   "0.0.0.0:8181",
		},
		client: &http.Client{},
		status: service.StatusDisconnected,
	}
	return c
}

type connAdapter struct {
	baseURL url.URL
	client  *http.Client

	mu     sync.RWMutex
	ws     *websocket.Conn
	status service.ConnStatus
	ctx    context.Context
	cancel context.CancelFunc
}

// FetchChats implements service.ChatService.
func (c *connAdapter) FetchChats() {
}

// CreateChat implements service.ChatService.
func (c *connAdapter) CreateChat(chatType int, name string) error {
	u := c.baseURL
	u.Path = "api/v1/chats"

	jsonData := map[string]interface{}{
		"name": name,
		"type": chatType,
	}

	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonDataBytes))
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
	u := c.baseURL
	u.Path = "/api/v1/chats"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	if c.ws != nil {
		err := c.ws.Close()
		c.ws = nil
		c.status = service.StatusDisconnected
		return err
	}
	c.status = service.StatusDisconnected
	return nil
}

// Connect implements service.ChatService.
// Starts connection attempt in background and returns immediately.
// Connection status can be checked with GetStatus().
func (c *connAdapter) Connect(ctx context.Context) error {
	c.mu.Lock()
	// Create child context for connection management
	connCtx, cancel := context.WithCancel(ctx)
	c.ctx = connCtx
	c.cancel = cancel
	c.status = service.StatusConnecting
	c.mu.Unlock()

	go c.connectLoop()

	return nil
}

func (c *connAdapter) connectLoop() {
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "0.0.0.0:8181",
		Path:   "/ws",
	}

	reconnectDelay := 1 * time.Second
	maxReconnectDelay := 30 * time.Second

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		ws, err := websocket.Dial(wsURL.String(), "", "http://0.0.0.0:8181")

		c.mu.Lock()
		if err == nil {
			c.ws = ws
			c.status = service.StatusConnected
			reconnectDelay = 1 * time.Second
			c.mu.Unlock()
			fmt.Println("Connected to server")

			// Wait for connection to close
			<-c.ctx.Done()

			c.mu.Lock()
			c.ws = nil
			c.status = service.StatusDisconnected
			c.mu.Unlock()
			fmt.Println("Disconnected from server")
		} else {
			c.status = service.StatusConnecting
			c.mu.Unlock()
			fmt.Printf("Connection failed: %v. Retrying in %v...\n", err, reconnectDelay)

			// Wait before retrying
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(reconnectDelay):
				// Exponential backoff
				reconnectDelay *= 2
				if reconnectDelay > maxReconnectDelay {
					reconnectDelay = maxReconnectDelay
				}
			}
		}
	}
}

// GetStatus returns current connection status.
func (c *connAdapter) GetStatus() service.ConnStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// IsConnected returns true if currently connected.
func (c *connAdapter) IsConnected() bool {
	return c.GetStatus() == service.StatusConnected
}

// getWS returns the websocket connection safely.
func (c *connAdapter) getWS() *websocket.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ws
}

// ReceiveMessages implements service.ChatService.
func (c *connAdapter) ReceiveMessages(chatID int64) (<-chan *chatdomain.Message, error) {
	msgCh := make(chan *chatdomain.Message)

	go func() {
		defer close(msgCh)

		for {
			ws := c.getWS()
			if ws == nil {
				select {
				case <-c.ctx.Done():
					return
				case <-time.After(500 * time.Millisecond):
					continue
				}
			}

			var msg chatdomain.Message
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				// Connection lost, will retry when reconnected
				select {
				case <-c.ctx.Done():
					return
				case <-time.After(500 * time.Millisecond):
					continue
				}
			}

			select {
			case msgCh <- &msg:
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return msgCh, nil
}

// SendMessage implements service.ChatService.
func (c *connAdapter) SendMessage(msg chatdomain.Message) error {
	ws := c.getWS()
	if ws == nil {
		return fmt.Errorf("not connected")
	}

	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// StatusWatcher returns a channel that receives status updates.
func (c *connAdapter) StatusWatcher() <-chan service.ConnStatus {
	statusCh := make(chan service.ConnStatus, 10)

	go func() {
		lastStatus := service.StatusDisconnected
		for {
			select {
			case <-c.ctx.Done():
				close(statusCh)
				return
			default:
				currentStatus := c.GetStatus()
				if currentStatus != lastStatus {
					select {
					case statusCh <- currentStatus:
						lastStatus = currentStatus
					default:
						// Channel full, skip
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return statusCh
}
