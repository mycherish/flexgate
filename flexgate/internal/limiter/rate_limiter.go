package limiter

import (
	"sync"
	"time"
)

// TokenBucket 令牌桶结构体
type TokenBucket struct {
	rate     int64      // 每秒生成的令牌数 (每秒充能多少)
	capacity int64      // 桶的最大容量 (允许的瞬间并发峰值)
	tokens   int64      // 当前桶内剩余的令牌数
	lastTick time.Time  // 上次放令牌的时间点
	mu       sync.Mutex // 保证并发安全，防止令牌超卖
}

// NewTokenBucket 创建一个新的限流器实例
func NewTokenBucket(rate int64, capacity int64) *TokenBucket {
	return &TokenBucket{
		rate:     rate,
		capacity: capacity,
		tokens:   capacity,   // 初始状态桶是满的
		lastTick: time.Now(), // 记录初始化时间
	}
}

// Allow 尝试获取 1 个令牌
// 返回 true 表示允许通过，false 表示已被限流
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	// 1. 计算自上次获取令牌以来，这段时间应该补充多少令牌
	// 补充量 = 时间间隔(秒) * 每秒速率
	duration := now.Sub(b.lastTick)
	newTokens := int64(duration.Seconds()) * b.rate

	if newTokens > 0 {
		b.tokens += newTokens
		// 确保令牌数不会超过桶的上限
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
		// 只有真正增加了令牌，才更新时间戳，避免微小时间碎片导致计算损失
		b.lastTick = now
	}

	// 2. 判断桶里是否还有令牌
	if b.tokens > 0 {
		b.tokens-- // 消耗一个令牌
		return true
	}

	return false
}
