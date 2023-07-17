package filter

import "homework10/internal/model/ads"

type Filter struct {
	PublishedBy bool

	AuthorBy bool
	AuthorID int64

	DateBy bool
	Date   ads.Date
}
