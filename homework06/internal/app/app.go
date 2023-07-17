package app

import (
	"context"
	"homework6/internal/ads"
)

type App interface {

	// CreateAd makes new ad with given parameters and inserts it to the repository
	CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error)

	// ChangeAdStatus changes status of ad with ID adID
	ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (ads.Ad, error)

	// UpdateAd updates title and text of ad with ID adID
	UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (ads.Ad, error)
}

type Repository interface {

	// GetAdByID returns a pointer to ad with specified ID for making some changes with it
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)

	// AddAd inserts ad to repository
	AddAd(ctx context.Context, ad ads.Ad) (ads.Ad, error)
}

func NewApp(repo Repository) App {
	return &MyApp{Repository: repo}
}
