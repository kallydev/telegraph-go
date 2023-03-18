package telegraph

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL = "https://telegra.ph/%s"
	apiURL  = "https://api.telegra.ph/%s"
)

type ClientOption struct {
	Proxy   string
	Timeout time.Duration
}

type Client struct {
	AccessToken string
	option      *ClientOption
	httpClient  *http.Client
}

func NewClient(accessToken string, option *ClientOption) (client *Client, err error) {
	client = &Client{
		AccessToken: accessToken,
		option:      option,
		httpClient:  new(http.Client),
	}
	if client.option != nil {
		if client.option.Timeout != 0 {
			client.httpClient.Timeout = client.option.Timeout
		}
		if len(client.option.Proxy) > 0 {
			proxyURL, err := url.Parse(client.option.Proxy)
			if err != nil {
				return nil, err
			}
			client.httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}
	return client, nil
}

func (client *Client) post(method string, parm url.Values) (response []byte, err error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf(apiURL, method), strings.NewReader(parm.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
