package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculatorService_Sum(t *testing.T) {
	svc := &CalculatorService{}

	var a int32 = 1
	var b int32 = 1
	c := a + b

	assert.Equal(t, c, svc.Sum(a, b))
}

func TestCalculatorService_PrimeDecomposition(t *testing.T) {
	svc := &CalculatorService{}

	var a int32 = 120
	var b = []int32{2, 2, 2, 3, 5}
	i := 0
	ch := svc.PrimeDecomposition(a)

	for v := range ch {
		assert.Equal(t, v, b[i])
		i++
	}
}

func TestCalculatorService_ComputeAverage(t *testing.T) {
	svc := &CalculatorService{}
	numbers := []int32{1, 2, 3, 4}
	var res float32 = 2.5
	input := make(chan int32)
	average := make(chan float32)

	go svc.ComputeAverage(input, average)
	for _, n := range numbers {
		input <- n
	}
	close(input)

	assert.Equal(t, res, <-average)

}

func TestCalculatorService_FindMax(t *testing.T) {
	var actual []int32
	numbers := []int32{1, 20, 15, 21}
	maxs := []int32{1, 20, 20, 21}

	in := make(chan int32)
	out := make(chan int32)

	svc := &CalculatorService{}

	go svc.FindMax(in, out)
	go func() {
		for _, n := range numbers {
			in <- n
		}

		close(in)
	}()

	for m := range out {
		actual = append(actual, m)
	}

	assert.Equal(t, maxs, actual)
}
