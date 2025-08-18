package client

import (
	"fmt"
	"log"
	"time"
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
	duration := ""
	if c.connectedAt != nil {
		elapsed := time.Since(*c.connectedAt)
		years := int(elapsed.Hours() / 24 / 365)
		months := int(elapsed.Hours()/24/30) % 12
		days := int(elapsed.Hours()/24) % 30
		hours := int(elapsed.Hours()) % 24
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		if years > 0 {
			duration += fmt.Sprintf("%d years ", years)
		}
		if months > 0 {
			duration += fmt.Sprintf("%d months ", months)
		}
		if days > 0 {
			duration += fmt.Sprintf("%dd:", days)
		}
		if hours > 0 {
			duration += fmt.Sprintf("%dh:", hours)
		}
		if minutes > 0 {
			duration += fmt.Sprintf("%dm:", minutes)
		}
		if seconds > 0 || duration == "" {
			duration += fmt.Sprintf("%ds", seconds)
		}
	}

	return duration
}
