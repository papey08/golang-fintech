package app

import (
	"context"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"homework10/internal/model/users"

	valid "github.com/papey08/golang-fintech/validation"
)

type MyApp struct {
	AdRepository
	UserRepository
}

func (a MyApp) CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error) {
	if _, err := a.UserRepository.GetUserByID(ctx, userID); err != nil {
		return ads.Ad{}, err
	}
	adToAdd := ads.Ad{
		ID:           0,
		Title:        title,
		Text:         text,
		AuthorID:     userID,
		Published:    false,
		CreationDate: ads.CurrentDate(),
		UpdateDate:   ads.CurrentDate(),
	}
	if err := valid.Validate(adToAdd); err != nil {
		return ads.Ad{}, errs.ValidationError
	}
	return a.AddAd(ctx, adToAdd)
}

func (a MyApp) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (ads.Ad, error) {
	if _, err := a.UserRepository.GetUserByID(ctx, userID); err != nil {
		return ads.Ad{}, err
	}
	if ad, err := a.AdRepository.GetAdByID(ctx, adID); err != nil {
		return ads.Ad{}, err
	} else if ad.AuthorID != userID {
		return ads.Ad{}, errs.WrongUserError
	} else {
		ad.Published = published
		ad.UpdateDate = ads.CurrentDate()
		return a.UpdateAdFields(ctx, adID, ad)
	}
}

func (a MyApp) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (ads.Ad, error) {
	if _, err := a.UserRepository.GetUserByID(ctx, userID); err != nil {
		return ads.Ad{}, err
	}
	if ad, err := a.AdRepository.GetAdByID(ctx, adID); err != nil {
		return ads.Ad{}, err
	} else if ad.AuthorID != userID {
		return ads.Ad{}, errs.WrongUserError
	} else if validErr := valid.Validate(ads.Ad{
		ID:           ad.ID,
		Title:        title,
		Text:         text,
		AuthorID:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: ad.CreationDate,
		UpdateDate:   ad.UpdateDate,
	}); validErr != nil {
		return ads.Ad{}, errs.ValidationError
	} else {
		ad.Title = title
		ad.Text = text
		ad.UpdateDate = ads.CurrentDate()
		return a.UpdateAdFields(ctx, adID, ad)
	}
}

func (a MyApp) GetAdByID(ctx context.Context, adID int64) (ads.Ad, error) {
	return a.AdRepository.GetAdByID(ctx, adID)
}

func (a MyApp) GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error) {
	return a.AdRepository.GetAdsList(ctx, f)
}

func (a MyApp) SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error) {
	return a.AdRepository.SearchAds(ctx, pattern)
}

func (a MyApp) DeleteAd(ctx context.Context, userID int64, adID int64) error {
	user, err := a.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	ad, err := a.AdRepository.GetAdByID(ctx, adID)
	if err != nil {
		return err
	}

	if user.ID != ad.AuthorID {
		return errs.WrongUserError
	} else {
		return a.AdRepository.DeleteAd(ctx, adID)
	}
}

func (a MyApp) GetUserByID(ctx context.Context, userID int64) (users.User, error) {
	return a.UserRepository.GetUserByID(ctx, userID)
}

func (a MyApp) CreateUser(ctx context.Context, nickname string, email string) (users.User, error) {
	userToAdd := users.User{
		ID:       0,
		Nickname: nickname,
		Email:    email,
	}
	return a.AddUser(ctx, userToAdd)
}

func (a MyApp) UpdateUser(ctx context.Context, userID int64, nickname string, email string) (users.User, error) {
	userToUpdate := users.User{
		ID:       0,
		Nickname: nickname,
		Email:    email,
	}
	return a.UserRepository.UpdateUserFields(ctx, userID, userToUpdate)
}

func (a MyApp) DeleteUser(ctx context.Context, userID int64) error {
	return a.UserRepository.DeleteUser(ctx, userID)
}
