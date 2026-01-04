package handler

import (
	"sync"

	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/websocket/v2"
)

type WSClient struct {
	Conn   *websocket.Conn
	RoomID uint
	UserID uint
}

type ChatWSHandler struct {
	MessageService services.MessageService
}

func NewChatWSHandler(messageService services.MessageService) *ChatWSHandler {
	return &ChatWSHandler{
		MessageService: messageService,
	}
}

var (
	rooms = make(map[uint]map[*WSClient]bool)
	mu    sync.RWMutex
)

func (h *ChatWSHandler) Handle(c *websocket.Conn) {
	roomID := c.Locals("room_id").(uint)
	userID := c.Locals("user_id").(uint)

	client := &WSClient{
		Conn:   c,
		RoomID: roomID,
		UserID: userID,
	}

	mu.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*WSClient]bool)
	}
	rooms[roomID][client] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(rooms[roomID], client)
		if len(rooms[roomID]) == 0 {
			delete(rooms, roomID)
		}
		mu.Unlock()
		c.Close()
	}()

	for {
		var payload struct {
			Message string `json:"message"`
		}

		if err := c.ReadJSON(&payload); err != nil {
			break
		}

		if payload.Message == "" {
			continue
		}

		msg, err := h.MessageService.SaveMessage(payload.Message, int(roomID), int(userID))
		if err != nil {
			continue
		}

		mu.RLock()
		clients := make([]*WSClient, 0, len(rooms[roomID]))
		for cl := range rooms[roomID] {
			clients = append(clients, cl)
		}
		mu.RUnlock()

		for _, cl := range clients {
			if err := cl.Conn.WriteJSON(msg); err != nil {
				mu.Lock()
				delete(rooms[roomID], cl)
				mu.Unlock()
				cl.Conn.Close()
			}
		}
	}
}
