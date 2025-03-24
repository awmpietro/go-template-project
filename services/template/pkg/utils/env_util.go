package utils

import (
	"os"
	"strconv"
)

func GetEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func GetEnvAsFloat(name string, defaultVal float64) float64 {
	valStr := os.Getenv(name)
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := os.Getenv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
