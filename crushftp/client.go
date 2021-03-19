package crushftp

import (
	"net/http"
	"net/url"
)

type client struct {
	baseURL url.URL
	client  http.Client
}

func NewClient(protocol, hostAndPort, username, password string) *client {
	return &client{
		baseURL: url.URL{
			Scheme: protocol,
			Host:   hostAndPort,
			User:   url.UserPassword(username, password),
		},
		client: http.Client{},
	}
}

func (c *client) request(path string, params map[string]string) (*http.Response, error) {
	url := c.BaseUrl()
	url.Path = path

	q := url.Query()
	for param, value := range params {
		q.Add(param, value)
	}
	url.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	return c.client.Do(req)
}

func (c *client) command(command string) (*http.Response, error) {
	return c.request("/WebInterface/function/", map[string]string{"command": command})
}

func (c *client) BaseUrl() url.URL {
	return c.baseURL
}
