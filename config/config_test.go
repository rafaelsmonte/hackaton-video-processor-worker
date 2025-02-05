package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	godotenv.Load(".env")

	cfg := NewConfig()

	assert.IsType(t, Config{}, cfg)
}

func TestGetConfig(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")

	cfg := NewConfig()

	value := cfg.GetConfig("TEST_KEY")
	assert.Equal(t, "test_value", value)
}

func TestGetConfig_EmptyKey(t *testing.T) {
	cfg := NewConfig()

	value := cfg.GetConfig("NON_EXISTENT_KEY")
	assert.Empty(t, value)
}

func TestGetConfig_WithDotenv(t *testing.T) {
	os.Setenv("TEST_ENV", "loaded_from_env")

	godotenv.Load(".env")
	cfg := NewConfig()

	value := cfg.GetConfig("TEST_ENV")
	assert.Equal(t, "loaded_from_env", value)
}
