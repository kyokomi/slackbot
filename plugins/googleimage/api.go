package googleimage

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/bitly/go-simplejson"
)

type GoogleImageAPIClient interface {
	GetImageLinks(query string) ([]string, error)
}

type googleImageAPIClient struct {
	client *http.Client
	cx     string
	apiKey string
}

func NewGoogleImageAPIClient(httpClient *http.Client, cx string, apiKey string) GoogleImageAPIClient {
	return &googleImageAPIClient{client: httpClient, cx: cx, apiKey: apiKey}
}

func (g googleImageAPIClient) GetImageLinks(query string) ([]string, error) {
	params := url.Values{}
	params.Set("searchType", "image")
	params.Set("alt", "json")
	params.Set("cx", g.cx)
	params.Set("key", g.apiKey)
	params.Set("q", query)

	resp, err := g.client.Get(fmt.Sprintf("%s?%s", endpointURL, params.Encode()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("get request error statusCode %s \n%s", resp.Status, string(data))
	}

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	items, err := j.Get("items").Array()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var links []string
	for _, item := range items {
		link := item.(map[string]interface{})["link"].(string)
		links = append(links, link)
	}

	return links, nil
}
