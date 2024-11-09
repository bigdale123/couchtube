package dbmodels

type Video struct {
	ID           int    `db:"id" json:"id"`
	ChannelID    int    `db:"channel_id" json:"channelId"`
	URL          string `db:"url" json:"url"`
	SectionStart int    `db:"section_start" json:"sectionStart"`
	SectionEnd   int    `db:"section_end" json:"sectionEnd"`
}
