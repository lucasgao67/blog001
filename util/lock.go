package util

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

var RedisClient *redis.Client

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: ":32771",
	})
}

func LockBlock(key string) error {

	// 不过期
	for {
		if cmd := RedisClient.SetNX(key, time.Now().Unix(), 0); cmd.Val() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("")

}

func UnLodk(key string) {
	RedisClient.Del(key)
}
