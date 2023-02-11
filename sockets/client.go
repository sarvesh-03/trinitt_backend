package sockets

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/trinitt/utils"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	mux  sync.Mutex

	// The fields below are for use by the hub.
	userId uint
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Time{})
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var M Message
		err := c.conn.ReadJSON(&M)
		if err != nil {
			log.Println(err)
			return
		}

		if M.Type == "AUTH_USER" {
			token := M.Data.(string)

			id, err := utils.GetUserIDFromToken(token)

			fmt.Println(id, "id")
			if err != nil {
				log.Println(err)
				return
			}

			c.userId = id
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			c.mux.Lock()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.mux.Unlock()
				return
			}
			c.mux.Unlock()
		}
	}
}

var wsHub *Hub

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn}
	client.hub.register <- client
	wsHub = hub

	go client.writePump()
	go client.readPump()
}
