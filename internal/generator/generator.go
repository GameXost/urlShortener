package generator

import (
	"math/rand"
	"strings"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	aliasLength = 10
)

func GenerateShortAlias() string {
	var alias strings.Builder
	alias.Grow(aliasLength)
	for range 10 {
		alias.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return alias.String()
}
