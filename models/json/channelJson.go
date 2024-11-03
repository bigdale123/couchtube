package jsonmodels

type VideoJson struct {
	Url          string `json:"url"`
	SegmentStart int    `json:"segment_start"`
	SegmentEnd   int    `json:"segment_end"`
}

type ChannelJson struct {
	Name   string      `json:"name"`
	Videos []VideoJson `json:"videos"`
}

type ChannelsJson struct {
	Channels []ChannelJson `json:"channels"`
}
