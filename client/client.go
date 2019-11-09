package client

import "time"

type Client struct {
	ClientID, MovementSen, TempSen, AirCond, LightMain, LightSec string
	LastSeen                                                     time.Time
}

type Response struct {
	MovementSen string `json:"movementSen"`
	TempSen     string `json:"tempSen"`
	AirCond     string `json:"airCond"`
	LightMain   string `json:"lightMain"`
	LightSec    string `json:"lightSec"`
}

func (c *Client) UpdateState(r Response) {
	c.MovementSen = r.MovementSen
	c.TempSen = r.TempSen
	c.AirCond = r.AirCond
	c.LightMain = r.LightMain
	c.LightSec = r.LightSec
	c.LastSeen = time.Now()
}

func NewClient(clientID string) *Client {
	return &Client{
		ClientID: clientID,
	}
}
