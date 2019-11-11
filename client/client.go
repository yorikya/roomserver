package client

import "time"

type Client struct {
	ClientID, MovementSen, TempSen, AirCond, LightMain, LightSec string
	LastSeen                                                     time.Time
}

type Message struct {
	MovementSen string `json:"movementSen"`
	TempSen     string `json:"tempSen"`
	AirCond     string `json:"airCond"`
	LightMain   string `json:"lightMain"`
	LightSec    string `json:"lightSec"`
	Action      string `json:"action"`
}

func (c *Client) UpdateState(m Message) {
	c.MovementSen = m.MovementSen
	c.TempSen = m.TempSen
	c.AirCond = m.AirCond
	c.LightMain = m.LightMain
	c.LightSec = m.LightSec
	c.LastSeen = time.Now()
}

func NewClient(clientID string) *Client {
	return &Client{
		ClientID: clientID,
	}
}
