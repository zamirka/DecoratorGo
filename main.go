package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Adder interface {
	Add(x, y int) int
}

type AdderFunc func(x, y int) int

func (a AdderFunc) Add(x, y int) int {
	return a(x, y)
}

type AdderMiddleware func(Adder) Adder

func WrapLogger(logger *log.Logger) AdderMiddleware {
	return func(a Adder) Adder {
		fn := func(x, y int) (result int) {
			defer func(t time.Time) {
				logger.Printf("took=%v, x=%v, y=%v, result=%v", time.Since(t), x, y, result)
			}(time.Now())
			return a.Add(x, y)
		}
		return AdderFunc(fn)
	}
}

func WrapperCache(cache *sync.Map) AdderMiddleware {
	return func(a Adder) Adder {
		fn := func(x, y int) int {
			key := fmt.Sprintf("x=%dy=%d", x, y)
			val, ok := cache.Load(key)
			if ok {
				return val.(int)
			}
			result := a.Add(x, y)
			cache.Store(key, result)
			return result
		}
		return AdderFunc(fn)
	}
}

func Chain(outer AdderMiddleware, middleware ...AdderMiddleware) AdderMiddleware {
	return func(a Adder) Adder {
		topIndex := len(middleware) - 1
		for i := range middleware {
			a = middleware[topIndex-i](a)
		}
		return outer(a)
	}
}

func main() {
	var a Adder = AdderFunc(
		func(x, y int) int {
			return x + y
		},
	)
	a = Chain(
		WrapLogger(log.New(os.Stdout, "DecoratorGo: ", 1)),
		WrapperCache(&sync.Map{}),
	)(a)

	for i := 1; i < 10; i++ {
		x := rand.Intn(2)
		y := rand.Intn(2)
		a.Add(x, y)
	}
}
