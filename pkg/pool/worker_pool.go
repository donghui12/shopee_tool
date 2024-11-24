package pool

import (
    "sync"

    "github.com/panjf2000/ants"
)

var (
	pool *ants.Pool
	once sync.Once
)

// 初始化全局 WorkerPool
func InitPool(size int) {
	once.Do(func() {
		var err error
		pool, err = ants.NewPool(size)
		if err != nil {
			panic(err)
		}
	})
}

// 获取全局 WorkerPool
func GetPool() *ants.Pool {
    return pool
}