package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	correctInput := []string{
		"abcd",
		"\\4\\5ab5c1d2e4",
		"qwe\\\\4",
		"qwe\\45",
		"",
	}
	uncorrectInpur := []string{
		"45",
		"qwe45",
		"45qwe",
		"qwe\\\\\\\\\\",
	}

	correctAnswers := []string{
		"abcd",
		"45abbbbbcddeeee",
		"qwe\\\\\\\\",
		"qwe44444",
		"",
	}

	a, _ := unpack("")
	if a == "" {
		fmt.Println("yeee")
	}
	answers := make([]string, 0, len(correctInput))
	for _, v := range correctInput {
		ans, err := unpack(v)
		if err != nil {
			t.Errorf("expected nil, got %s", err.Error())
		}
		answers = append(answers, ans)
	}
	assert.Equal(t, answers, correctAnswers, "значения не совпадают")

	for _, v := range uncorrectInpur {
		_, err := unpack(v)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	}
}
