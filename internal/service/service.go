package service

type ICalculator interface {
	Sum(a, b int32) int32
	PrimeDecomposition(a int32) chan int32
}

type CalculatorService struct {
}

func (s CalculatorService) PrimeDecomposition(a int32) chan int32 {
	prime := make(chan int32, 5)

	go func() {
		var k int32 = 2

		for {
			if a%k == 0 {
				a = a / k
				prime <- k
			} else {
				k = k + 1
			}

			if a <= 1 {
				close(prime)
				break
			}
		}
	}()

	return prime
}

func (s CalculatorService) Sum(a, b int32) int32 {
	return a + b
}
