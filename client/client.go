package client

import (
	"log"
	"time"

	"proxy-stream/helper/datetime"
)

type Client struct {
	indentificator string
	connectedAt    *time.Time
}

func New(identificator string) Client {
	return Client{
		indentificator: identificator,
	}
}

func (c *Client) Connect() {
	now := time.Now()
	c.connectedAt = &now
	log.Printf("client:🟩 %s", c.indentificator)

}

func (c *Client) Disconnect() {
	duration := c.duration()

	log.Printf("client:🟥 %s %s", c.indentificator, duration)
}

func (c *Client) GetID() string {
	return c.indentificator
}

func (c *Client) duration() string {
	return datetime.Datetime2duration(c.connectedAt)
}
