package tests

import (
	"homework9/internal/model/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestGetAdByID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	var id0, id1 int64 = 0, 1

	ad0, err := client.getAdByID(id0)
	assert.NoError(t, err)

	ad1, err := client.getAdByID(id1)
	assert.NoError(t, err)

	assert.Equal(t, ad0.Data.ID, id0)
	assert.Equal(t, ad1.Data.ID, id1)
}

func TestSearchAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello world", "привет мир")
	assert.NoError(t, err)

	ads, err := client.searchAds("hello")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("papey08", "email@golang.com")
	assert.NoError(t, err)

	ad, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	err = client.deleteAd(user.Data.ID, ad.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAdByID(ad.Data.ID)
	assert.Error(t, err, errs.AdNotExist)
}
