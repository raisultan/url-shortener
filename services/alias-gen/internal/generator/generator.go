package generator

import "strings"

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateAlias(n int64) string {
	return base62Encode(n)
}

func base62Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}
	var chars []string
	base := int64(len(alphabet))
	for n > 0 {
		rem := n % base
		chars = append([]string{string(alphabet[rem])}, chars...)
		n = n / base
	}
	return strings.Join(chars, "")
}
