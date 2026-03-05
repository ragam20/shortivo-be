package link

import (
	"crypto/sha1"
	"errors"
	"net/url"
)

// Generate unique code for short url
func encodeBase62(num int64) (string, error) {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if num == 0 {
		return "", errors.New("num can't be zero")
	}

	var result []byte
	for num > 0 {
		reminder := num % 62
		result = append(result, base62Chars[reminder])
		num /= 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

// validate url
func validateURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}

func visitorHash(ipAddr, userAgent, referer, device string) [20]byte {
	ref := referer
	if ref == "" {
		ref = "none"
	}

	raw := ipAddr + "|" + userAgent + "|" + ref + "|" + device
	return sha1.Sum([]byte(raw))
}
