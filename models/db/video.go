package dbmodels

type Video struct {
	ID           int    `db:"id" json:"id"`
	ChannelID    int    `db:"channel_id" json:"channel_id"`
	URL          string `db:"url" json:"url"`
	SegmentStart int    `db:"segment_start" json:"segment_start"`
	SegmentEnd   int    `db:"segment_end" json:"segment_end"`
}
