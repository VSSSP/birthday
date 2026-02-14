package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
	assert.True(t, CheckPassword(password, hashed))
}

func TestCheckPassword_Wrong(t *testing.T) {
	hashed, _ := HashPassword("correct-password")
	assert.False(t, CheckPassword("wrong-password", hashed))
}
