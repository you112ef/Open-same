package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client

	// Content-specific rooms
	rooms map[string]map[*Client]bool

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Client registered: %s", client.ID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				// Remove client from all rooms
				for roomID, clients := range h.rooms {
					if clients[client] {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: %s", client.ID)

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// JoinRoom adds a client to a specific content room
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true

	// Notify other clients in the room
	joinMessage := Message{
		Type:      "user_joined",
		RoomID:    roomID,
		UserID:    client.UserID,
		Username:  client.Username,
		Timestamp: time.Now(),
	}

	h.broadcastToRoom(roomID, joinMessage)
}

// LeaveRoom removes a client from a specific content room
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if clients, exists := h.rooms[roomID]; exists {
		if clients[client] {
			delete(clients, client)
			
			// Notify other clients in the room
			leaveMessage := Message{
				Type:      "user_left",
				RoomID:    roomID,
				UserID:    client.UserID,
				Username:  client.Username,
				Timestamp: time.Now(),
			}

			h.broadcastToRoom(roomID, leaveMessage)

			// Remove room if empty
			if len(clients) == 0 {
				delete(h.rooms, roomID)
			}
		}
	}
}

// BroadcastToRoom sends a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, exists := h.rooms[roomID]; exists {
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}

		for client := range clients {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// broadcastToRoom is an internal method for broadcasting to a room
func (h *Hub) broadcastToRoom(roomID string, message Message) {
	if clients, exists := h.rooms[roomID]; exists {
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}

		for client := range clients {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// GetRoomClients returns all clients in a specific room
func (h *Hub) GetRoomClients(roomID string) []*Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, exists := h.rooms[roomID]; exists {
		result := make([]*Client, 0, len(clients))
		for client := range clients {
			result = append(result, client)
		}
		return result
	}
	return []*Client{}
}

// GetRoomCount returns the number of clients in a specific room
func (h *Hub) GetRoomCount(roomID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, exists := h.rooms[roomID]; exists {
		return len(clients)
	}
	return 0
}

// GetTotalClients returns the total number of connected clients
func (h *Hub) GetTotalClients() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetTotalRooms returns the total number of active rooms
func (h *Hub) GetTotalRooms() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.rooms)
}

// BroadcastToUser sends a message to a specific user across all their connections
func (h *Hub) BroadcastToUser(userID string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- messageBytes
}