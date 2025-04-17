package main

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis   *redis.Client
	redsync *redsync.Redsync
}

func NewCache() *Cache {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "xxxxx",          // Redis 密码
		DB:       0,                // 使用默认的数据库
	})

	// 检查连接是否正常
	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}

	// 初始化 Redsync
	rs := redsync.New(goredis.NewPool(client))

	return &Cache{
		redis:   client,
		redsync: rs,
	}
}

// 实现一个分布式锁
func (c *Cache) LockResource(lockName string) (*redsync.Mutex, error) {
	// 创建分布式锁
	mutex := c.redsync.NewMutex(lockName)

	// 尝试加锁，最多等待 10 秒，锁的过期时间为 30 秒
	if err := mutex.LockContext(context.Background()); err != nil {
		return nil, fmt.Errorf("无法获得锁: %v", err)
	}

	// 如果成功获得锁，返回成功
	return mutex, nil
}

// 释放分布式锁
func (c *Cache) UnlockResource(mutex *redsync.Mutex) error {
	if _, err := mutex.Unlock(); err != nil {
		return fmt.Errorf("无法释放锁: %v", err)
	}
	return nil
}
