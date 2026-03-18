package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    _ "github.com/jackc/pgx/v5/stdlib"
)

var (
	db        *sql.DB
	clients   = make(map[*websocket.Conn]bool)
	mu        sync.Mutex               // Защита карты клиентов от гонок данных
	broadcast = make(chan Message)     // Канал для рассылки
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	var err error
	// Подключение к БД
	db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}

	r := gin.Default()

	// Раздача статики
	r.StaticFile("/", "./static/index.html")

	// Эндпоинт для истории сообщений
	r.GET("/history", func(c *gin.Context) {
		rows, err := db.Query("SELECT username, content FROM messages ORDER BY id DESC LIMIT 50")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var history []Message
		for rows.Next() {
			var m Message
			if err := rows.Scan(&m.Username, &m.Content); err != nil {
				continue
			}
			history = append(history, m)
		}
		
		// Разворачиваем срез, чтобы старые сообщения были сверху
		for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
			history[i], history[j] = history[j], history[i]
		}

		c.JSON(200, history)
	})

	// WebSocket роут
	r.GET("/ws", handleConnections)

	// Воркер рассылки
	go handleMessages()

	log.Println("Chat started on :8080")
	r.Run(":8080")
}

func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		// Сохраняем в Postgres
		_, err := db.Exec("INSERT INTO messages (username, content) VALUES ($1, $2)", msg.Username, msg.Content)
		if err != nil {
			log.Println("DB Write Error:", err)
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("WS Write Error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
