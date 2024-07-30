package urlshortener

import (
	"crypto/sha1"
	"encoding/hex"
)

func Shorten(url string) string {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha[:8] // short URL of length 8
}
