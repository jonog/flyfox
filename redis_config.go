package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func RedisInit() {
	pool = NewPool("flyfox_redis:6379")
}

func NewPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
