package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	// 短链码字符集
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 短链码长度
	shortCodeLength = 7
)

var (
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mutex      sync.Mutex
)

// GenerateShortCode 生成短链码
func GenerateShortCode(originalURL string) string {
	// 使用MD5哈希
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hasher.Write([]byte(time.Now().String()))
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	// 取前7位作为短码
	shortCode := md5Hash[:shortCodeLength]

	// 确保只包含允许的字符
	for i, c := range shortCode {
		if !strings.ContainsRune(charset, c) {
			shortCode = shortCode[:i] + string(charset[seededRand.Intn(len(charset))]) + shortCode[i+1:]
		}
	}

	return shortCode
}

// GenerateRandomShortCode 生成随机短链码
func GenerateRandomShortCode() string {
	mutex.Lock()
	defer mutex.Unlock()

	b := make([]byte, shortCodeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IsValidShortCode 验证短链码是否合法
func IsValidShortCode(code string) bool {
	if len(code) != shortCodeLength {
		return false
	}

	for _, c := range code {
		if !strings.ContainsRune(charset, c) {
			return false
		}
	}

	return true
}
