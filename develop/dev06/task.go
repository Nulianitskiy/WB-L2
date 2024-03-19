package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cut struct {
	F            string
	D            string
	S            bool
	Filename     string
	fieldIndices []int
	tillEnd      bool
	start        int
}

func Start(args []string) error {
	cu := Cut{Filename: args[len(args)-1], fieldIndices: make([]int, 0, 10)}
	fs := flag.NewFlagSet("cutflags", flag.ContinueOnError)
	fs.StringVar(&cu.F, "f", "", "columns")
	fs.StringVar(&cu.D, "d", "\t", "delimeter")
	fs.BoolVar(&cu.S, "s", false, "only with delimeter")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	file, err := os.Open(cu.Filename)
	if err != nil {
		return err
	}
	defer file.Close()
	err = cu.ParseF()
	if err != nil {
		return err
	}
	err = cu.Process(file)
	if err != nil {
		return err
	}

	return nil
}

func (cu *Cut) ParseF() error {
	if strings.Contains(cu.F, ",") {
		fields := strings.Split(cu.F, ",")
		for _, v := range fields {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			cu.fieldIndices = append(cu.fieldIndices, n-1)
		}
		return nil
	}
	if len(cu.F) == 3 {
		num1, err := strconv.Atoi(string(cu.F[0]))
		if err != nil {
			return err
		}
		num2, err := strconv.Atoi(string(cu.F[2]))
		if err != nil {
			return err
		}
		for i := num1; i <= num2; i++ {
			cu.fieldIndices = append(cu.fieldIndices, i-1)
		}
		return nil
	}
	if len(cu.F) == 2 {
		if string(cu.F[0]) == "-" {
			num1, err := strconv.Atoi(string(cu.F[1]))
			if err != nil {
				return err
			}
			for i := 1; i <= num1; i++ {
				cu.fieldIndices = append(cu.fieldIndices, i-1)
			}
			return nil
		}
		if string(cu.F[1]) == "-" {
			num1, err := strconv.Atoi(string(cu.F[0]))
			if err != nil {
				return err
			}
			cu.tillEnd = true
			cu.start = num1 - 1
			return nil
		}
		return nil
	}

	num, err := strconv.Atoi(string(cu.F[0]))
	if err != nil {
		return err
	}
	cu.fieldIndices = append(cu.fieldIndices, num-1)
	return nil
}

func (cu *Cut) Process(f *os.File) error {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), cu.D)
		if len(fields) == 1 && cu.S {
			continue
		}
		if len(fields) == 1 && !cu.S {
			fmt.Printf("%s\r\n", fields[0])
			continue
		}
		if cu.tillEnd {
			for i := cu.start; i < len(fields); i++ {
				if i == len(fields)-1 {
					fmt.Printf("%s\r\n", fields[i])
					continue
				}
				fmt.Printf("%s%s", fields[i], cu.D)
			}
			continue
		}
		for i, v := range cu.fieldIndices {
			if v < len(fields) {
				if v == len(fields)-1 || i == len(cu.fieldIndices)-1 {
					fmt.Printf("%s", fields[v])
					continue
				}
				if i+1 < len(cu.fieldIndices) && cu.fieldIndices[i+1] >= len(fields) {
					fmt.Printf("%s", fields[v])
					continue
				}
				fmt.Printf("%s%s", fields[v], cu.D)
			}
		}
		fmt.Printf("\r\n")
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
