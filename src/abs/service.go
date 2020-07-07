package abs

import (
	app "src/config"
	"time"
)

type Service struct {
	Redis
}

// RedisLock 加锁
func (s *Service) RedisLock(lookName string, expire time.Duration) bool {
	if lookName == "" {
		return false
	}
	key := app.Redis.Key.Bbs.StringRedisLock + lookName
	Time := time.Now()
	nowTime := Time.Unix()
	expireTimeNow := Time.Add(expire).Unix()
	cache := s.GetCache()
	res := cache.SetNX(key, expireTimeNow, 0).Val()
	if res {
		cache.Expire(key, expire)
		return true
	}
	oldExpireTime, _ := cache.Get(key).Int64()
	if nowTime > oldExpireTime {
		cache.Set(key, expireTimeNow, 0)
		return true
	}
	return false
}

// RedisUnLock 解锁
func (s *Service) RedisUnLock(lookName string) bool {
	if lookName == "" {
		return false
	}
	key := app.Redis.Key.Bbs.StringRedisLock + lookName
	redisCache := s.GetCache()
	if redisCache.Exists(key).Val() > 0 {
		redisCache.Del(key)
	}
	return true
}
