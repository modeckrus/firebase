package postedition

import "time"

//PostPub is post what will be publiched and displayable
type PostPub struct {
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Nick        string    `json:"nick"`
	Avatar      string    `json:"avatar"`
	Likes       int       `json:"likes"`
	LastComment string    `json:"lastcomment"`
	UserID      string    `json:"userid"`
	Hasattach   bool      `json:"hasattach"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"createdat"`
}
