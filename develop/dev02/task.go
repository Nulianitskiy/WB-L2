package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func unpack(s string) (string, error) {

	var res strings.Builder
	var ram rune
	res.Grow(len(s))
	var isEscapeOn bool
	for i, char := range s {
		if !isEscapeOn && string(char) == "\\" {
			if i == len([]rune(s))-1 {
				return "", fmt.Errorf("uncorrect string")
			}
			isEscapeOn = true
			continue
		}
		if isEscapeOn || unicode.IsLetter(char) {
			res.WriteRune(char)
			ram = char
			isEscapeOn = false
			continue
		}
		if unicode.IsDigit(char) {
			if ram == 0 {
				return "", fmt.Errorf("bad string")
			}
			num, err := strconv.Atoi(string(char))
			if num == 0 {
				return "", fmt.Errorf("zero repeating symbol")
			}
			if err != nil {
				return "", err
			}
			for i := 0; i < num-1; i++ {
				res.WriteRune(ram)
			}
			ram = 0
		}
	}
	return res.String(), nil
}

func main() {
	var s string
	fmt.Scan(&s)
	str, err := unpack(s)
	if err != nil {
		fmt.Println(err.Error())
	}
	if str == "" {
		fmt.Println("yess")
	}
	fmt.Println(str)
}
