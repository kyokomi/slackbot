package plugins

import (
	"fmt"
	"strings"
)

type PluginManager interface {
	AddPlugin(key interface{}, val BotMessagePlugin)
	ExecPlugins(botEvent BotEvent)
	StopReply()
	StartReply()
	IsReply() bool
	GetPlugins() []Plugin // deepCopy
}

type plugins struct {
	plugins       []Plugin
	isReply       bool
	messageSender MessageSender
}

var _ PluginManager = (*plugins)(nil)

type MessageSender interface {
	SendMessage(message string, channel string)
}

func NewPluginManager(sender MessageSender) PluginManager {
	return &plugins{
		plugins:       []Plugin{},
		isReply:       true,
		messageSender: sender,
	}
}

func (ps *plugins) AddPlugin(key interface{}, val BotMessagePlugin) {
	ps.plugins = append(ps.plugins, Plugin{key, val})
}

func (ps *plugins) StopReply() {
	ps.isReply = false
}

func (ps *plugins) StartReply() {
	ps.isReply = true
}

func (ps *plugins) IsReply() bool {
	return ps.isReply
}

func (ps *plugins) ExecPlugins(botEvent BotEvent) {
	for _, p := range ps.plugins {
		ok, m := p.CheckMessage(botEvent, botEvent.text)
		if !ok {
			continue
		}

		next := p.DoAction(botEvent, m)
		if !next {
			break
		}
	}
}

func (ps *plugins) SendMessage(message string, channel string) {
	ps.messageSender.SendMessage(message, channel)
}

func (ps *plugins) GetPlugins() []Plugin {
	deepCopy := make([]Plugin, len(ps.plugins))
	copy(deepCopy, ps.plugins)
	return deepCopy
}

type BotMessagePlugin interface {
	CheckMessage(event BotEvent, message string) (bool, string)
	DoAction(event BotEvent, message string) bool
	Help() string
}

type Plugin struct {
	Key interface{}
	BotMessagePlugin
}

func (p Plugin) Name() string {
	return fmt.Sprintf("%s", p.Key)
}

type BotID string

func (b BotID) Equal(bot string) bool {
	if string(b) == bot {
		return true
	}

	if b.LinkID() == bot {
		return true
	}

	return false
}

func (b BotID) LinkID() string {
	return fmt.Sprintf("<@%s>", b)
}

type BotEvent struct {
	messageSender MessageSender

	domain string

	botID   BotID
	botName string

	senderID    string
	senderName  string
	text        string
	channelID   string
	channelName string
	timestamp   string
}

var _ MessageSender = (*BotEvent)(nil)

func NewBotEvent(sender MessageSender, domain, botID, botName, senderID, senderName, text, channelID, channelName, timestamp string) BotEvent {
	return BotEvent{
		messageSender: sender,
		domain:        domain,
		botID:         BotID(botID),
		botName:       botName,
		senderID:      senderID,
		senderName:    senderName,
		text:          text,
		channelID:     channelID,
		channelName:   channelName,
		timestamp:     timestamp,
	}
}

func (b *BotEvent) Reply(message string) {
	b.SendMessage(message, b.Channel())
}

func (b *BotEvent) SendMessage(message string, channel string) {
	b.messageSender.SendMessage(message, channel)
}

func (b *BotEvent) BaseText() string {
	return b.text
}

// Channel is Deprecated
func (b *BotEvent) Channel() string {
	return b.ChannelID()
}

func (b *BotEvent) ChannelID() string {
	return b.channelID
}

func (b *BotEvent) ChannelName() string {
	return b.channelName
}

func (b *BotEvent) BotID() string {
	return string(b.botID)
}

func (b *BotEvent) BotName() string {
	return string(b.botName)
}

func (b *BotEvent) BotLinkID() string {
	return b.botID.LinkID()
}
func (b *BotEvent) BotLinkIDForClient() string {
	return b.botID.LinkID() + ":"
}

func (b *BotEvent) SenderID() string {
	return b.senderID
}

func (b *BotEvent) SenderName() string {
	return b.senderName
}

// ArchivesURL return copy link archives url
func (b *BotEvent) ArchivesURL() string {
	// Slackの仕様が変わると使えなくなるのであまり推奨しない方法
	return fmt.Sprintf("https://%s.slack.com/archives/%s/p%s", b.domain, b.channelName, strings.Replace(b.timestamp, ".", "", 1))
}

func (b *BotEvent) BotCmdArgs(message string) ([]string, bool) {
	trimMessage, ok := b.BotCmdMessage(message)
	return strings.Fields(trimMessage), ok
}

func (b *BotEvent) BotCmdMessage(message string) (string, bool) {
	switch {
	case strings.HasPrefix(message, b.BotLinkIDForClient()):
		return message[len(b.BotLinkIDForClient()):], true
	case strings.HasPrefix(message, b.BotLinkID()):
		return message[len(b.BotLinkID()):], true
	case strings.HasPrefix(message, b.BotName()):
		return message[len(b.BotName()):], true
	case strings.HasPrefix(message, b.BotID()):
		return message[len(b.BotID()):], true
	default:
		return "", false
	}
}

func (b *BotEvent) GetMessageSender() MessageSender {
	return b.messageSender
}
