package internal

import (
	"fmt"
	"sync"
)

type loan struct {
	id    int
	price int
	rate  float64
}

func NewLoan(id, price int, rate float64) *loan {
	return &loan{
		id:    id,
		price: price,
		rate:  rate,
	}
}

func (l *loan) calc(year int) {
	months := 12 * year
	interest := 0
	for i := 0; i < months; i++ {
		balance := l.price * (months - i) / months
		interest += int(float64(balance) * l.rate / 12)
	}
	fmt.Printf("{id: %d, year: %d, total: %d, interest: %d}\n",
		l.id, year, l.price+interest, interest)
}

func Worker(l *loan, years chan int, wg *sync.WaitGroup) {
	for y := range years {
		l.calc(y)
		wg.Done()
	}
}

var Generator = func(done <-chan struct{}, nums ...int) <-chan int {
	intStream := make(chan int, len(nums))
	go func() {
		defer close(intStream)
		for _, i := range nums {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func Multiply(done <-chan struct{}, intStream <-chan int, multiplier int) <-chan int {
	mulStream := make(chan int)
	go func() {
		defer close(mulStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case mulStream <- i * multiplier:
			}
		}
	}()
	return mulStream
}

func Add(done <-chan struct{}, intStream <-chan int, additive int) <-chan int {
	addStream := make(chan int)
	go func() {
		defer close(addStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addStream <- i + additive:
			}
		}
	}()
	return addStream
}
