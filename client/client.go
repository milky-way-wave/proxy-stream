package client

import (
	"log"
	"time"
)

var cookie_key_name string = "cid"

type Client struct {
	indentificator string
	ticker         int64
}

func New(identificator string) Client {
	return Client{
		indentificator: identificator,
		ticker:         0,
	}
}

func (c Client) Connect() {
	c.ticker = time.Now().Unix()
	log.Printf("%s - connected", c.indentificator)

}

func (c Client) Disconnect() {
	if c.ticker == 0 {
		return
	}

	now := time.Now().Unix()
	c.ticker = now - c.ticker

}
