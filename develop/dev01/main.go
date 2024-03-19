package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

// getExactTime получает точное время с использованием библиотеки NTP.
func getExactTime(serverName string) {
	time, err := ntp.Time(serverName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err.Error())
		os.Exit(1)
	}
	fmt.Println(time)
}

// "pool.ntp.org"
func main() {
	getExactTime("pool.ntp.org")
}
