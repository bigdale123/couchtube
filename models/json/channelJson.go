package jsonmodels

type VideoJson struct {
	Url          string `json:"url"`
	SegmentStart int    `json:"segmentStart"`
	SegmentEnd   int    `json:"segmentEnd"`
}

type ChannelJson struct {
	Name   string      `json:"name"`
	Videos []VideoJson `json:"videos"`
}

type ChannelsJson struct {
	Channels []ChannelJson `json:"channels"`
}
