package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SortUtil struct { //util
	Filename string
	K        int
	N        bool
	R        bool
	U        bool
	data     [][]string
}

func Start(args []string) error {
	su := SortUtil{Filename: args[0]}
	fs := flag.NewFlagSet("sortflags", flag.ContinueOnError)
	fs.IntVar(&su.K, "k", 0, "column for sorting")
	fs.BoolVar(&su.N, "n", false, "numerical sort")
	fs.BoolVar(&su.U, "u", false, "remove duplicates")
	fs.BoolVar(&su.R, "r", false, "reversed order")

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	file, err := os.Open(su.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = su.Process(file)
	if err != nil {
		return err
	}
	su.Sort()
	su.Print()

	return nil
}

func (su *SortUtil) Process(file *os.File) error {
	var hasEnoughCols bool

	su.data = make([][]string, 0)

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		su.data = append(su.data, make([]string, 0, 10))
		su.data[i] = append(su.data[i], line)
		su.data[i] = append(su.data[i], strings.Split(line, " ")...)
		if len(su.data[i])-1 >= su.K {
			hasEnoughCols = true
		}
		i++
	}

	if !hasEnoughCols {
		su.K = 0
	}

	if su.U {
		su.data = su.removeDupls()
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (su *SortUtil) removeDupls() [][]string {
	mp := make(map[string]struct{}, len(su.data))
	res := make([][]string, 0, len(su.data))
	var str string
	var err error
	for _, v := range su.data {
		if len(v) <= su.K {
			str = ""
			if su.N {
				str = "0"
			}
		} else {
			str = v[su.K]
			_, err = strconv.Atoi(str)
			if err != nil {
				str = "0"
			}
		}

		if _, ok := mp[str]; !ok {
			mp[str] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}

func (su *SortUtil) lessNum(i, j int) bool {
	var num1, num2 int
	var err error
	if len(su.data[i]) > su.K {
		num1, err = strconv.Atoi(su.data[i][su.K])
		if err != nil {
			num1 = 0
		}
	}
	if len(su.data[j]) > su.K {
		num2, err = strconv.Atoi(su.data[j][su.K])
		if err != nil {
			num2 = 0
		}
	}
	if num1 != num2 {
		return num1 < num2
	}
	return su.data[i][0] < su.data[j][0]
}

func (su *SortUtil) lessString(i, j int) bool {
	if len(su.data[i]) <= su.K && len(su.data[j]) <= su.K {
		return true
	}
	if len(su.data[i]) <= su.K && len(su.data[j]) > su.K {
		return true
	}
	if len(su.data[i]) > su.K && len(su.data[j]) <= su.K {
		return false
	}
	return su.data[i][su.K] < su.data[j][su.K]
}

func (su *SortUtil) Sort() {
	suLess := su.lessString

	if su.N {
		suLess = su.lessNum
	}

	if su.R {
		sort.Slice(su.data, func(i, j int) bool {
			return !suLess(i, j)
		})
		return
	}
	sort.Slice(su.data, suLess)

}

func (su *SortUtil) Print() {
	for i := 0; i < len(su.data); i++ {
		fmt.Printf("%s\r\n", su.data[i][0])
	}
}

func main() {
	fmt.Println(os.Args)
	err := Start(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
