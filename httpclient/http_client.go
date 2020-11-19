package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"getui_gosdk/consts"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Method string
	Path   string
	Body   []byte
	Header []HTTPOption
}

type Response struct {
	Code           int             `json:"code"`
	Msg            string          `json:"msg"`
	Data           json.RawMessage `json:"data"`
	HttpStatusCode int             `json:"http_status_code"`
}

type RetryConfig struct {
	MaxRetryTimes int
	RetryInterval time.Duration
}

type HTTPOption func(r *http.Request)

func SetHeader(key string, value string) HTTPOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}

type HTTPClient struct {
	retryTimes int
	cli        *http.Client
}

func NewClient() *HTTPClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{Transport: tr}
	return &HTTPClient{
		retryTimes: 3,
		cli:        c,
	}
}

func (r *Request) buildHTTPRequest() (*http.Request, error) {
	var body io.Reader

	if r.Body != nil {
		body = bytes.NewBuffer(r.Body)
	}
	fmt.Println(consts.BaseUrl + r.Path)
	req, err := http.NewRequest(r.Method, consts.BaseUrl+r.Path, body)
	if err != nil {
		return nil, err
	}

	for _, opt := range r.Header {
		opt(req)
	}

	return req, nil
}

func (c *HTTPClient) do(req *Request) (*Response, error) {
	request, err := req.buildHTTPRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.cli.Do(request)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	res.HttpStatusCode = resp.StatusCode
	return &res, nil
}

func (c *HTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	var (
		result *Response
		err    error
	)
	for i := 0; i < c.retryTimes; i++ {
		result, err = c.do(req)
		if err == nil {
			break
		}
	}
	return result, err
}
