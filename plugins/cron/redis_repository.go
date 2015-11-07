package cron

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"gopkg.in/redis.v3"
	"gopkg.in/vmihailenco/msgpack.v1"
)

const (
	redisCronTaskKey = "heroku:cron:task:%s"
)

type RedisRepository struct {
	redisDB *redis.Client
}

func (s RedisRepository) Close() error {
	if s.redisDB != nil {
		return s.redisDB.Close()
	}
	return nil
}

func NewRedisRepository(addr, password string, db int64) CronRepository {
	s := &RedisRepository{}
	s.redisDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return s
}

func (s RedisRepository) Load() (map[string]CronTaskMap, error) {
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

			log.Printf("%#v\n", task)

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

func (s RedisRepository) Save(taskMap map[string]CronTaskMap) error {
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

var _ CronRepository = (*RedisRepository)(nil)
