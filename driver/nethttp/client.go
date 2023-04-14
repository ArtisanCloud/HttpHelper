package nethttp

import (
	"crypto/tls"
	"github.com/artisancloud/httphelper/client"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

// Client 是 net/http 的封装
type Client struct {
	conf       *client.Config
	coreClient *http.Client
}

func NewHttpClient(config *client.Config) (*Client, error) {
	if config == nil {
		config = &client.Config{}
		config.Default()
	}

	coreClient, err := newCoreClient(*config)
	if err != nil {
		return nil, err
	}

	return &Client{
		conf:       config,
		coreClient: coreClient,
	}, nil
}

// SetConfig 配置客户端
func (c *Client) SetConfig(config client.Config) error {
	config.Default()
	c.conf = &config

	coreClient, err := newCoreClient(config)
	if err != nil {
		return errors.Wrap(err, "failed to create core client use new config")
	}

	c.coreClient = coreClient
	return nil
}

// GetConfig 返回配置副本
func (c *Client) GetConfig() client.Config {
	return *c.conf
}

func (c *Client) DoRequest(request *http.Request) (response *http.Response, err error) {
	return c.coreClient.Do(request)
}

func newCoreClient(config client.Config) (*http.Client, error) {
	coreClient := http.Client{
		Timeout: config.Timeout,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()

	if config.Cert.CertFile != "" && config.Cert.KeyFile != "" {
		certPair, err := tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load certificate")
		}
		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{certPair},
		}
	}

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse proxy URL")
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &coreClient, nil
}
