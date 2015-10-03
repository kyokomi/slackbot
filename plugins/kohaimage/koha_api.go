package kohaimage

import (
	"io/ioutil"
	"net/http"
)

//go:generate mockgen -package kohaimage -source koha_api.go -destination koha_api_mock.go

type KohaAPI interface {
	GetImageURL() string
}

type kohaAPI struct {
	client *http.Client
}

func (k *kohaAPI) GetImageURL() string {
	resp, err := k.client.Get("https://koha-api.appspot.com/v1/api/image")
	if err != nil {
		return "image not found"
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "image not found"
	}
	return string(data)
}

func NewKohaAPI() KohaAPI {
	koha := &kohaAPI{
		client: &http.Client{},
	}
	return koha
}
