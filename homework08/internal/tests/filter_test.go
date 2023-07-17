package tests

import (
	"github.com/stretchr/testify/assert"
	"homework8/internal/model/ads"
	"testing"
)

func TestNotPublishedOnlyFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedOnly: false,
		AuthorBy:      false,
		AuthorID:      0,
		DateBy:        false,
		Date:          ads.Date{},
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 2)
}

func TestAuthorByFilter(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)
	_, err = client.createUser(100, "pokemosha", "email2@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(100, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(123, "привет", "мир")
	assert.NoError(t, err)

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedOnly: false,
		AuthorBy:      true,
		AuthorID:      response.Data.ID,
		DateBy:        false,
		Date:          ads.Date{},
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 2)
}

func TestDateByFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(123, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(123, "привет", "мир")
	assert.NoError(t, err)

	curDate := ads.CurrentDate()

	adsList, err := client.listAdsWithFilters(filterRequest{
		PublishedOnly: false,
		AuthorBy:      false,
		AuthorID:      0,
		DateBy:        true,
		Date:          curDate,
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 3)

	curDate.Year--

	adsList, err = client.listAdsWithFilters(filterRequest{
		PublishedOnly: false,
		AuthorBy:      false,
		AuthorID:      0,
		DateBy:        true,
		Date:          curDate,
	})
	assert.NoError(t, err)
	assert.Len(t, adsList.Data, 0)
}
