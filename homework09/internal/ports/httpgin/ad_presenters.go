package httpgin

import "homework9/internal/model/ads"

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type searchAdRequest struct {
	Pattern string `json:"pattern"`
}

type filterRequest struct {
	PublishedBy bool `json:"published_by"`

	AuthorBy bool  `json:"author_by"`
	AuthorID int64 `json:"author_id"`

	DateBy bool     `json:"date_by"`
	Date   ads.Date `json:"date"`
}

type deleteAdRequest struct {
	UserID int64 `json:"user_id"`
}
