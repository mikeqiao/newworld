// Author: mike.qiao
// File:world_test
// Date:2022/3/14 16:56

package newworld

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	Start(nil)
	//	ChanClose()
}

func Benchmark_11(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func ChanClose() {
	a := make(chan int, 100)
	for i := 1; i < 100; i++ {
		a <- i
	}
	close(a)
	go func() {
		for {
			d, ok := <-a
			if ok {
				fmt.Println("a:", d)
			} else {
				fmt.Println("closed")
				break
			}
		}
	}()

	for {
		d, ok := <-a
		if ok {
			fmt.Println("b:", d)
		} else {
			fmt.Println("closed")
			break
		}
	}
}
