package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

func redisDial() redis.PubSubConn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	return redis.PubSubConn{Conn: c}
}
