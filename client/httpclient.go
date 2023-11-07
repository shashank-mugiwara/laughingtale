package client

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var http_client *resty.Client = resty.
	New().
	SetTimeout(10 * time.Second)

func HttpClient() *resty.Client {
	return http_client
}
