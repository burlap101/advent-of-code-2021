package main

import (
	"fmt"
	"os"

	"github.com/burlap101/advent-of-code-2021/day16/transport"
)

func ProcessFile(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func main() {
	var err error
	fmt.Println("Advent of code day 16")
	hs := ProcessFile("input.txt")
	p := transport.Packet{
		HexString: hs,
		Parent: nil,
	}
	p.Version, err = transport.ExtractVersion(p.HexString)
	if err != nil {
		panic(err)
	}
	_, _, err = p.Children()
	if err != nil {
		panic(err)
	}
	fmt.Println("Version sum =", int(p.Version) + transport.VersionSum)
}