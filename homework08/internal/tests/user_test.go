package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(123))
	assert.Equal(t, response.Data.Nickname, "papey08")
	assert.Equal(t, response.Data.Email, "email@golang.com")
}

func TestGetUserByID(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err = client.getUserByID(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(123))
	assert.Equal(t, response.Data.Nickname, "papey08")
	assert.Equal(t, response.Data.Email, "email@golang.com")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser(123, "papey08", "email@golang.com")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, response.Data.Nickname, "email2@golang.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(123))
	assert.Equal(t, response.Data.Nickname, "papey08")
	assert.Equal(t, response.Data.Email, "email2@golang.com")
}
