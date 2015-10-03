package kohaimage_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/kohaimage"
)

var testEvent = plugins.NewTestEvent("koha")

func TestCheckMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := kohaimage.NewPlugin(nil)

	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPI := kohaimage.NewMockKohaAPI(ctrl)
	mockAPI.EXPECT().GetImageURL().Return("http://pbs.twimg.com/media/CQJLPe7UAAAH20f.png")
	p := kohaimage.NewPlugin(mockAPI)

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != true")
	}
}
