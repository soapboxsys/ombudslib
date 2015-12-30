package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//record_type := flag.String("type", "", "Either bltn or endo")
	flag.Parse()

	//pipe := false

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to stdin")
		//pipe = true
	} else {
		fmt.Println("stdin is from a terminal")
	}
}
