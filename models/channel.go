package models

type Video struct {
	Url string `json:"url"`
}

type Channel struct {
	Id     int8    `json:"id"`
	Type   string  `json:"type"`
	Videos []Video `json:"videos"`
}

type Channels struct {
	Channels []Channel `json:"channels"`
}
