package verification

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateVerificationCode() string {
	minRange, maxRange := 100000, 999999

	return strconv.Itoa(rand.Intn(maxRange-minRange+1) + minRange)
}

func GenerateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}
