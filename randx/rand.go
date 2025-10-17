package randx

import (
	"math/rand"
)

// 生成验证码Token
func GenerateRandomToken(length int, charset string) string {
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}

	return string(token)
}

func GenerateRandomTokenWithSeed(length int, charset string, seed int64) string {
	source := rand.NewSource(seed)
	r := rand.New(source)
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[r.Intn(len(charset))]
	}
	return string(token)
}
