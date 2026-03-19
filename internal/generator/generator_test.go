package generator

import (
	"strings"
	"testing"
)

// проверка на соответствие длин алиаса и заданной длины
func TestGenerateShortAliasLength(t *testing.T) {
	alias := GenerateShortAlias()
	if len(alias) != aliasLength {
		t.Errorf("wanted len %d, got %d", aliasLength, len(alias))
	}
}

// проверка на наличие неразрешенного символа в алиасе
func TestGenerateShortAliasSymbols(t *testing.T) {
	alias := GenerateShortAlias()
	for _, char := range alias {
		if !strings.Contains(alphabet, string(char)) {
			t.Errorf("symbol: %v is not from alphabet %s", char, alphabet)
		}
	}
}
