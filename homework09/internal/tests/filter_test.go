package tests

import (
	"github.com/stretchr/testify/assert"
	"homework9/internal/model/ads"
	"testing"
)

func TestNotPublishedOnlyFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedBy: true,
		AuthorBy:    false,
		AuthorID:    0,
		DateBy:      false,
		Date:        ads.Date{},
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 2)
}

func TestAuthorByFilter(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)
	_, err = client.createUser("pokemosha", "email2@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(0, "привет", "мир")
	assert.NoError(t, err)

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedBy: true,
		AuthorBy:    true,
		AuthorID:    response.Data.ID,
		DateBy:      false,
		Date:        ads.Date{},
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 2)
}

func TestDateByFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(0, "привет", "мир")
	assert.NoError(t, err)

	curDate := ads.CurrentDate()

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedBy: true,
		AuthorBy:    false,
		AuthorID:    0,
		DateBy:      true,
		Date:        curDate,
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 3)

	curDate.Year--

	adsList, err = client.listAdsWithFilters(filterRequest{
		PublishedBy: true,
		AuthorBy:    false,
		AuthorID:    0,
		DateBy:      true,
		Date:        curDate,
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 0)
}
