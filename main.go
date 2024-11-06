package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// web socket upgrader

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client struct to hold connection

type Client struct {
	Conn     *websocket.Conn
	UserName string
}

// creating global var

var (
	clients    = make(map[*Client]bool)
	messages   = []string{}
	messageCha = make(chan string)
	clientMux  sync.Mutex
)

func main() {

	// create router

	router := gin.Default()

	// for js file
	router.Static("/static", "./static")

	// for html

	router.LoadHTMLGlob("template/*")
	router.GET("/", showChatPage)

	// web scoket handler
	router.GET("/ws", handleWebSocket)

	// for message broad casting

	go broadCastMessages()

	// for running the server

	router.Run(":8080")

}

// chat page

func showChatPage(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)
}

// handel web socket connection

func handleWebSocket(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Faild to upgrade web socket :", err)

		return
	}

	defer conn.Close()

	// creating a new client

	client := &Client{Conn: conn}
	clientMux.Lock()
	clients[client] = true
	clientMux.Unlock()

	// to persist chat histroy

	for _, message := range messages {

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Faild to persist the chat histroy :", err)
			client.Conn.Close()
			return

		}
	}

	// wait for the message

	for {

		var msg map[string]string

		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("error while reading the message:", err)
			break
		}

		//Handel message based on action

		switch msg["action"] {

		// when the user join the chat
		case "join":

			client.UserName = msg["username"]

			broadCastMessage(fmt.Sprintf("%s has join the chat", client.UserName))

			// for broadcasting the message
		case "message":

			broadCastMessage(fmt.Sprintf("%s:%s", client.UserName, msg["message"]))

			// when user leave the chat room

		case "leave":

			broadCastMessage(fmt.Sprintf("%s has leave the chat ", client.UserName))

			clientMux.Lock()
			delete(clients, client)
			clientMux.Unlock()
			client.Conn.Close()
			return
		}

	}
}

func broadCastMessages() {

	for {

		msg := <-messageCha

		messages = append(messages, msg)

		clientMux.Lock()
		for c := range clients {

			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))

			if err != nil {

				fmt.Println("error while sending the message:", err)
				c.Conn.Close()
				delete(clients, c)
			}
		}
		clientMux.Unlock()
	}

}

// for broadcasting message to all the client  and put the message in a channel
func broadCastMessage(msg string) {

	messageCha <- msg
}
