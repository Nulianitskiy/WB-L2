package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupAnagrams(t *testing.T) {
	input := [][]string{
		{"ток", "пятак", "тяпКа", "КОТ", "столик", "листок", "пяткА", "слиток", "молоток"},
		{"ток", "пот", "гот"},
	}
	exp := []map[string][]string{
		{
			"ток":    {"кот"},
			"пятак":  {"пятка", "тяпка"},
			"столик": {"листок", "слиток"},
		},
		{},
	}

	ans := make([]map[string][]string, 0, len(exp))
	for _, v := range input {
		ans = append(ans, GroupAnagrams(v))
	}
	assert.Equal(t, exp, ans, "not equal")
}
