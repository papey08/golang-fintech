package adrepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"homework_extra/internal/app"
	"homework_extra/internal/model/ads"
	"homework_extra/internal/model/errs"
	"homework_extra/internal/model/filter"
	"time"
)

const (
	getAdByIDQuery = `
		SELECT * FROM ads
		WHERE id = ($1);`

	addAdQuery = `
		INSERT INTO ads(title, text, author_id, published, creation_date, update_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`

	updateAdFieldsQuery = `
		UPDATE ads
		SET title = $1,
			text = $2,
			author_id = $3,
			published = $4,
			creation_date = $5,
			update_date = $6
		WHERE id = $7;`

	getAdsListQuery = `
		SELECT * FROM ads
		WHERE
		    (($1) OR published) AND
		    ((NOT $2) OR author_id = $3) AND
		    ((NOT $4) OR creation_date = $5);`

	searchAdsQuery = `
		SELECT * FROM ads
		WHERE (title LIKE $1 OR text LIKE $1);`

	deleteAdQuery = `
		DELETE FROM ads
		WHERE id = $1;
	`
)

type AdRepository struct {
	pgx.Conn
}

func (r *AdRepository) GetAdByID(ctx context.Context, id int64) (ads.Ad, error) {
	row := r.QueryRow(ctx, getAdByIDQuery, id)
	var ad ads.Ad
	var creationDate, updateDate time.Time
	if err := row.Scan(&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID, &ad.Published, &creationDate, &updateDate); err == pgx.ErrNoRows {
		return ads.Ad{}, errs.AdNotExist
	} else if err != nil {
		return ads.Ad{}, errs.AdRepositoryError
	} else {
		ad.CreationDate.Year, ad.CreationDate.Month, ad.CreationDate.Day = creationDate.Date()
		ad.UpdateDate.Year, ad.UpdateDate.Month, ad.UpdateDate.Day = updateDate.Date()
		return ad, nil
	}
}

func (r *AdRepository) AddAd(ctx context.Context, ad ads.Ad) (ads.Ad, error) {
	var addedAdID int64
	err := r.QueryRow(ctx, addAdQuery,
		ad.Title,
		ad.Text,
		ad.AuthorID,
		ad.Published,
		fmt.Sprintf("%d-%d-%d", ad.CreationDate.Year, ad.CreationDate.Month, ad.CreationDate.Day),
		fmt.Sprintf("%d-%d-%d", ad.UpdateDate.Year, ad.UpdateDate.Month, ad.UpdateDate.Day)).Scan(&addedAdID)
	if err != nil {
		return ads.Ad{}, errs.AdRepositoryError
	}
	ad.ID = addedAdID
	return ad, nil
}

func (r *AdRepository) UpdateAdFields(ctx context.Context, idToUpdate int64, newAd ads.Ad) (ads.Ad, error) {
	e, err := r.Exec(ctx, updateAdFieldsQuery,
		newAd.Title,
		newAd.Text,
		newAd.AuthorID,
		newAd.Published,
		fmt.Sprintf("%d-%d-%d", newAd.CreationDate.Year, newAd.CreationDate.Month, newAd.CreationDate.Day),
		fmt.Sprintf("%d-%d-%d", newAd.UpdateDate.Year, newAd.UpdateDate.Month, newAd.UpdateDate.Day),
		idToUpdate)
	if err != nil {
		return ads.Ad{}, errs.AdRepositoryError
	} else if e.RowsAffected() == 0 {
		return ads.Ad{}, errs.AdNotExist
	} else {
		return ads.Ad{
			ID:           idToUpdate,
			Title:        newAd.Title,
			Text:         newAd.Text,
			AuthorID:     newAd.AuthorID,
			Published:    newAd.Published,
			CreationDate: newAd.CreationDate,
			UpdateDate:   newAd.UpdateDate,
		}, nil
	}
}

func (r *AdRepository) GetAdsList(ctx context.Context, f filter.Filter) ([]ads.Ad, error) {
	rows, err := r.Query(ctx, getAdsListQuery,
		f.PublishedBy,
		f.AuthorBy,
		f.AuthorID,
		f.DateBy,
		fmt.Sprintf("%d-%d-%d", f.Date.Year, f.Date.Month, f.Date.Day))
	defer rows.Close()
	if err != nil {
		return nil, errs.AdRepositoryError
	}
	adList := make([]ads.Ad, 0)
	for rows.Next() {
		var tempAd ads.Ad
		var creationDate, updateDate time.Time
		_ = rows.Scan(&tempAd.ID, &tempAd.Title, &tempAd.Text, &tempAd.AuthorID, &tempAd.Published, &creationDate, &updateDate)
		tempAd.CreationDate.Year, tempAd.CreationDate.Month, tempAd.CreationDate.Day = creationDate.Date()
		tempAd.UpdateDate.Year, tempAd.UpdateDate.Month, tempAd.UpdateDate.Day = updateDate.Date()
		adList = append(adList, tempAd)
	}
	return adList, nil
}

func (r *AdRepository) SearchAds(ctx context.Context, pattern string) ([]ads.Ad, error) {
	rows, err := r.Query(ctx, searchAdsQuery, pattern+"%")
	defer rows.Close()
	if err != nil {
		return nil, errs.AdRepositoryError
	}
	adList := make([]ads.Ad, 0)
	for rows.Next() {
		var tempAd ads.Ad
		var creationDate, updateDate time.Time
		_ = rows.Scan(&tempAd.ID, &tempAd.Title, &tempAd.Text, &tempAd.AuthorID, &tempAd.Published, &creationDate, &updateDate)
		tempAd.CreationDate.Year, tempAd.CreationDate.Month, tempAd.CreationDate.Day = creationDate.Date()
		tempAd.UpdateDate.Year, tempAd.UpdateDate.Month, tempAd.UpdateDate.Day = updateDate.Date()
		adList = append(adList, tempAd)
	}
	return adList, nil
}

func (r *AdRepository) DeleteAd(ctx context.Context, idToDelete int64) error {
	e, err := r.Exec(ctx, deleteAdQuery, idToDelete)
	if err != nil {
		return errs.AdRepositoryError
	} else if e.RowsAffected() == 0 {
		return errs.AdNotExist
	} else {
		return nil
	}
}

func New(conn *pgx.Conn) app.AdRepository {
	return &AdRepository{
		Conn: *conn,
	}
}
