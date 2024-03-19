package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Grep struct {
	A        int
	B        int
	C        int
	Cc       bool
	I        bool
	V        bool
	F        bool
	N        bool
	data     []string
	Filename string
	Pattern  string
	count    int
}

func Start(args []string) error {
	gp := Grep{Pattern: args[len(args)-2], Filename: args[len(args)-1], data: make([]string, 0, 1000)}
	fs := flag.NewFlagSet("grepflags", flag.ContinueOnError)
	fs.IntVar(&gp.A, "A", 0, "strings after")
	fs.IntVar(&gp.B, "B", 0, "strings before")
	fs.IntVar(&gp.C, "C", 0, "context")
	fs.BoolVar(&gp.Cc, "c", false, "count")
	fs.BoolVar(&gp.I, "i", false, "ignore-case")
	fs.BoolVar(&gp.V, "v", false, "invert")
	fs.BoolVar(&gp.F, "F", false, "exact match")
	fs.BoolVar(&gp.N, "n", false, "string number")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	file, err := os.Open(gp.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gp.Process(file)
	if err != nil {
		return err
	}

	if gp.F {
		gp.Pattern = regexp.QuoteMeta(gp.Pattern)
	}
	if gp.I {
		gp.Pattern = fmt.Sprintf(`(?i)%s`, gp.Pattern)
	}

	re := regexp.MustCompile(gp.Pattern)

	matched := gp.Find(re)

	gp.PrintOut(matched)

	return nil
}

func (gp *Grep) PrintOut(mp map[int]struct{}) {
	if gp.Cc {
		fmt.Printf("%d\r\n", gp.count)
		return
	}
	for i, v := range gp.data {
		if _, ok := mp[i]; ok {
			if gp.N {
				fmt.Printf("%d:%s\r\n", i+1, v)
			} else {
				fmt.Printf("%s\r\n", v)
			}
		}
	}
}

func (gp *Grep) AddIndexes(direction bool, cur int, outlen int, mp *map[int]struct{}) {
	if direction {
		cur++
	} else {
		cur--
	}
	for j, c := cur, 0; j < len(gp.data) && c < outlen && j >= 0; {
		(*mp)[j] = struct{}{}
		if direction {
			j++
		} else {
			j--
		}
		c++
	}
}

func (gp *Grep) Find(re *regexp.Regexp) map[int]struct{} {

	matchedIndexes := make(map[int]struct{}, len(gp.data)/2)

	for i, v := range gp.data {
		if re.MatchString(v) && !gp.V || !re.MatchString(v) && gp.V {
			matchedIndexes[i] = struct{}{}
			gp.count++
			if gp.A > 0 {
				gp.AddIndexes(true, i, gp.A, &matchedIndexes)
			}
			if gp.B > 0 {
				gp.AddIndexes(false, i, gp.B, &matchedIndexes)
			}
			if gp.C > 0 {
				gp.AddIndexes(true, i, gp.C, &matchedIndexes)
				gp.AddIndexes(false, i, gp.C, &matchedIndexes)
			}
		}
	}
	return matchedIndexes

}

func (gp *Grep) Process(f *os.File) error {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		gp.data = append(gp.data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := Start(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
