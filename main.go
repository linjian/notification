package main

import (
	"log"
	"net/http"
	// "fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients = make(map[string][]*websocket.Conn)
	pubsub  redis.PubSubConn
)

func main() {
	pubsub = redisDial()
	defer pubsub.Close()

	go notify()

	http.HandleFunc("/", handleConnections)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	register(r.FormValue("user_id"), c)
}
