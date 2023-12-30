package app

import (
	"context"
	"homework8/internal/model/ads"
	"homework8/internal/model/filter"
	"homework8/internal/model/users"
)

type App interface {

	// CreateAd makes new ad with given parameters and inserts it to the repository
	CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error)

	// ChangeAdStatus changes status of ad with ID adID
	ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (ads.Ad, error)

	// UpdateAd updates title and text of ad with ID adID
	UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (ads.Ad, error)

	// GetAdByID returns ad with adID
	GetAdByID(ctx context.Context, adID int64) (ads.Ad, error)

	// GetAdsList returns slice of ads passed through filter f
	GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error)

	// SearchAds returns slice of ads which title is equal to pattern
	SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error)

	// GetUserByID returns ad with userID
	GetUserByID(ctx context.Context, userID int64) (users.User, error)

	// CreateUser makes new user with given parameters and inserts it to the repository
	CreateUser(ctx context.Context, id int64, nickname string, email string) (users.User, error)

	// UpdateUser updates title and text of user with ID userID
	UpdateUser(ctx context.Context, userID int64, nickname string, email string) (users.User, error)
}

type AdRepository interface {

	// GetAdByID returns ad with specified ID
	GetAdByID(ctx context.Context, id int64) (ads.Ad, error)

	// AddAd inserts ad to repository
	AddAd(ctx context.Context, ad ads.Ad) (ads.Ad, error)

	// UpdateAdFields replaces all fields of ad in repository with id equal to idToUpdate to fields of newAd
	UpdateAdFields(ctx context.Context, idToUpdate int64, newAd ads.Ad) (ads.Ad, error)

	// GetAdsList returns slice of ads passed through filter f
	GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error)

	// SearchAds returns slice of ads which title is equal to pattern
	SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error)
}

type UserRepository interface {

	// GetUserByID returns users with specified ID
	GetUserByID(ctx context.Context, id int64) (users.User, error)

	// AddUser inserts users to repository
	AddUser(ctx context.Context, user users.User) (users.User, error)

	// UpdateUserFields replaces all fields of users in repository with id equal to idToUpdate to fields of newUser
	UpdateUserFields(ctx context.Context, idToUpdate int64, newUser users.User) (users.User, error)
}

func NewApp(adRepo AdRepository, userRepo UserRepository) App {
	return &MyApp{
		AdRepository:   adRepo,
		UserRepository: userRepo,
	}
}