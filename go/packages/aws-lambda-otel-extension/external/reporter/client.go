package reporter

import (
	"aws-lambda-otel-extension/external/lib"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type PostData struct {
	body   []byte
	path   string
	trying bool
	lock   *sync.Mutex
	next   *PostData
	prev   *PostData
}

type HttpClient struct {
	HttpClient *http.Client
	settings   *lib.UserSettings
	eg         *errgroup.Group
	stackLast  *PostData
	stackLock  *sync.Mutex
}

type transformDataType func() ([]byte, error)

func NewHttpClient(settings *lib.UserSettings) *HttpClient {
	return &HttpClient{
		HttpClient: &http.Client{},
		settings:   settings,
		eg:         &errgroup.Group{},
		stackLock:  &sync.Mutex{},
	}
}

func (c *HttpClient) Flush() {
	current := c.stackLast
	for {
		if current == nil {
			return
		}
		if !current.trying {
			param := current
			c.eg.Go(func() error {
				return c.syncPost(param)
			})
		}
		current = current.prev
	}
}

func (c *HttpClient) Post(path string, body []byte) error {
	c.stackLock.Lock()
	defer c.Flush()
	defer c.stackLock.Unlock()
	data := PostData{body: body, path: path, lock: &sync.Mutex{}, trying: false, prev: c.stackLast}
	c.stackLast = &data
	return nil
}

func (c *HttpClient) removeStack(postData *PostData) {
	c.stackLock.Lock()
	defer c.stackLock.Unlock()
	if postData.prev == nil {
		c.stackLast = postData.next
	} else {
		postData.prev.next = postData.next
	}
	if postData.next != nil {
		postData.next.prev = postData.prev
	}
}

func (c *HttpClient) syncPost(postData *PostData) (err error) {
	postData.lock.Lock()
	postData.trying = true
	defer func() {
		if err == nil {
			c.removeStack(postData)
		} else {
			postData.trying = false
			postData.lock.Unlock()
		}
	}()
	fmt.Printf(">> sending post '%s'\n", postData.path)
	req, err := http.NewRequest("POST", postData.path, bytes.NewBuffer(postData.body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	extraParams, err := url.ParseQuery(c.settings.Common.Destination.RequestHeaders)
	if err != nil {
		return err
	}

	for key, value := range extraParams {
		req.Header.Set(key, value[0])
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with status %s", resp.Status)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf(">> post sent '%s'\n", postData.path)
	return nil
}

func (c *HttpClient) PostLogs(transformData transformDataType) {
	c.eg.Go(func() error {
		body, err := transformData()
		if err != nil {
			return err
		}
		// fmt.Printf("sending body '%s\n'", body)
		return c.Post(c.settings.Logs.Destination, body)
	})
}

func (c *HttpClient) PostMetric(metrics string) {
	c.eg.Go(func() error {
		return c.Post(c.settings.Metrics.Destination, []byte(metrics))
	})
}

func (c *HttpClient) PostRequest(request []byte) {
	c.eg.Go(func() error {
		return c.Post(c.settings.Request.Destination, request)
	})
}

func (c *HttpClient) Shutdown() error {
	for {
		// TODO: implement a timeout for this loop, less than a second
		if c.stackLast == nil {
			break
		}
		c.Flush()
		time.Sleep(20 * time.Millisecond)
	}
	c.HttpClient.CloseIdleConnections()
	return c.eg.Wait()
}
