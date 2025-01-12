package cache

import "github.com/redis/go-redis/v9"

type BruteForceChecker interface{}

type RedisBruteForceChecker struct {
	cmd redis.Cmdable
}

func NewBruteForceChecker(cmd redis.Cmdable) BruteForceChecker {
	return &RedisBruteForceChecker{}
}
