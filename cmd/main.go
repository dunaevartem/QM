package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
    "github.com/gorilla/websocket"
)

var (
	db       *sql.DB
	clients  = make(map[*websocket.Conn]bool) // Список подключенных клиентов
	broadcast = make(chan Message)           // Канал для рассылки
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	// Подключение к БД через env (для K8s)
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.GET("/ws", handleConnections)

	go handleMessages() // Воркер рассылки

	log.Println("Chat started on :8080")
	r.Run(":8080")
}

func handleConnections(c *gin.Context) {
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			delete(clients, ws)
			break
		}
		// Сохраняем в Postgres
		db.Exec("INSERT INTO messages (username, content) VALUES ($1, $2)", msg.Username, msg.Content)
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}
