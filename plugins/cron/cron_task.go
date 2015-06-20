package cron

import (
	"log"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
)

type CronTask struct {
	Active  bool
	Command CronCommand
}

type CronTaskMap map[string]CronTask

func (c CronTaskMap) AddTask(key string, task CronTask) {
	c[key] = task
}

type CronStore interface {
	Load() (map[string]CronTaskMap, error)
	Save(cronTaskMap map[string]CronTaskMap) error
}

var (
	cronClient  map[string]*cron.Cron
	cronStore   CronStore
	cronTaskMap = map[string]CronTaskMap{}
)

func Setup() {
	cronClient = map[string]*cron.Cron{}
}

func SetupStore(store CronStore) {
	cronStore = store
	data, err := cronStore.Load()
	if err != nil {
		log.Println(err)
	} else {
		cronTaskMap = data
	}
}

func RefreshCron(ctx context.Context) {
	for channelID, _ := range cronTaskMap {
		log.Println("RefreshCron channelID", channelID)
		refreshCron(ctx, channelID)
	}
}

func Stop() {
	if cronClient != nil {
		for _, c := range cronClient {
			if c == nil {
				continue
			}
			c.Stop()
		}
	}
}

func refreshCron(ctx context.Context, channelID string) {
	if cronClient[channelID] != nil {
		cronClient[channelID].Stop()
		cronClient[channelID] = nil
	}

	c := cron.New()
	for _, activeCron := range getTaskMap(channelID) {
		if !activeCron.Active {
			continue
		}

		message := activeCron.Command.Message
		c.AddFunc(activeCron.Command.CronSpec, func() {
			plugins.SendChannelMessage(ctx, channelID, message)
		})
	}
	c.Start()

	cronClient[channelID] = c

	if cronStore != nil {
		cronStore.Save(cronTaskMap)
	}
}

func startTask(channelID string, c CronCommand) {
	if cronTaskMap[channelID] == nil {
		cronTaskMap[channelID] = CronTaskMap{}
	}
	cronTaskMap[channelID].AddTask(c.CronKey(), CronTask{true, c})
}

func stopTask(channelID string, c CronCommand) {
	if cronTaskMap[channelID] == nil {
		cronTaskMap[channelID] = CronTaskMap{}
	}
	cronTaskMap[channelID].AddTask(c.CronKey(), CronTask{false, c})
}

func getTaskMap(channelID string) map[string]CronTask {
	if cronTaskMap[channelID] == nil {
		cronTaskMap[channelID] = CronTaskMap{}
	}
	return cronTaskMap[channelID]
}
