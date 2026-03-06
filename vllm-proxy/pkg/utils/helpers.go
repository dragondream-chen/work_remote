package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func FormatAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func FormatURL(host string, port int, path string) string {
	return fmt.Sprintf("http://%s:%d%s", host, port, path)
}

func CalculatePrefillScore(requestLength int) float64 {
	lengthScore := float64(requestLength) / 4.0
	return lengthScore*0.0345 + 120.0745
}

func CalculateDecodeScore(requestLength int) float64 {
	return float64(requestLength)
}

func ExponentialBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay * time.Duration(1<<(attempt-1))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

func ContainsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, str string) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != str {
			result = append(result, v)
		}
	}
	return result
}
