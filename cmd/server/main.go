package main

import (
	"fmiis/internal/app"
	"fmiis/internal/config"
	"log"
	"net/http"
	"os"
	"path/filepath"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config :", err)
	}

	/* -----------------------------------------------------------
	   🔌 SIMPLE & CLEAN SOCKET.IO SERVER (like Node.js)
	----------------------------------------------------------- */
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool { return true },
			},
		},
	})

	// ---- Event Handlers ----
	server.OnConnect("/", func(conn socketio.Conn) error {
		log.Println("🔌 Socket connected:", conn.ID())
		conn.Emit("connected", "You are connected!")
		return nil
	})

	server.OnEvent("/", "ping", func(conn socketio.Conn, msg string) {
		log.Println("📩 PING received:", msg)
		conn.Emit("pong", "PONG from server!")
	})

	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		log.Println("❌ Socket disconnected:", conn.ID(), "reason:", reason)
	})

	server.OnError("/", func(conn socketio.Conn, err error) {
		log.Println("🔥 Socket error:", err)
	})

	go server.Serve()
	defer server.Close()

	/* ----------------------------------------------------------- */

	// Create uploads directory
	//Get project root (where main.go is)
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	uploadBaseDir := filepath.Join(projectRoot, "uploads")
	// Or if running from project root directly:
	// uploadBaseDir := filepath.Join(projectRoot, "uploads")

	application, err := app.NewApp(cfg, *server, uploadBaseDir)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// Start server
	application.StartServer()
}
