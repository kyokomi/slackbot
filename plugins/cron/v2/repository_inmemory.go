package cron

// NewInMemoryRepository create in memory repository
func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		taskMap: map[string]TaskMap{},
	}
}

type inMemoryRepository struct {
	taskMap map[string]TaskMap
}

func (r *inMemoryRepository) Load() (map[string]TaskMap, error) {
	return r.taskMap, nil
}

func (r *inMemoryRepository) Save(cronTaskMap map[string]TaskMap) error {
	r.taskMap = cronTaskMap
	return nil
}

func (r *inMemoryRepository) Close() error {
	return nil
}
