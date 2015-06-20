package cron

import (
	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
)

type cronTask struct {
	active  bool
	command cronCommand
}

var (
	cronClient  *cron.Cron
	cronTaskMap = map[string]cronTask{}
)

func Setup() {
	cronClient = nil
	// TODO: Redisとかに保存してるやつをcronTaskMapにセットしてrefreshCronする
}

func Stop() {
	if cronClient != nil {
		cronClient.Stop()
	}
}

func refreshCron(ctx context.Context) {
	if cronClient != nil {
		cronClient.Stop()
	}

	c := cron.New()
	for _, activeCron := range cronTaskMap {
		if !activeCron.active {
			continue
		}

		message := activeCron.command.Message
		c.AddFunc(activeCron.command.CronSpec, func() {
			plugins.SendMessage(ctx, message)
		})
	}
	c.Start()

	cronClient = c
}

func startTask(c cronCommand) {
	cronTaskMap[c.CronKey()] = cronTask{true, c}
}

func stopTask(c cronCommand) {
	cronTaskMap[c.CronKey()] = cronTask{false, c}
}

func getTaskMap() map[string]cronTask {
	return cronTaskMap
}
