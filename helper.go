package httphelper

import (
	"github.com/artisancloud/httphelper/client"
	"github.com/artisancloud/httphelper/dataflow"
	"github.com/artisancloud/httphelper/driver/nethttp"
)

type Helper interface {
	SetClient(client client.Client)
	GetClient() client.Client

	WithMiddleware(middlewares ...dataflow.RequestMiddleware)

	Df() dataflow.RequestDataflow
}

type RequestHelper struct {
	client           client.Client
	middlewareHandle dataflow.RequestMiddleware
	config           *Config
}

type Config struct {
	*client.Config
	BaseUrl string
}

func NewRequestHelper(conf *Config) (Helper, error) {
	c, err := nethttp.NewHttpClient(conf.Config)
	if err != nil {
		return nil, err
	}
	return &RequestHelper{
		client: c,
		middlewareHandle: func(handle dataflow.RequestHandle) dataflow.RequestHandle {
			return handle
		},
		config: conf,
	}, nil
}

func (r *RequestHelper) SetClient(client client.Client) {
	r.client = client
}

func (r *RequestHelper) GetClient() client.Client {
	return r.client
}

func (r *RequestHelper) WithMiddleware(middlewares ...dataflow.RequestMiddleware) {
	if len(middlewares) == 0 {
		return
	}
	var buildHandle func(md dataflow.RequestMiddleware, appendMd dataflow.RequestMiddleware) dataflow.RequestMiddleware
	buildHandle = func(md dataflow.RequestMiddleware, appendMd dataflow.RequestMiddleware) dataflow.RequestMiddleware {
		return func(handle dataflow.RequestHandle) dataflow.RequestHandle {
			return md(appendMd(handle))
		}
	}
	for _, middleware := range middlewares {
		r.middlewareHandle = buildHandle(r.middlewareHandle, middleware)
	}
}

func (r *RequestHelper) Df() dataflow.RequestDataflow {
	return dataflow.NewDataflow(r.client, r.middlewareHandle, &dataflow.Option{
		BaseUrl: r.config.BaseUrl,
	})
}
