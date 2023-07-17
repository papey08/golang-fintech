package filter

import "homework8/internal/model/ads"

type Filter struct {
	PublishedOnly bool

	AuthorBy bool
	AuthorID int64

	DateBy bool
	Date   ads.Date
}
