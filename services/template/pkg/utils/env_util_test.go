package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	assert.Equal(t, 42, GetEnvAsInt("TEST_INT", 0))
	assert.Equal(t, 0, GetEnvAsInt("NON_EXISTENT_ENV", 0))
}

func TestGetEnvAsFloat(t *testing.T) {
	os.Setenv("TEST_FLOAT", "3.14")
	defer os.Unsetenv("TEST_FLOAT")
	assert.Equal(t, 3.14, GetEnvAsFloat("TEST_FLOAT", 0.0))
	assert.Equal(t, 0.0, GetEnvAsFloat("NON_EXISTENT_ENV", 0.0))
}

func TestGetEnvAsBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")
	assert.Equal(t, true, GetEnvAsBool("TEST_BOOL", false))
	assert.Equal(t, false, GetEnvAsBool("NON_EXISTENT_ENV", false))
}
