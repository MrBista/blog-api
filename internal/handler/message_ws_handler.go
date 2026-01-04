package handler

import (
	"sync"

	"github.com/MrBista/blog-api/internal/services"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	Conn   *websocket.Conn
	RoomID int
	UserID int
}

type ChatWSHandler struct {
	MessageService services.MessageService
}

func NewChatWsHandler(messageService services.MessageService) *ChatWSHandler {
	return &ChatWSHandler{
		MessageService: messageService,
	}
}

var (
	rooms = make(map[uint]map[*WSClient]bool)
	mu    sync.Mutex
)

func (h *ChatWSHandler) Handle(c *websocket.Conn) {

	roomID := c.Locals("room_id").(int)
	userID := c.Locals("user_id").(int)

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

		msg, err := h.MessageService.SaveMessage(payload.Message, roomID, userID)
		if err != nil {
			continue
		}

		mu.Lock()
		for cl := range rooms[roomID] {
			cl.Conn.WriteJSON(msg)
		}
		mu.Unlock()
	}
}
