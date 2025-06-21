package services_test

import (
	"hr-system/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateFormat_Invalid(t *testing.T) {
	_, err1 := time.Parse("2006-01-02", "invalid")
	_, err2 := time.Parse("2006-01-02", "2024-01-01")
	err := services.ValidateFormat(err1, err2)
	assert.Equal(t, services.ErrInvalidDateFormat, err)
}

func TestValidateFormat_Valid(t *testing.T) {
	_, err1 := time.Parse("2006-01-02", "2024-01-01")
	_, err2 := time.Parse("2006-01-02", "2024-01-15")
	err := services.ValidateFormat(err1, err2)
	assert.Nil(t, err)
}

func TestValidateCycle_EndBeforeStart(t *testing.T) {
	start := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	err := services.ValidateCycle(start, end)
	assert.Equal(t, services.ErrDateRange, err)
}

func TestValidateCycle_ValidRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	err := services.ValidateCycle(start, end)
	assert.Nil(t, err)
}
