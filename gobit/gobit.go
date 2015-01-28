package gobit

import (
	"net/http"
	"net"
	"time"
	"bytes"
	"fmt"
	"encoding/json"
)

// TODO: とりあえずgobrake参考
var defaultHttpClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 3*time.Second)
		},
		ResponseHeaderTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// Client is an errbit client.
type Client struct {
	client  *http.Client
	noticeURL string
	options Options
}

func New(opt Options) *Client {
	
	c := Client{}
	c.client = defaultHttpClient
	c.noticeURL = opt.createNoticeBaseURL()
	c.options = opt
	
	return &c
}

func (c Client) Send(notice *Notice) error {
	data, err := json.Marshal(notice)
	if err != nil {
		return err
	}

	u := c.options.createNoticeBaseURL()
	
	res, err := c.client.Post(u, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response code %d", res.StatusCode)
	}

	return nil
}
