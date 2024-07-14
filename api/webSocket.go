package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mhghw/fara-message/db"
)

const (
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
)

type Hub struct {
	clients     map[uint64]*Client
	chatClients map[uint64][]*Client
	broadcast   chan db.Message
	register    chan *Client
	unregister  chan *Client
}

type Client struct {
	userID  uint64
	chatID  uint64
	hub     *Hub
	conn    *websocket.Conn
	send    chan db.Message
	receive chan db.Message
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[uint64]*Client),
		chatClients: make(map[uint64][]*Client),
		broadcast:   make(chan db.Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			err := h.AddClientToHub(client)
			if err != nil {
				log.Printf("Error adding client to hub: %v", err)
				return
			}
		case client := <-h.unregister:
			h.RemoveClient(client)
		case message := <-h.broadcast:
			err := h.SendMessageToClients(message)
			if err != nil {
				log.Printf("Error sending message to clients: %v", err)
				return
			}
		}
	}
}

func (h *Hub) AddClientToHub(c *Client) error {
	flag, err := db.Mysql.IsAChatContact(c.userID, c.chatID)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New("this ID is not a member of this chat")
	}
	h.clients[c.userID] = c
	h.chatClients[c.chatID] = append(h.chatClients[c.chatID], c)
	return nil
}

func (h *Hub) RemoveClient(c *Client) {
	delete(h.clients, c.userID)
	if userIDs, ok := h.chatClients[c.chatID]; ok {
		for i, client := range userIDs {
			if client == c {
				h.chatClients[c.chatID] = append(userIDs[:i], userIDs[i+1:]...)
				break
			}
		}
		if len(h.chatClients[c.chatID]) == 0 {
			delete(h.chatClients, c.chatID)
		}
	}
	c.conn.Close()
	close(c.send)
	close(c.receive)
}

func (h *Hub) SendMessageToClients(message db.Message) error {
	err := db.Mysql.SendMessage(message)
	if err != nil {
		return err
	}
	for key, clients := range h.chatClients {
		if key == message.ChatID {
			for _, client := range clients {
				client.receive <- message
			}
		}
	}
	return nil
}

func ServeWs(hub *Hub, c *gin.Context) {
	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.Status(400)
		return
	}
	intOfUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf("Error converting user ID to int: %v", err)
		c.Status(400)
		return
	}
	chatID := c.Param("chatID")
	intOfChatID, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		log.Printf("Error converting chat ID to int: %v", err)
		c.Status(400)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		c.Status(400)
		return
	}
	client := &Client{
		userID:  intOfUserID,
		chatID:  intOfChatID,
		hub:     hub,
		conn:    conn,
		send:    make(chan db.Message),
		receive: make(chan db.Message),
	}
	client.hub.register <- client
	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message db.Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		message.SenderID = c.userID
		message.ChatID = c.chatID
		messageID, err := generateID()
		if err != nil {
			log.Printf("Error generating message ID: %v", err)
			return
		}
		message.ID = messageID
		message.Time = time.Now()
		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message := <-c.receive:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteJSON(message); err != nil {
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
