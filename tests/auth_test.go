package tests

import (
	"github.com/ElvinEga/gofiber_starter/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "test123"
	hashed := utils.HashPassword(password)

	assert.NotEmpty(t, hashed)
	assert.True(t, utils.CheckPasswordHash(password, hashed))
	assert.False(t, utils.CheckPasswordHash("wrongpassword", hashed))
}

func TestGenerateUsername(t *testing.T) {
	name := "John Doe"
	username := utils.GenerateUsername(name)

	assert.NotEmpty(t, username)
	assert.Contains(t, username, "johndoe")
}

func TestGenerateSecureToken(t *testing.T) {
	token := utils.GenerateSecureToken(32)

	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes = 64 hex characters
}

func TestGenerateUUID(t *testing.T) {
	uuid := utils.GenerateUUID()

	assert.NotEmpty(t, uuid)
	assert.Len(t, uuid.String(), 36) // Standard UUID length
}
