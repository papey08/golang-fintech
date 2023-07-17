package adrepo

import (
	"context"
	"homework9/internal/app"
	"homework9/internal/model/ads"
	"homework9/internal/model/errs"
	"homework9/internal/model/filter"
	"strings"
	"sync"
)

type AdRepository struct {
	data   map[int64]ads.Ad
	freeID int64
	mu     sync.RWMutex
}

func (r *AdRepository) GetAdByID(ctx context.Context, id int64) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, errs.AdRepositoryError
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()
		if ad, ok := r.data[id]; !ok {
			return ads.Ad{}, errs.AdNotExist
		} else {
			return ad, nil
		}
	}
}

func (r *AdRepository) AddAd(ctx context.Context, ad ads.Ad) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, errs.AdRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		ad.ID = r.freeID
		r.data[r.freeID] = ad
		r.freeID++
		return ad, nil
	}
}

func (r *AdRepository) UpdateAdFields(ctx context.Context, idToUpdate int64, newAd ads.Ad) (ads.Ad, error) {
	select {
	case <-ctx.Done():
		return ads.Ad{}, errs.AdRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.data[idToUpdate]; !ok {
			return ads.Ad{}, errs.AdNotExist
		} else {
			ad := ads.Ad{
				ID:           idToUpdate,
				Title:        newAd.Title,
				Text:         newAd.Text,
				AuthorID:     newAd.AuthorID,
				Published:    newAd.Published,
				CreationDate: newAd.CreationDate,
				UpdateDate:   newAd.UpdateDate,
			}
			r.data[idToUpdate] = ad
			return ad, nil
		}
	}
}

func (r *AdRepository) GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, errs.AdRepositoryError
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()
		adList := make([]ads.Ad, 0, len(r.data))
		for _, ad := range r.data {
			if !f.PublishedBy && !ad.Published {
				continue
			}
			if f.DateBy && ad.CreationDate != f.Date {
				continue
			}
			if f.AuthorBy && ad.AuthorID != f.AuthorID {
				continue
			}
			adList = append(adList, ad)
		}
		return adList, nil
	}
}

func (r *AdRepository) SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error) {
	select {
	case <-ctx.Done():
		return nil, errs.AdRepositoryError
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()
		resList := make([]ads.Ad, 0, len(r.data))
		for _, ad := range r.data {
			if ad.Title == pattern || strings.HasPrefix(ad.Title, pattern) {
				resList = append(resList, ad)
			}
		}
		return resList, nil
	}
}

func (r *AdRepository) DeleteAd(ctx context.Context, idToDelete int64) error {
	select {
	case <-ctx.Done():
		return errs.AdRepositoryError
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.data[idToDelete]; !ok {
			return errs.AdNotExist
		} else {
			delete(r.data, idToDelete)
			return nil
		}
	}
}

func New() app.AdRepository {
	return &AdRepository{
		data:   make(map[int64]ads.Ad),
		freeID: 0,
	}
}
