package structural

import (
	"log"
	"math/big"
	"os"
	"time"
)

type factFunc func(n int) *big.Int

func factorialWrapper(factFunc factFunc, logger *log.Logger) factFunc {
	return func(n int) *big.Int {
		defer func(t time.Time) {
			logger.Printf("took=%v, n=%v, result=%v", time.Since(t), n, "factFunc(n)")
		}(time.Now())

		return factFunc(n)
	}
}

func factorialWithGoRoutine(n int) *big.Int {
	var result = big.NewInt(1)
	ch := make(chan *big.Int)
	if n == 1 || n == 0 {
		return result
	}

	for i := 1; i <= n; i++ {
		go func(n int) {
			ch <- big.NewInt(int64(n))
		}(i)
	}

	for i := 0; i < n; i++ {
		result.Mul(result, <-ch)
	}
	return result
}

func factorialWithoutGoRoutine(n int) *big.Int {
	var result = big.NewInt(1)
	if n == 1 || n == 0 {
		return result
	}

	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

func factorialWithBatching(n int) *big.Int {
	batchSize := 80
	if n <= batchSize {
		return factorialWithoutGoRoutine(n)
	}

	var result = big.NewInt(1)
	ch := make(chan *big.Int)
	if n == 1 || n == 0 {
		return result
	}

	for i := 0; i < n; i += batchSize {
		go func(start, end int) {
			chunkResult := big.NewInt(1)
			for j := start; j <= end; j++ {
				chunkResult.Mul(chunkResult, big.NewInt(int64(j)))
			}
			ch <- chunkResult
		}(i+1, min(i+batchSize, n))
	}

	for i := 0; i < n; i += batchSize {
		result.Mul(result, <-ch)
	}
	return result
}

func DecoratorRunner(n int) {
	factorialWrapper(factorialWithGoRoutine, log.New(os.Stdout, "Factorial with go routine: ", log.LstdFlags))(n)
	factorialWrapper(factorialWithoutGoRoutine, log.New(os.Stdout, "Factorial without go routine: ", log.LstdFlags))(n)
	factorialWrapper(factorialWithBatching, log.New(os.Stdout, "Factorial with batching aproach: ", log.LstdFlags))(n)
}
