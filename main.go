package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Adder interface {
	Add(x, y int) int
}

type AdderFunc func(x, y int) int

func (a AdderFunc) Add(x, y int) int {
	return a(x, y)
}

func main() {
	a := AdderFunc(
		func(x, y int) (result int) {
			defer func(t time.Time) {
				log.SetOutput(os.Stdout)
				log.Printf("took=%v, x=%v, y=%v, result=%v", time.Since(t), x, y, result)
			}(time.Now())
			return x + y
		},
	)
	fmt.Println(Do(a))
}

func Do(adder Adder) int {
	return adder.Add(1, 2)
}
