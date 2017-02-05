package slackbot

import (
	"time"

	"gopkg.in/redis.v3"
)

// Repository is slackbot common repository interface
type Repository interface {
	Load(key string) (string, error)
	LoadList(key string) ([]string, error)
	Save(key string, value string) error
	SaveList(key string, values []string) error
	Close() error
}

// NewOnMemoryRepository create on memory repository
func NewOnMemoryRepository() Repository {
	return &inMemoryRepository{
		memory: map[string]interface{}{},
	}
}

type inMemoryRepository struct {
	memory map[string]interface{}
}

func (r *inMemoryRepository) LoadList(key string) ([]string, error) {
	return r.memory[key].([]string), nil
}

func (r *inMemoryRepository) Load(key string) (string, error) {
	return r.memory[key].(string), nil
}

func (r *inMemoryRepository) Save(key string, value string) error {
	r.memory[key] = value
	return nil
}

func (r *inMemoryRepository) SaveList(key string, values []string) error {
	r.memory[key] = values
	return nil
}

func (r *inMemoryRepository) Close() error {
	return nil
}

// NewRedisRepository create on redis repository
func NewRedisRepository(addr, password string, db int64) Repository {
	repo := &redisRepository{}
	repo.redisDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return repo
}

type redisRepository struct {
	redisDB *redis.Client
}

func (r *redisRepository) LoadList(key string) ([]string, error) {
	return r.redisDB.LRange(key, 0, -1).Result()
}

func (r *redisRepository) Load(key string) (string, error) {
	cmd := r.redisDB.Get(key)
	if cmd.Err() == redis.Nil {
		return "", nil
	}
	return cmd.Result()
}

func (r *redisRepository) Save(key string, value string) error {
	return r.redisDB.Set(key, value, time.Duration(0)).Err()
}

func (r *redisRepository) SaveList(key string, values []string) error {
	return r.redisDB.LPush(key, values...).Err()
}

func (r *redisRepository) Close() error {
	return r.redisDB.Close()
}
