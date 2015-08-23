package cron

func NewOnMemoryCronRepository() CronRepository {
	return OnMemoryCronRepository{
		taskMap: map[string]CronTaskMap{},
	}
}

type OnMemoryCronRepository struct {
	taskMap map[string]CronTaskMap
}

func (r OnMemoryCronRepository) Load() (map[string]CronTaskMap, error) {
	return r.taskMap, nil
}

func (r OnMemoryCronRepository) Save(cronTaskMap map[string]CronTaskMap) error {
	r.taskMap = cronTaskMap
	return nil
}

func (r OnMemoryCronRepository) Close() error {
	return nil
}
