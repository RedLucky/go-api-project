package cache

import (
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

type RedisConn struct {
	Pool *redis.Pool
}

func New(redisAddress string) *RedisConn {
	return &RedisConn{
		Pool: &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", redisAddress)
				if err != nil {
					log.Printf("ERROR: fail init redis pool: %s", err.Error())
					os.Exit(1)
				}
				return conn, err
			},
		},
	}
}
