package api

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mhghw/fara-message/db"
)

type Hub struct {
	clients     map[uint64]*Client
	chatClients map[uint64][]*Client
	broadcast   chan db.Message
	register    chan *Client
	unregister  chan *Client
}

type Client struct {
	userID  uint64 /////bebin mishe db ro bardasht
	chatID  uint64
	hub     *Hub
	conn    *websocket.Conn
	send    chan db.Message /////bebin mishe db ro bardasht
	receive chan db.Message /////bebin mishe db ro bardasht
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
}

func (h *Hub) Run() {
	log.Printf("hereeee11111")  ////////////////////////////////////////////////////////
	for {
		log.Printf("hereeee2000000")  ////////////////////////////////////////////////////////

		select {
		case client := <-h.register:
			log.Printf("hereeee100000000000")  ////////////////////////////////////////////////////////

			err := h.AddClientToHub(client)
			if err != nil {
				log.Printf("error adding client to hub:%v", err)
				// c.Status(400)                ????????????????????????????????????????????????
				return
			}
		case client := <-h.unregister:
			h.RemoveClient(client)
		case message := <-h.broadcast:
			err := h.SendMessageToClients(message)
			if err != nil {
				log.Printf("error adding client to hub:%v", err)
				// c.Status(400)                ????????????????????????????????????????????????
				return
			}
		}
	}
}

func (h *Hub) AddClientToHub(c *Client) error {
	log.Printf("hereeee22222")  ////////////////////////////////////////////////////////

	flag, err := db.Mysql.IsAChatContact(c.userID, c.chatID)
	if err != nil {
		return err
	}
	if !flag {
		return errors.New("this ID is not a member of this chat")
	}
	h.clients[c.userID] = c
	h.chatClients[c.chatID] = append(h.chatClients[c.chatID], c) //check!!!!!!!!!!!!!!!!!!!!!!!!
	return nil
}

func (h *Hub) RemoveClient(c *Client) {
	log.Printf("hereeee3333333")  ////////////////////////////////////////////////////////

	delete(h.clients, c.userID)
	//inja check beshe!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1
	if userIDs, ok := h.chatClients[c.chatID]; ok {
		for i, id := range userIDs {
			if id == c {
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
	log.Printf("hereeee444444")  ////////////////////////////////////////////////////////

	//chon tooye add client to hub yedor check kardam ke karbar ozve chat hast ya na, inja check nekardam
	err := db.Mysql.SendMessage(message)
	if err != nil {
		return err
	}
	for key, value := range h.chatClients {
		if key == message.ChatID {
			for _, v := range value {
				log.Printf("i send message to :%v",value)  ////////////////////////////////////////////////////////

				v.receive <- message
			}
		}
	}
	return nil
}

func ServeWs(hub *Hub, c *gin.Context) {
	log.Printf("hereeee55555")  ////////////////////////////////////////////////////////
	log.Printf("this is hub:%v",hub)

	userID, err := GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		log.Printf("error get user ID:%v", err)
		c.Status(400)
		return
	}
	intOfUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	chatID := c.Param("chatID")
	intOfChatID, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		log.Printf("error converting to int:%v", err)
		c.Status(400)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("error upgrading connection to websocket:%v", err)
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
	// users[client.user.ID] = conn
	client.hub.register <- client
	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) ReadPump() {
	log.Printf("hereeee66666")  ////////////////////////////////////////////////////////

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var message db.Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error reading message from websocket:%v", err)
			// c.Status(400)     ?????????????????????????????????????????
			return
		}

		message.SenderID = c.userID
		messageID, err := generateID()
		if err != nil {
			log.Printf("error in generating ID:%v", err)
			// c.Status(400)                ?????????????????????????????????????
			return
		}
		message.ID = messageID
		message.Time = time.Now()
		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	log.Printf("hereeee77777")  ////////////////////////////////////////////////////////

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for{
		message:=<-c.receive
		err := c.conn.WriteJSON(message)
		if err != nil {
			log.Printf("error writing message to websocket: %v", err)
			// c.Status(400)                ?????????????????????????????????????
			return
		}
	}
}

// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.hub.unregister <- c
// 		c.conn.Close()
// 	}()
// 	c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		var message db.Message
// 		err := c.conn.ReadJSON(&message)
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}

// 		message.SenderID = c.userID
// 		messageID, err := generateID()
// 		if err != nil {
// 			log.Printf("error in generating ID: %v", err)
// 			return
// 		}
// 		message.ID = messageID
// 		message.Time = time.Now()
// 		c.hub.broadcast <- message
// 	}
// }

// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.hub.unregister <- c
// 		c.conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message := <-c.receive:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteJSON(message); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
