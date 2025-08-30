package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should implement proper origin checking
		return true
	},
}

// Client represents a connected WebSocket client
type Client struct {
	// Unique identifier for the client
	ID string

	// User information
	UserID   string
	Username string

	// The hub
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Current room
	currentRoom string
}

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	RoomID    string                 `json:"room_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// HandleWebSocket handles the WebSocket connection upgrade and client registration
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create new client
	client := &Client{
		ID:       uuid.New().String(),
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		UserID:   r.URL.Query().Get("user_id"),
		Username: r.URL.Query().Get("username"),
	}

	// Register client with hub
	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle message based on type
		c.handleMessage(msg)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "join_room":
		c.handleJoinRoom(msg)
	case "leave_room":
		c.handleLeaveRoom(msg)
	case "content_change":
		c.handleContentChange(msg)
	case "cursor_move":
		c.handleCursorMove(msg)
	case "selection_change":
		c.handleSelectionChange(msg)
	case "chat_message":
		c.handleChatMessage(msg)
	case "ping":
		c.handlePing()
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleJoinRoom handles room joining
func (c *Client) handleJoinRoom(msg Message) {
	roomID := msg.RoomID
	if roomID == "" {
		return
	}

	// Leave current room if any
	if c.currentRoom != "" {
		c.hub.LeaveRoom(c, c.currentRoom)
	}

	// Join new room
	c.currentRoom = roomID
	c.hub.JoinRoom(c, roomID)

	// Send confirmation
	response := Message{
		Type:      "room_joined",
		RoomID:    roomID,
		UserID:    c.UserID,
		Username:  c.Username,
		Timestamp: time.Now(),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes
}

// handleLeaveRoom handles room leaving
func (c *Client) handleLeaveRoom(msg Message) {
	if c.currentRoom != "" {
		c.hub.LeaveRoom(c, c.currentRoom)
		c.currentRoom = ""
	}

	// Send confirmation
	response := Message{
		Type:      "room_left",
		UserID:    c.UserID,
		Username:  c.Username,
		Timestamp: time.Now(),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes
}

// handleContentChange handles content changes
func (c *Client) handleContentChange(msg Message) {
	if c.currentRoom == "" {
		return
	}

	// Broadcast change to other clients in the room
	changeMessage := Message{
		Type:      "content_change",
		RoomID:    c.currentRoom,
		UserID:    c.UserID,
		Username:  c.Username,
		Data:      msg.Data,
		Timestamp: time.Now(),
	}

	c.hub.BroadcastToRoom(c.currentRoom, changeMessage)
}

// handleCursorMove handles cursor movement
func (c *Client) handleCursorMove(msg Message) {
	if c.currentRoom == "" {
		return
	}

	// Broadcast cursor position to other clients in the room
	cursorMessage := Message{
		Type:      "cursor_move",
		RoomID:    c.currentRoom,
		UserID:    c.UserID,
		Username:  c.Username,
		Data:      msg.Data,
		Timestamp: time.Now(),
	}

	c.hub.BroadcastToRoom(c.currentRoom, cursorMessage)
}

// handleSelectionChange handles text selection changes
func (c *Client) handleSelectionChange(msg Message) {
	if c.currentRoom == "" {
		return
	}

	// Broadcast selection to other clients in the room
	selectionMessage := Message{
		Type:      "selection_change",
		RoomID:    c.currentRoom,
		UserID:    c.UserID,
		Username:  c.Username,
		Data:      msg.Data,
		Timestamp: time.Now(),
	}

	c.hub.BroadcastToRoom(c.currentRoom, selectionMessage)
}

// handleChatMessage handles chat messages
func (c *Client) handleChatMessage(msg Message) {
	if c.currentRoom == "" {
		return
	}

	// Broadcast chat message to other clients in the room
	chatMessage := Message{
		Type:      "chat_message",
		RoomID:    c.currentRoom,
		UserID:    c.UserID,
		Username:  c.Username,
		Content:   msg.Content,
		Timestamp: time.Now(),
	}

	c.hub.BroadcastToRoom(c.currentRoom, chatMessage)
}

// handlePing handles ping messages
func (c *Client) handlePing() {
	response := Message{
		Type:      "pong",
		Timestamp: time.Now(),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes
}

// SendMessage sends a message to this specific client
func (c *Client) SendMessage(message Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		close(c.send)
		delete(c.hub.clients, c)
	}
}

// GetCurrentRoom returns the current room ID
func (c *Client) GetCurrentRoom() string {
	return c.currentRoom
}

// IsInRoom checks if the client is in a specific room
func (c *Client) IsInRoom(roomID string) bool {
	return c.currentRoom == roomID
}