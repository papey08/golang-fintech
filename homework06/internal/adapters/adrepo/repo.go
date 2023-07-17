package adrepo

import (
	"context"
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

var AdNotExist = errors.New("no ad in repository with required id")
var ContextErr = errors.New("context is done")

type Repository struct {
	data []ads.Ad
}

func (r *Repository) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, ContextErr
	default:
		if id >= int64(len(r.data)) {
			return nil, AdNotExist
		} else {
			return &r.data[id], nil
		}
	}
}

func (r *Repository) AddAd(ctx context.Context, ad ads.Ad) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, ContextErr
	default:
		ad.ID = int64(len(r.data))
		r.data = append(r.data, ad)
		return ad, nil
	}
}

func New() app.Repository {
	return &Repository{
		data: make([]ads.Ad, 0),
	}
}
