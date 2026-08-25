package internal

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
)

// Clamp restricts a value between min and max
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Must panics if error is not nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Check logs error if present
func Check(err error) bool {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	return true
}

// ReadFile reads entire file into byte slice
func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// WriteFile writes byte slice to file
func WriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// RandomSeed generates a cryptographic random seed
func RandomSeed(length int) ([]byte, error) {
	seed := make([]byte, length)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random seed: %w", err)
	}
	return seed, nil
}

// Shuffle randomizes the order of integers using Fisher-Yates algorithm
func Shuffle(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		arr[i], arr[j.Int64()] = arr[j.Int64()], arr[i]
	}
}

// AbsInt returns absolute value of integer
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat returns absolute value of float64
func AbsFloat(x float64) float64 {
	return math.Abs(x)
}

// Pow returns base raised to power
func Pow(base, exp int) int {
	return int(math.Pow(float64(base), float64(exp)))
}
