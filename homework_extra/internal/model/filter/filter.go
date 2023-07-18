package filter

import "homework_extra/internal/model/ads"

type Filter struct {
	PublishedBy bool

	AuthorBy bool
	AuthorID int64

	DateBy bool
	Date   ads.Date
}
