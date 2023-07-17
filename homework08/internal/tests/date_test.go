package tests

import (
	"github.com/stretchr/testify/assert"
	"homework8/internal/model/ads"
	"testing"
)

func TestCreationDate(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	_, err = client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err := client.getAdByID(0)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.CreationDate.Year, ads.CurrentDate().Year)
	assert.Equal(t, response.Data.CreationDate.Month, ads.CurrentDate().Month)
	assert.Equal(t, response.Data.CreationDate.Day, ads.CurrentDate().Day)
}

func TestUpdateDate(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err := client.createAd(123, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(123, response.Data.ID, true)
	assert.NoError(t, err)

	assert.Equal(t, response.Data.UpdateDate.Year, ads.CurrentDate().Year)
	assert.Equal(t, response.Data.UpdateDate.Month, ads.CurrentDate().Month)
	assert.Equal(t, response.Data.UpdateDate.Day, ads.CurrentDate().Day)
}
