package util

import (
	"github.com/go-redis/redis"
	"time"
)

var RedisClient *redis.Client

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: ":32771",
	})
}

func LockBlock(key string) {

	// 不过期
	for {
		if cmd := RedisClient.SetNX(key, time.Now().Unix(), 0); cmd.Val() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func UnLodk(key string) {
	RedisClient.Del(key)
}
