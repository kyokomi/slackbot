package docomo

import (
	"strings"

	"fmt"

	"github.com/kyokomi/go-docomo/docomo"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
)

const repositoryDocomoContextKey = "docomo.context"

type plugin struct {
	client     *docomo.Client
	repository slackbot.Repository
}

func NewPlugin(docomoClient *docomo.Client, repository slackbot.Repository) plugins.BotMessagePlugin {
	return &plugin{
		client:     docomoClient,
		repository: repository,
	}
}

func (r plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	cmdArgs, ok := event.BotCmdArgs(message)
	if !ok {
		return false, message
	}
	return true, strings.Join(cmdArgs, " ")
}

func (p *plugin) Help() string {
	return `docomo: ドコモ雑談
	ドコモ雑談APIで会話します。部屋ごとにcontextを保持しています。
`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)

func (r plugin) DoAction(event plugins.BotEvent, message string) bool {
	// Contextを取得
	redisKey := buildRoomContextKey(event.Channel())
	dialogueCtx, _ := r.repository.Load(redisKey)

	// 雑談API叩く
	name := event.SenderName()
	req := docomo.DialogueRequest{
		Utt:      &message,
		Context:  &dialogueCtx,
		Nickname: &name,
	}
	res, err := r.client.Dialogue.Get(req, false)
	if err != nil {
		event.Reply("はて？")
		return true
	} else if res.RequestError.PolicyException.MessageID != "" {
		event.Reply(fmt.Sprintf("%s: %s",
			res.RequestError.PolicyException.MessageID,
			res.RequestError.PolicyException.Text),
		)
		return true
	}

	// 結果を返す
	event.Reply(res.Utt)

	// Contextを保存
	if dialogueCtx == "" {
		r.repository.Save(redisKey, res.Context)
	}

	return false // next ok
}

func buildRoomContextKey(room string) string {
	return fmt.Sprintf("%s.%s", repositoryDocomoContextKey, room)
}
