package hub

import (
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	Username string
	Message  string
	Type     string
	Userlist []string
	Room     string
	Time     string
}

func NewHub() *Hub {
	return &Hub{
		//broadcast:  make(chan []byte),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.Type == "botmessage" {
				for client := range h.clients {
					if client.room == message.Room {
						if client.username == message.Username {
							welcomemsg := Message{Username: "ChatBot", Room: client.room,
								Message: "Welcome to the chat room " + message.Username, Type: "botmessage", Time: time.Now().Format("3:04 pm")}
							select {
							case client.send <- welcomemsg:
							default:
								close(client.send)
								delete(h.clients, client)
							}
						} else {
							var msg Message
							if message.Message == "welcome" {
								msg = Message{Username: "ChatBot", Room: client.room,
									Message: message.Username + " has entered the chat", Type: "botmessage", Time: time.Now().Format("3:04 pm")}
							} else if message.Message == "leave" {
								msg = Message{Username: "ChatBot", Room: client.room,
									Message: message.Username + " has left the chat", Type: "botmessage", Time: time.Now().Format("3:04 pm")}
							} else if message.Message == "stock" {
								msg = Message{Username: "ChatBot", Room: client.room,
									Message: CmdResult, Type: "botmessage", Time: time.Now().Format("3:04 pm")}
							} else if message.Message == "wrong" {
								msg = Message{Username: "ChatBot", Room: client.room,
									Message: "I do not have this command. Try another.", Type: "botmessage", Time: time.Now().Format("3:04 pm")}
							}
							select {
							case client.send <- msg:
							default:
								close(client.send)
								delete(h.clients, client)
							}
						}
					}
				}
			} else if message.Type == "chatmessage" {
				for client := range h.clients {
					if client.room == message.Room {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}
