// Package websocket — WebSocket hub для табло в реальном времени.
package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Hub управляет WebSocket-соединениями.
type Hub struct {
	logger     *zap.Logger
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client представляет WebSocket-клиента.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Message — структура сообщения для клиентов (поля: 8-байтные, затем 16-байтные, затем time для fieldalignment).
type Message struct {
	Timestamp    time.Time              `json:"timestamp"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Type         string                 `json:"type"`
	TripID       string                 `json:"trip_id,omitempty"`
	Status       string                 `json:"status,omitempty"`
	DelayMinutes int                    `json:"delay_minutes,omitempty"`
}

// NewHub создаёт новый Hub.
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// Run запускает hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Debug("Client connected", zap.Int("total_clients", len(h.clients)))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Debug("Client disconnected", zap.Int("total_clients", len(h.clients)))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast отправляет сообщение всем клиентам.
func (h *Hub) Broadcast(msg *Message) {
	msg.Timestamp = time.Now()
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal message", zap.Error(err))
		return
	}
	h.broadcast <- data
}

// GetClientCount возвращает количество подключенных клиентов.
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// readPump читает сообщения от клиента.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		if err := c.conn.Close(); err != nil {
			c.hub.logger.Warn("failed to close conn in readPump", zap.Error(err))
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.hub.logger.Warn("SetReadDeadline", zap.Error(err))
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			c.hub.logger.Warn("SetReadDeadline in pong", zap.Error(err))
		}
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// writePump отправляет сообщения клиенту.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			c.hub.logger.Warn("failed to close conn in writePump", zap.Error(err))
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.hub.logger.Warn("SetWriteDeadline", zap.Error(err))
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.hub.logger.Warn("WriteMessage close", zap.Error(err))
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				c.hub.logger.Warn("Write", zap.Error(err))
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.hub.logger.Warn("SetWriteDeadline ping", zap.Error(err))
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs обрабатывает WebSocket-соединение.
func ServeWs(hub *Hub, conn *websocket.Conn) {
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
