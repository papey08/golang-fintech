package ads

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCurrentDate(t *testing.T) {
	date := CurrentDate()
	assert.Equal(t, time.Now().UTC().Day(), date.Day)
	assert.Equal(t, time.Now().UTC().Month(), date.Month)
	assert.Equal(t, time.Now().UTC().Year(), date.Year)
}
