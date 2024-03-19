package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	testTime := [][]time.Duration{
		{
			10 * time.Hour,
			10 * time.Second,
			10 * time.Minute,
			10 * time.Hour,
		},
		{
			5 * time.Second,
			5 * time.Second,
			5 * time.Second,
		},
	}
	gaps := []int{11, 6}

	for i, v := range testTime {
		ar := make([]<-chan interface{}, 0, len(v))
		st := time.Now().Second()
		for _, k := range v {
			ar = append(ar, sig(k))
		}
		<-or(ar...)
		end := time.Now().Second()

		if end-st > gaps[i] {
			t.Errorf("working more time than need")
		}
	}
}
