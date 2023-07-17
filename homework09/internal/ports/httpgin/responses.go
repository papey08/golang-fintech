package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/model/ads"
	"homework9/internal/model/users"
)

type adResponse struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Text         string   `json:"text"`
	AuthorID     int64    `json:"author_id"`
	Published    bool     `json:"published"`
	CreationDate ads.Date `json:"creation_date"`
	UpdateDate   ads.Date `json:"update_date"`
}

type adsResponse []adResponse

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(ads []ads.Ad) *gin.H {
	resp := make(adsResponse, 0, len(ads))
	for _, ad := range ads {
		resp = append(resp, adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: ad.CreationDate,
			UpdateDate:   ad.UpdateDate,
		})
	}
	return &gin.H{
		"data":  resp,
		"error": nil,
	}
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
