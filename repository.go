package slackbot

import (
	"time"

	"gopkg.in/redis.v3"
)

// Repository is slackbot common repository interface
type Repository interface {
	Load(key string) (string, error)
	Save(key string, value string) error
	Close() error
}

// NewOnMemoryRepository create on memory repository
func NewOnMemoryRepository() Repository {
	return &inMemoryRepository{
		memory: map[string]string{},
	}
}

type inMemoryRepository struct {
	memory map[string]string
}

func (r *inMemoryRepository) Load(key string) (string, error) {
	return r.memory[key], nil
}

func (r *inMemoryRepository) Save(key string, value string) error {
	r.memory[key] = value
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

func (r *redisRepository) Load(key string) (string, error) {
	return r.redisDB.Get(key).Result()
}

func (r *redisRepository) Save(key string, value string) error {
	return r.redisDB.Set(key, value, time.Duration(0)).Err()
}

func (r *redisRepository) Close() error {
	return r.redisDB.Close()
}
