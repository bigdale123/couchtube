package dbmodels

type Video struct {
	ID           int    `db:"id" json:"id"`
	ChannelID    int    `db:"channel_id" json:"channelId"`
	URL          string `db:"url" json:"url"`
	SegmentStart int    `db:"segment_start" json:"segmentStart"`
	SegmentEnd   int    `db:"segment_end" json:"segmentEnd"`
}
