package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cleverreach.com/crtools/crconfig"
)

// CatchRequest for catching the request yourself, e.g. for testing. Example:
// rest.CatchRequest = func(req *http.Request) *http.Response {
// 	w := httptest.NewRecorder()
// 	w.WriteString("Hello World")
// 	return w.Result()
// }
var CatchRequest func(*http.Request) *http.Response

// Client is a simple and small REST Client.
// Use NewClient() or NewClientFromURL() to get an instance ready to be used.
type Client struct {
	url               string
	params            map[string]string
	header            map[string]string
	timeout           time.Duration
	retries           int
	maxRetryTime      time.Duration
	ignoreCertificate bool
	token             struct {
		method string
		value  string
		isset  bool
	}
}

// NewClientFromURL creates a yarc.Client by parsing given url it into the internal structure
func NewClientFromURL(urlPath string) *Client {
	Logger.Debugln("NewClientFromUrl ", urlPath)
	pos := strings.Index(urlPath, "?")

	uri, query := urlPath, ""
	if pos > 0 {
		uri = urlPath[:pos]
		query = urlPath[pos+1:]
	}

	c := Client{
		url:    uri,
		params: make(map[string]string),
		header: make(map[string]string),

		timeout:           crconfig.GetDuration("REST_CLIENT_TIMEOUT", 0) * time.Millisecond,
		retries:           int(crconfig.GetInt("REST_CLIENT_RETRIES", 0)),
		maxRetryTime:      crconfig.GetDuration("REST_CLIENT_MAX_RETRY_TIME", 60) * time.Second,
		ignoreCertificate: crconfig.GetBool("REST_CLIENT_IGNORE_CERTIFICATE", IgnoreCertificate),
	}

	for _, kv := range strings.Split(query, "&") {
		pos := strings.Index(kv, "=")
		if pos >= 0 {
			val := kv[pos+1:]
			if v, e := url.QueryUnescape(val); e == nil {
				val = v
			}
			c.params[kv[:pos]] = val
		} else if kv != "" {
			c.params[kv] = ""
		}
	}

	return &c
}

// NewClient Initiates a new REST request.
func NewClient(urlParts ...interface{}) *Client {
	urlPath := crconfig.Get("REST_CLIENT_BASE_PATH", BasePath)

	if len(urlParts) > 0 {
		if urlPath == "" {
			urlPath = format2string(urlParts[0], false)
			urlParts = urlParts[1:]
		}

		u, _ := url.Parse(urlPath)
		hasProto, slashpos := (u.Scheme != ""), 0
		if u.Path != "" {
			u.Path = strings.Trim(u.Path, "/")
			urlPath = u.String()
		} else if hasProto && u.Host == "" {
			slashpos = 1
		}
		prependSlash := (u.Scheme == "" && u.Host == "" && !strings.Contains(u.Path, "."))

		// join all trimmed parts together
		for i, part := range urlParts {
			if i >= slashpos {
				urlPath += "/"
			}
			urlPath += format2string(part, true)
		}
		if prependSlash {
			urlPath = "/" + urlPath
		}
	}
	Logger.Debugln("NewClient ", urlPath)

	c := Client{
		url:    urlPath,
		params: make(map[string]string),
		header: make(map[string]string),

		timeout:           crconfig.GetDuration("REST_CLIENT_TIMEOUT", 0) * time.Millisecond,
		retries:           int(crconfig.GetInt("REST_CLIENT_RETRIES", 0)),
		maxRetryTime:      crconfig.GetDuration("REST_CLIENT_MAX_RETRY_TIME", 60) * time.Second,
		ignoreCertificate: crconfig.GetBool("REST_CLIENT_IGNORE_CERTIFICATE", IgnoreCertificate),
	}

	return &c
}

func format2string(part interface{}, trimSlashes bool) string {
	switch p := part.(type) {
	case string:
		if trimSlashes {
			return strings.Trim(p, "/")
		}
		return p
	case float64:
		return strconv.FormatFloat(p, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(p), 'f', -1, 32)
	case int64:
		return strconv.FormatInt(p, 10)
	case int32:
		return strconv.FormatInt(int64(p), 10)
	case int16:
		return strconv.FormatInt(int64(p), 10)
	case int8:
		return strconv.FormatInt(int64(p), 10)
	case int:
		return strconv.FormatInt(int64(p), 10)
	case bool:
		return strconv.FormatBool(p)
	}
	return ""
}

func (c *Client) request(method string, postData *bytes.Buffer, urlParams bool) *Response {
	var reqBody io.Reader

	if method != http.MethodGet && postData != nil {
		reqBody = bytes.NewReader(postData.Bytes())
	}

	url := c.url

	if urlParams {
		url += c.getParamString(true)
	}

	Logger.Debugln(method, "to", url)

	if method != http.MethodGet {
		if postData != nil {
			Logger.Debugln(postData.String())
		} else {
			Logger.Debugln("postData is nil")
		}
	}

	resp := Response{}

	var req *http.Request
	req, resp.Error = http.NewRequest(method, url, reqBody)
	if resp.Error == nil {

		c.Header("Content-Type", "application/json")
		if c.token.isset {
			c.Header("Authorization", c.token.method+" "+c.token.value)
		}

		for k, v := range c.header {
			req.Header.Add(k, v)
			Logger.Debugln("add header " + k + ": " + v)
		}

		Logger.Debugln("do request now")

		var response *http.Response
		if CatchRequest != nil {
			response = CatchRequest(req)
		} else {
			client := http.Client{Timeout: c.timeout}

			if c.ignoreCertificate {
				client.Transport = &http.Transport{
					MaxIdleConnsPerHost: 10,
					TLSClientConfig: &tls.Config{
						MaxVersion:         tls.VersionTLS13,
						InsecureSkipVerify: true,
					},
				}
			}

			response, resp.Error = client.Do(req)
		}

		if response != nil {
			retryTime := 500 * time.Millisecond
			for i := 0; i < c.retries+1; i++ {
				resp.Status = response.StatusCode
				if resp.Error == nil {
					Logger.Debugln("reading response now")
					_, resp.Error = resp.Data.ReadFrom(response.Body)
					response.Body.Close()
					if resp.Status < 500 {
						break
					} else {
						time.Sleep(retryTime)
						if retryTime *= 2; retryTime > c.maxRetryTime {
							retryTime = c.maxRetryTime
						}
					}
				}
			}
		} else {
			Logger.Debugln("response nil")
		}
	}

	return &resp
}

func (c *Client) getParamString(forGET bool) string {
	var params string

	for k, v := range c.params {
		if params != "" {
			params += "&"
		}
		params += k
		params += "="
		if forGET {
			params += url.QueryEscape(v)
		} else {
			params += v
		}
	}

	if forGET && params != "" {
		if strings.Contains(c.url, "?") {
			params = "&" + params
		} else {
			params = "?" + params
		}
	}

	return params
}

// URL is the URL the request is made to
func (c *Client) URL() string {
	return c.url
}

// Param provides you to specify parameters
func (c *Client) Param(key, value interface{}) *Client {
	k := format2string(key, false)
	v := format2string(value, false)
	c.params[k] = v
	return c
}

// DeleteParam deletes a parameter from the list
func (c *Client) DeleteParam(key string) *Client {
	delete(c.params, key)
	return c
}

// Header provides you to specify parameters
func (c *Client) Header(key, value string) *Client {
	c.header[key] = value
	return c
}

// Timeout sets the requests timeout
func (c *Client) Timeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// Retries sets the number of times the same request is retried, if 5xx status is returned.
// The wait time is raised times two on every try. It starts with 500ms.
// Parameter max defines how the max wait time is (ms).
func (c *Client) Retries(retries, max int) *Client {
	c.retries = retries
	c.maxRetryTime = time.Duration(max)
	return c
}

// IgnoreCertificate ignores cerrtifacte errors
func (c *Client) IgnoreCertificate(flag bool) *Client {
	c.ignoreCertificate = flag
	return c
}

// Token provides you to specify a token.
// method is e.g. Bearer.
// this is put into header.
func (c *Client) Token(method, value string) *Client {
	c.token.method = method
	c.token.value = value
	c.token.isset = true
	return c
}

// Get makes a simple get request to an url, given to Client.
// Returns a pointer to a Buffer and error.
func (c *Client) Get() *Response {
	return c.request(http.MethodGet, nil, true)
}

// GetObject makes a simple get request and unmarshals the response into a given struct.
// Returns error on any serious errors.
func (c *Client) GetObject(result interface{}) *Response {

	resp := c.Get()

	if resp.Error == nil {
		if resp.Data.Len() > 0 {
			err := json.Unmarshal(resp.Data.Bytes(), result)
			if err != nil {
				// unmarshal nags about everything.
				// so just put out a warning, ignore the error
				// and let the caller decide...
				Logger.Warnln(err)
				Logger.Debugln(resp.Data.String())
				err = nil
			}
		} else {
			resp.Error = errors.New("no data received, nothing to unmarshal")
		}
	}

	return resp
}

// GetValues makes a simple get request and tries to parse the response to url values.
func (c *Client) GetValues(values *url.Values) *Response {

	resp := c.Get()

	if resp.Error == nil && resp.Data.Len() > 0 {
		*values, resp.Error = url.ParseQuery(resp.Data.String())
	}

	return resp
}

// Delete makes a simple DELETE request
func (c *Client) Delete() *Response {
	return c.request(http.MethodDelete, nil, true)
}

// DeletePostParams makes a DELETE request by sending params as post data
func (c *Client) DeletePostParams() *Response {
	params := c.getParamString(false)
	data := bytes.NewBufferString(params)

	return c.request(http.MethodDelete, data, false)
}

// Post makes a simple POST request
func (c *Client) Post(data *bytes.Buffer) *Response {
	return c.request(http.MethodPost, data, true)
}

// PostString makes a simple POST request, sending a string as data
func (c *Client) PostString(str string) *Response {
	return c.Post(bytes.NewBufferString(str))
}

// PostObject makes a simple POST request, marshaling the given struct to JSON
func (c *Client) PostObject(data interface{}) *Response {
	json, err := json.Marshal(data)
	if err == nil {
		data := bytes.NewBuffer(json)
		return c.Post(data)
	}
	return &Response{Error: err}
}

// PostParams makes a simple POST request, sending the given params as post data
func (c *Client) PostParams() *Response {
	params := c.getParamString(false)
	data := bytes.NewBufferString(params)

	return c.request(http.MethodPost, data, false)
}

// Put makes a simple PUT request
func (c *Client) Put(data *bytes.Buffer) *Response {
	return c.request(http.MethodPut, data, true)
}

// PutString makes a simple PUT request, sending a string as data
func (c *Client) PutString(str string) *Response {
	return c.Put(bytes.NewBufferString(str))
}

// PutObject makes a simple PUT request, marshaling the given struct to JSON
func (c *Client) PutObject(data interface{}) *Response {
	json, err := json.Marshal(data)
	if err == nil {
		data := bytes.NewBuffer(json)
		return c.Put(data)
	}
	return &Response{Error: err}
}

// PutParams makes a simple PUT request, sending the given params as post data
func (c *Client) PutParams() *Response {
	params := c.getParamString(false)
	data := bytes.NewBufferString(params)

	return c.request(http.MethodPut, data, false)
}
