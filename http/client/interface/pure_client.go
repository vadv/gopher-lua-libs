package http_interface

import (
	"net/http"
	"net/url"
)

func NewPureClient() *pureClient {
	return &pureClient{Client: &http.Client{}}
}

type pureClient struct {
	*http.Client
}

func (c *pureClient) DoRequest(req *http.Request) (*http.Response, error) {
	return c.Do(req)
}

func (c *pureClient) PostFormRequest(url string, data url.Values) (*http.Response, error) {
	return c.PostForm(url, data)
}
