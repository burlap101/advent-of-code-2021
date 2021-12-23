package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"io"

	"github.com/burlap101/advent-of-code-2021/day6/fish"
)

func populateFishes() []fish.LanternFish {
	dat, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var fishes []fish.LanternFish
	for _, inputTime := range strings.Split(string(dat), ",") {
		initialTime, err := strconv.Atoi(inputTime)
		if err != nil {
			panic(err)
		}
		fishes = append(fishes, fish.LanternFish{Timer: initialTime})
	}

	err = os.WriteFile("fishlog.txt", dat, 0o777)
	if err != nil {
		panic(err)
	}

	return fishes
}

/* Storing fishes in memory */
func serviceFishes(fishes []fish.LanternFish) ([]fish.LanternFish, error) {
	for i := range fishes {
		err := fishes[i].ServiceTimer(false)
		if err != nil {
			newBorn, err := fishes[i].SpawnNew()
			if err != nil {
				panic(err)
			}
			fishes = append(fishes, newBorn)
		}
	}
	return fishes, nil
}

/* Processes provided reading and returns number of new spawns */
func processReading(dat []byte) (int, fish.LanternFish, error) {
	newSpawns := 0
	timer, err := strconv.Atoi(string(dat[0]))
	if err != nil {
		panic(err)
	}
	fish := fish.LanternFish{Timer: timer}
	err = fish.ServiceTimer(false)
	if err != nil {
		newSpawns += 1
	}
	return newSpawns, fish, nil
}

/* Storing fish state in logfile. returns new spawns for a day */
func serviceFishLog() (int, error) {
	fp, err := os.OpenFile("fishlog.txt", os.O_RDWR, 0o0777)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	dat := make([]byte, 2)
	newSpawns := 0
	totalFish := 0

	for err = nil; err != io.EOF; totalFish++ {
		
		readed, readErr := fp.Read(dat)
		_, seekErr := fp.Seek(int64(-readed), 1)

		if seekErr != nil {
			panic(seekErr)
		}
		
		if readed == 1 {
			fmt.Println("here")
			spawnCount, fishy, err := processReading(dat)
			if err != nil {
				panic(err)
			}
			newSpawns += spawnCount
			dat[1] = ','
			dat[0] = strconv.Itoa(fishy.Timer)[0]
			_, err = fp.Write(dat)
			if err != nil {
				panic(err)
			}
			break
		} else if readErr == io.EOF {
			fmt.Println("end of file reached")
			break
		} else if readed == 2 && dat[1] == ',' {
			spawnCount, fishy, err := processReading(dat)
			if err != nil {
				panic(err)
			}
			newSpawns += spawnCount
			dat[0] = strconv.Itoa(fishy.Timer)[0]
			_, err = fp.Write(dat)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(readed)
			panic("i'm not sure what is going on here...")
		}
	}
	for i := 0; i < newSpawns; i++ {
		_, err := fp.Write([]byte{'8', ','})
		if err != nil {
			panic(err)
		}
	}

	return totalFish + newSpawns, nil
}

func main() {
	fmt.Println("Advent of code day 6")
	fmt.Println("Part1")
	fishes := populateFishes()
	
	fmt.Println("Initial count fishes:", len(fishes))

	for i := 0; i < 80; i++ {
		var err error
		fishes, err = serviceFishes(fishes)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Day%d count: %d\n", i+1, len(fishes))
	}

	fmt.Println("\nPart2")
	fishes = populateFishes()
	timers := make([]uint64, 9)
	for _, fishy := range fishes {
		timers[fishy.Timer] += 1
	}
	for i := 0; i < 256; i++ {
		motherCount := timers[0]
		timers = timers[1:]
		timers[6] += motherCount
		timers = append(timers, motherCount)
		fullCount := uint64(0)
		for _, fcount := range timers {
			fullCount += fcount
		}
		fmt.Printf("Day%d count: %d\n", i+1, fullCount)
	}
}