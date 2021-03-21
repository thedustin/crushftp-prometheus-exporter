package crushftp

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    url.URL
	httpClient http.Client
	logger     *log.Logger
}

type ClientOptions struct {
	HostAndPort string
	Http        HttpClientOptions
	Logger      *log.Logger
	Password    string
	PathBase    string
	Scheme      string
	Username    string
}

type HttpClientOptions struct {
	Insecure bool
}

const DefaultScheme = "https"

func NewClient(opts ClientOptions) *Client {
	if opts.Scheme == "" {
		opts.Scheme = DefaultScheme
	}

	opts.PathBase = strings.TrimSuffix(opts.PathBase, "/")
	opts.PathBase = opts.PathBase + "/"

	if opts.Logger == nil {
		opts.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return &Client{
		baseURL: url.URL{
			Scheme: opts.Scheme,
			Host:   opts.HostAndPort,
			Path:   opts.PathBase,
			User:   url.UserPassword(opts.Username, opts.Password),
		},
		httpClient: createHttpClient(opts.Http.Insecure),
		logger:     opts.Logger,
	}
}

func createHttpClient(insecure bool) http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.InsecureSkipVerify = insecure

	return http.Client{
		Transport: transport,
	}
}

func (c *Client) request(path string, params map[string]string) (*http.Response, error) {
	url := c.baseURL
	url.Path = url.Path + path

	q := url.Query()
	for param, value := range params {
		q.Add(param, value)
	}
	url.RawQuery = q.Encode()

	c.logger.Printf("%s %s", http.MethodGet, stripPassword(&url))

	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	return c.httpClient.Do(req)
}

func (c *Client) command(command string) (*http.Response, error) {
	return c.request("WebInterface/function/", map[string]string{"command": command})
}

func stripPassword(u *url.URL) string {
	_, passSet := u.User.Password()
	if passSet {
		return strings.Replace(u.String(), u.User.String()+"@", u.User.Username()+":***@", 1)
	}
	return u.String()
}
