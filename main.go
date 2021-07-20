package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
)

var users = []string{}
var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "3000"
	}

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println(s.ID())
		return nil
	})
	server.OnEvent("/", "list-include-me", func(s socketio.Conn, username string) {
		users = append(users, username)
		s.Emit("list-update", users)
	})
	server.OnEvent("/", "connected", func(s socketio.Conn, username string) {
		server.BroadcastToNamespace("/", "connected-user", username)
	})
	server.OnEvent("/", "msg-in", func(s socketio.Conn, msg string, username string) {
		var cMsg string
		fmt.Sprintf(cMsg, "%s: %s", username, msg)
		server.BroadcastToNamespace("/", "msg-out", cMsg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket listen error: %v", err)
		}
	}()
	defer server.Close()

	http.Handle("/", server)
	
	fmt.Println("Running on port " + port)
	http.ListenAndServe(":" + port, nil)
}