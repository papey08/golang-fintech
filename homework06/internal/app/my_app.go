package app

import (
	"context"
	"errors"
	valid "github.com/papey08/go_course_validation"
	"homework6/internal/ads"
)

var ValidationError = errors.New("invalid ad struct")
var WrongUserError = errors.New("different userID from ad.UserID")

type MyApp struct {
	Repository
}

func (a MyApp) CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error) {
	adToAdd := ads.Ad{
		ID:        0,
		Title:     title,
		Text:      text,
		AuthorID:  userID,
		Published: false,
	}
	if err := valid.Validate(adToAdd); err != nil {
		return ads.Ad{}, ValidationError
	}
	if ad, err := a.AddAd(ctx, adToAdd); err != nil {
		return ads.Ad{}, err
	} else {
		return ad, nil
	}
}

func (a MyApp) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (ads.Ad, error) {
	if ad, err := a.GetAdByID(ctx, adID); err != nil {
		return ads.Ad{}, err
	} else if ad.AuthorID != userID {
		return ads.Ad{}, WrongUserError
	} else {
		ad.Published = published
		return *ad, nil
	}
}

func (a MyApp) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (ads.Ad, error) {
	if ad, err := a.GetAdByID(ctx, adID); err != nil {
		return ads.Ad{}, err
	} else if ad.AuthorID != userID {
		return ads.Ad{}, WrongUserError
	} else if validErr := valid.Validate(ads.Ad{
		ID:        ad.ID,
		Title:     title,
		Text:      text,
		AuthorID:  ad.AuthorID,
		Published: ad.Published,
	}); validErr != nil {
		return ads.Ad{}, ValidationError
	} else {
		ad.Title = title
		ad.Text = text
		return *ad, nil
	}
}
