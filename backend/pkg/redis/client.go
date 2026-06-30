package redis

import (
	"context"

	goredis "github.com/go-redis/redis/v8"
)

type Client struct {
	client *goredis.Client
}

func NewClient(addr, password string, db int) *Client {
	return &Client{
		client: goredis.NewClient(&goredis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Raw() *goredis.Client {
	return c.client
}
