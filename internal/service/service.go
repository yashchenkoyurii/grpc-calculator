package service

type ICalculator interface {
	Sum(a, b int32) int32
	PrimeDecomposition(a int32) chan int32
	ComputeAverage(numbers chan int32, result chan float32)
	FindMax(numbers chan int32, max chan int32)
}

type CalculatorService struct {
}

func (s CalculatorService) FindMax(numbers chan int32, max chan int32) {
	var m = <-numbers
	max <- m

	for n := range numbers {
		if n > m {
			m = n
		}

		max <- m
	}

	close(max)
}

func (s CalculatorService) ComputeAverage(numbers chan int32, result chan float32) {
	var count float32
	var sum float32

	for n := range numbers {
		sum += float32(n)
		count++
	}

	result <- sum / count
}

func (s CalculatorService) PrimeDecomposition(a int32) chan int32 {
	prime := make(chan int32, 5)

	go func() {
		defer close(prime)

		var k int32 = 2

		for a > 1 {
			if a%k == 0 {
				a = a / k
				prime <- k
			} else {
				k = k + 1
			}
		}
	}()

	return prime
}

func (s CalculatorService) Sum(a, b int32) int32 {
	return a + b
}
