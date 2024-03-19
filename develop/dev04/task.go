package main

import (
	"fmt"
	"sort"
	"strings"
)

// В условии сказано только про русский алфавит :), но если хотим другие - то надо в мапе в качестве ключа хранить отсоритированную строку, в остальном все так же.
func MapLetter(char rune, set *[33]int) {
	switch char {
	case 'а':
		set[0]++
	case 'б':
		set[1]++
	case 'в':
		set[2]++
	case 'г':
		set[3]++
	case 'д':
		set[4]++
	case 'е':
		set[5]++
	case 'ё':
		set[6]++
	case 'ж':
		set[7]++
	case 'з':
		set[8]++
	case 'и':
		set[9]++
	case 'й':
		set[10]++
	case 'к':
		set[11]++
	case 'л':
		set[12]++
	case 'м':
		set[13]++
	case 'н':
		set[14]++
	case 'о':
		set[15]++
	case 'п':
		set[16]++
	case 'р':
		set[17]++
	case 'с':
		set[18]++
	case 'т':
		set[19]++
	case 'у':
		set[20]++
	case 'ф':
		set[21]++
	case 'х':
		set[22]++
	case 'ц':
		set[23]++
	case 'ч':
		set[24]++
	case 'ш':
		set[25]++
	case 'щ':
		set[26]++
	case 'ъ':
		set[27]++
	case 'ы':
		set[28]++
	case 'ь':
		set[29]++
	case 'э':
		set[30]++
	case 'ю':
		set[31]++
	case 'я':
		set[32]++
	}
}

func GroupAnagrams(words []string) map[string][]string {

	mp := make(map[[33]int][]string, len(words))

	for _, v := range words {
		setLetters := [33]int{}
		v = strings.ToLower(v)
		for _, char := range v {
			MapLetter(char, &setLetters)
		}
		if _, ok := mp[setLetters]; ok {
			mp[setLetters] = append(mp[setLetters], v)
		} else {
			mp[setLetters] = make([]string, 0)
			mp[setLetters] = append(mp[setLetters], v)
		}
	}

	res := make(map[string][]string, len(mp))
	for _, v := range mp {
		if len(v) != 1 {
			key := v[0]
			v = v[1:]
			sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
			res[key] = v
		}
	}
	return res

}

func main() {
	words := []string{"ток", "пятак", "тяпка", "кот", "столик", "листок", "пятка", "слиток", "молоток"}

	result := GroupAnagrams(words)

	for key, v := range result {
		fmt.Printf("Key: %v, value: %v\n", key, v)
	}
}
