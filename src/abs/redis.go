package abs

import (
	"fmt"
	"github.com/go-redis/redis"
	app "src/config"
	"src/library/tool"
	"sync"
	"time"
)

const (
	retryConnRdsMaxTimes = 3
)

type Redis struct {
}

type rdsCollector struct {
	once    sync.Once
	pointer *redis.Client
}

var cache, ssd, session rdsCollector

func (r *Redis) newClient(client string) (*redis.Client, error) {
	dt, ok := app.Redis.Client[client]
	if !ok {
		return nil, fmt.Errorf("redisClient of %s is not exist", client)
	}
	rds := redis.NewClient(&redis.Options{
		Addr:     dt.Host + ":" + dt.Port,
		Password: dt.Pwd,
		DB:       dt.Db,
	})
	return rds, nil
}

func (r *Redis) getClient(rdc *rdsCollector, client string) *redis.Client {
	var (
		err   error
		fatal bool
	)
	for i := 1; i <= retryConnRdsMaxTimes; i++ {
		//log.Printf("getClient %s try --->: %d\n", client, i)
		rdc.once.Do(func() {
			//log.Printf("newClient try --->: %s\n", client)
			rdc.pointer, err = r.newClient(client)
		})
		if err != nil {
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				rdc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		if err = rdc.pointer.Ping().Err(); err != nil {
			if i < retryConnRdsMaxTimes {
				time.Sleep(time.Second * time.Duration(i))
				rdc.once = sync.Once{}
				continue
			} else {
				fatal = true
				break
			}
		}
		break
	}
	if fatal {
		defer func() {
			errMsg := fmt.Sprintf("get redis client of %s occur error: %s", client, err.Error())
			app.Logger.Error(errMsg)
			tool.PushSimpleMessage(errMsg, true)
		}()
		return nil
	}
	return rdc.pointer
}

func (r *Redis) GetCache() *redis.Client {
	return r.getClient(&cache, "cache")
}

func (r *Redis) GetSSD() *redis.Client {
	return r.getClient(&ssd, "ssd")
}

func (r *Redis) GetSession() *redis.Client {
	return r.getClient(&session, "session")
}
