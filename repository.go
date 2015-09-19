package slackbot

type Repository interface {
	Load(key string) (string, error)
	Save(key string, value string) error
	Close() error
}

func NewOnMemoryRepository() Repository {
	return &OnMemoryRepository{
		memory: map[string]string{},
	}
}

type OnMemoryRepository struct {
	memory map[string]string
}

func (r *OnMemoryRepository) Load(key string) (string, error) {
	return r.memory[key], nil
}

func (r *OnMemoryRepository) Save(key string, value string) error {
	r.memory[key] = value
	return nil
}

func (r *OnMemoryRepository) Close() error {
	return nil
}
