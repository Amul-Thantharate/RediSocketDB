package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global variables
var ctx = context.Background()
var rdb *redis.Client
var db *gorm.DB

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Channel   string `gorm:"index"`
	Content   string
	Timestamp int64 `gorm:"autoCreateTime"`
}

func init() {
	var err error

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	dsn := "nora:root@tcp(localhost:3306)/chatdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}

	if err := db.AutoMigrate(&Message{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
}

func storeMessage(channel, content string) error {
	message := Message{Channel: channel, Content: content}
	return db.Create(&message).Error
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	channel := r.FormValue("channel")
	message := r.FormValue("message")

	if channel == "" || message == "" {
		http.Error(w, "Channel and message are required", http.StatusBadRequest)
		return
	}

	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		http.Error(w, "Error publishing message", http.StatusInternalServerError)
		return
	}

	err = storeMessage(channel, message)
	if err != nil {
		http.Error(w, "Error storing message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message published and stored in Mysql: %s", message)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	channels := r.URL.Query()["channel"] // Get multiple channel parameters
	if len(channels) == 0 {
		http.Error(w, "At least one 'channel' parameter is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	pubsub := rdb.Subscribe(ctx, channels...) // Subscribe to multiple channels
	defer pubsub.Close()

	fmt.Printf("Client subscribed to channels: %v\n", channels)

	ch := pubsub.Channel()
	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[%s] %s", msg.Channel, msg.Payload))); err != nil {
			fmt.Printf("WebSocket error: %v\n", err)
			break
		}
	}
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/publish", PublishHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

	server := &http.Server{Addr: ":8080"}

	go func() {
		fmt.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-signalChan
	fmt.Println("Shutting down server...")
	server.Close()
	rdb.Close()
	fmt.Println("Server shutdown complete")
}
