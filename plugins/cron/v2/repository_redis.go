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
	defaultRedisCronTaskKey = "kyokomi/slackbot:cron:task:%s"
)

type redisRepository struct {
	redisKeyFmt string
	redisDB     *redis.Client
}

func (s redisRepository) Close() error {
	if s.redisDB != nil {
		return s.redisDB.Close()
	}
	return nil
}

// NewRedisRepository create redis repository
func NewRedisRepository(addr, password string, db int64, redisKeyFmt string) Repository {
	if redisKeyFmt == "" {
		redisKeyFmt = defaultRedisCronTaskKey
	}
	return &redisRepository{
		redisKeyFmt: redisKeyFmt,
		redisDB: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (s *redisRepository) Load() (map[string]TaskMap, error) {
	keys, err := s.redisDB.Keys(fmt.Sprintf(s.redisKeyFmt, "*")).Result()
	if err != nil {
		return nil, err
	}

	result := map[string]TaskMap{}
	for _, redisKey := range keys {
		data, err := s.redisDB.HGetAllMap(redisKey).Result()
		if err != nil {
			return nil, err
		}

		taskMap := TaskMap{}
		for k, v := range data {
			task := Task{}
			if err := msgpack.Unmarshal(bytes.NewBufferString(v).Bytes(), &task); err != nil {
				log.Println(err) // TODO: error log
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

func (s *redisRepository) Save(taskMap map[string]TaskMap) error {
	var err error
	for channelID, data := range taskMap {
		for key, val := range data {
			redisKey := fmt.Sprintf(s.redisKeyFmt, channelID)

			if !val.Active {
				err = s.redisDB.HDel(redisKey, key).Err()
				continue
			}

			d, err := msgpack.Marshal(val)
			if err != nil {
				log.Println(err) // TODO: error log
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

var _ Repository = (*redisRepository)(nil)
