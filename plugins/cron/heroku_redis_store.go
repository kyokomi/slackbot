package cron

import (
	"bytes"
	"log"

	"golang.org/x/net/context"
	"gopkg.in/redis.v2"
	"gopkg.in/vmihailenco/msgpack.v1"

	"github.com/kyokomi/goroku"
	"fmt"
	"strings"
)

const (
	redisCronTaskKey = "heroku:cron:task:%s"
)

type HerokuRedisStore struct {
	redisDB *redis.Client
}

func (s HerokuRedisStore) Close() {
	if s.redisDB != nil {
		s.redisDB.Close()
	}
}

func NewHerokuRedisStore() HerokuRedisStore {
	s := HerokuRedisStore{}
	s.redisDB = goroku.MustRedis(goroku.OpenRedis(context.Background()))
	return s
}

func (s HerokuRedisStore) Load() (map[string]CronTaskMap, error) {
	keys, err := s.redisDB.Keys(fmt.Sprintf(redisCronTaskKey, "*")).Result()
	if err != nil {
		return nil, err
	}

	result := map[string]CronTaskMap{}
	for _, redisKey := range keys {
		data, err := s.redisDB.HGetAllMap(redisKey).Result()
		if err != nil {
			return nil, err
		}

		taskMap := CronTaskMap{}
		for k, v := range data {
			task := CronTask{}
			if err := msgpack.Unmarshal(bytes.NewBufferString(v).Bytes(), &task); err != nil {
				log.Println(err)
				continue
			}

			if !task.Active {
				continue
			}

			taskMap[k] = task
		}

		channelID := redisKey[strings.LastIndex(redisKey, ":")+1:]
		result[channelID] = taskMap
	}

	return result, err
}

func (s HerokuRedisStore) Save(taskMap map[string]CronTaskMap) error {
	var err error
	for channelID, data := range taskMap {
		for key, val := range data {
			redisKey := fmt.Sprintf(redisCronTaskKey, channelID)

			if !val.Active {
				err = s.redisDB.HDel(redisKey, key).Err()
				continue
			}

			d, err := msgpack.Marshal(val)
			if err != nil {
				log.Println(err)
				continue
			}

			err = s.redisDB.HSet(redisKey, key, string(d)).Err()
			if err != nil {
				break
			}
		}
	}
	return err
}

var _ CronStore = (*HerokuRedisStore)(nil)
