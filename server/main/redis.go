package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 定義一個全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) { // 因為希望初始化的動作一起放在main裡, 所以不寫init()

	pool = &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空閒鏈接數
		MaxActive:   maxActive,   //表示和數據庫的最大鏈接數, 0 表示沒有限制
		IdleTimeout: idleTimeout, // 最大空閒時間
		Dial: func() (redis.Conn, error) { // 初始化鏈接的代碼, 鏈接哪個ip
			return redis.Dial("tcp", address)
		},
	}
}
