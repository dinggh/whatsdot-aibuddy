package httpapi

import (
	"sync"
	"time"
)

type tokenBucket struct {
	Tokens     float64
	Capacity   float64
	RefillRate float64
	LastRefill time.Time
}

type DeviceLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	capacity float64
	refill   float64
}

func NewDeviceLimiter(capacity int, refillPerMinute int) *DeviceLimiter {
	if capacity <= 0 {
		capacity = 1
	}
	if refillPerMinute <= 0 {
		refillPerMinute = 1
	}
	return &DeviceLimiter{
		buckets:  make(map[string]*tokenBucket),
		capacity: float64(capacity),
		refill:   float64(refillPerMinute) / 60.0,
	}
}

func (d *DeviceLimiter) Allow(deviceID string) bool {
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()

	b, ok := d.buckets[deviceID]
	if !ok {
		d.buckets[deviceID] = &tokenBucket{
			Tokens:     d.capacity - 1,
			Capacity:   d.capacity,
			RefillRate: d.refill,
			LastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(b.LastRefill).Seconds()
	if elapsed > 0 {
		b.Tokens += elapsed * b.RefillRate
		if b.Tokens > b.Capacity {
			b.Tokens = b.Capacity
		}
		b.LastRefill = now
	}

	if b.Tokens < 1 {
		return false
	}
	b.Tokens--
	return true
}
