package main

import (
	"log"
	"regexp"
	"time"
	// "strconv"
	// "fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

// const notify_interval time.Duration = 1

func register(user_id string, c *websocket.Conn) {
	channel := getChannel(user_id)
	clients[channel] = append(clients[channel], c)
	pubsub.Subscribe(channel)
}

func notify() {
	for {
		switch v := pubsub.Receive().(type) {
		case redis.Message:
			// fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			writeMessage(v.Channel, []byte(string(v.Data)+": "+time.Now().String()))
		case redis.Subscription:
			// fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Fatal(v)
		}
	}
}

func writeMessage(channel string, message []byte) {
	conns := clients[channel]
	if len(conns) == 0 {
		pubsub.Unsubscribe(channel)
		user_id := extractUserId(channel)
		log.Println("Cannot find any websocket connections by user_id \"" + user_id + "\"")
		return
	}

	new_conns := make([]*websocket.Conn, len(conns))
	copy(new_conns, conns)

	for i, c := range conns {
		err := c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			c.Close()
			// 删除断开的连接
			new_conns = append(new_conns[:i], new_conns[i+1:]...)
		}
	}

	// unsubscribe channel
	if len(new_conns) == 0 {
		pubsub.Unsubscribe(channel)
	}

	if len(conns) != len(new_conns) {
		clients[channel] = new_conns
		// 清除引用，让GC回收已断开的连接，避免内存泄漏
		for i := 0; i < len(conns); i++ {
			conns[i] = nil
		}
		// fmt.Println(conns)
	}
	// fmt.Println(clients)
}

func extractUserId(s string) string {
	r, _ := regexp.Compile(`notification\.u_id_(\d+)\z`)
	match := r.FindStringSubmatch(s)
	if len(match) > 0 {
		return match[1]
	} else {
		return ""
	}
}

func getChannel(user_id string) string {
	// "{env}.notification.u_id_{user_id}"
	return "dev.notification.u_id_" + user_id
}
