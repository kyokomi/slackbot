package cron

type Task struct {
	Active  bool
	Command Command
}

type TaskMap map[string]Task

func (c TaskMap) AddTask(key string, task Task) {
	c[key] = task
}
