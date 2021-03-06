package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fehepe/chatbot-challenge/pkg/stock"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline   = []byte{'\n'}
	space     = []byte{' '}
	CmdResult = ""
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// send chan []byte
	send chan Message

	username string
	room     string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- Message{Username: c.username, Room: c.room, Message: string(message), Type: "chatmessage", Time: time.Now().Format("3:04 pm")}

		if strings.Contains(string(message), "/") {
			cmd := string(message)
			if strings.HasPrefix(cmd, "/stock=") {
				code := strings.ReplaceAll(cmd, "/stock=", "")

				cmdStock, err := stock.GetStockFromAPI(code)
				if err != nil {
					log.Print("GetStockFromAPI error:" + err.Error())

				}

				CmdResult = cmdStock.Response
				stockCmd := Message{Username: "ChatBot", Room: c.room,
					Message: "stock", Type: "botmessage", Time: time.Now().Format("3:04 pm")}
				c.hub.broadcast <- stockCmd
			} else {
				time.Sleep(2 * time.Second)
				wrongCmd := Message{Username: "ChatBot", Room: c.room,
					Message: "wrong", Type: "botmessage", Time: time.Now().Format("3:04 pm")}
				c.hub.broadcast <- wrongCmd
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.broadcast <- Message{Username: c.username, Room: c.room, Message: "leave", Type: "botmessage", Time: time.Now().Format("3:04 pm")}

		ticker.Stop()
		c.conn.Close()
	}()
	// broadcast welcome messages
	c.hub.broadcast <- Message{Username: c.username, Room: c.room, Message: "welcome", Type: "botmessage", Time: time.Now().Format("3:04 pm")}

	for {
		select {
		// Message send from hub.broadcast
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// update user list
			if msg.Type == "botmessage" {
				var userlist []string
				for client := range c.hub.clients {
					userlist = append(userlist, client.username)
				}
				msg.Userlist = userlist
			}

			b, _ := json.Marshal(msg)
			fmt.Println(msg)
			w.Write(b)

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

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	var username, room string
	if _, ok := r.URL.Query()["username"]; ok {
		fmt.Println(r.URL.Query()["username"])
		username = r.URL.Query()["username"][0]
	}
	if _, ok := r.URL.Query()["room"]; ok {
		fmt.Println(r.URL.Query()["room"])
		room = r.URL.Query()["room"][0]
	}

	client := &Client{hub: hub, conn: conn, send: make(chan Message), username: username, room: room}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}
