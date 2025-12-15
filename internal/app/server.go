package app

import (
	"fmt"
	"log"

	"fmiis/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (a *App) StartServer() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware())

	RegisterRoutes(r, a)

	// /* -----------------------------------------------------------
	//    🔌 SIMPLE & CLEAN SOCKET.IO SERVER (like Node.js)
	// ----------------------------------------------------------- */
	// server := socketio.NewServer(&engineio.Options{
	// 	Transports: []transport.Transport{
	// 		&websocket.Transport{
	// 			CheckOrigin: func(r *http.Request) bool { return true },
	// 		},
	// 	},
	// })

	// a.SocketServer = server

	// // ---- Event Handlers ----
	// server.OnConnect("/", func(conn socketio.Conn) error {
	// 	log.Println("🔌 Socket connected:", conn.ID())
	// 	conn.Emit("connected", "You are connected!")
	// 	return nil
	// })

	// server.OnEvent("/", "ping", func(conn socketio.Conn, msg string) {
	// 	log.Println("📩 PING received:", msg)
	// 	conn.Emit("pong", "PONG from server!")
	// })

	// server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
	// 	log.Println("❌ Socket disconnected:", conn.ID(), "reason:", reason)
	// })

	// server.OnError("/", func(conn socketio.Conn, err error) {
	// 	log.Println("🔥 Socket error:", err)
	// })

	// // Mount socket server into Gin
	// r.GET("/socket.io/*any", gin.WrapH(server))
	// r.POST("/socket.io/*any", gin.WrapH(server))

	// go server.Serve()
	// defer server.Close()

	// /* ----------------------------------------------------------- */

	// START APPLICATION SERVER
	port := "8000"
	if a != nil && a.Config != nil && a.Config.Port != "" {
		port = a.Config.Port
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("🚀 Findix server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}
