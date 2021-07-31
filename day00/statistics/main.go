package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

var meanFlag bool
var medianFlag bool
var modeFlag bool
var sdFlag bool

func init() {
	flag.BoolVar(&meanFlag, "mean", false, "display average")
	flag.BoolVar(&medianFlag, "median", false, "display median")
	flag.BoolVar(&modeFlag, "mode", false, "display mode")
	flag.BoolVar(&sdFlag, "sd", false, "display standard deviation")
}

// numbers - not empty sorted sequence
func median(numbers []int) (result float64) {
	if len(numbers)%2 == 0 {
		result = float64(numbers[len(numbers)/2]+numbers[len(numbers)/2-1]) / float64(2)
	} else {
		result = float64(numbers[len(numbers)/2])
	}
	return result
}

// numbers - not empty sequence
func mode(numbers []int) (result int) {
	countNumbers := make(map[int]int)
	maxCount := 1
	result = numbers[0]
	for _, number := range numbers {
		countNumbers[number]++
		if countNumbers[number] > maxCount && number < result {
			maxCount = countNumbers[number]
			result = number
		}
	}
	return result
}

// numbers - not empty sequence
func sd(numbers []int) (result float64) {
	arithmeticMean := mean(numbers)
	for _, number := range numbers {
		result += math.Pow(float64(number)-arithmeticMean, 2)
	}
	result /= float64(len(numbers))
	return math.Sqrt(result)
}

// numbers - not empty sequence
func mean(numbers []int) float64 {
	result := float64(0)
	for _, number := range numbers {
		result += float64(number)
	}
	result /= float64(len(numbers))
	return result
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		meanFlag = true
		medianFlag = true
		modeFlag = true
		sdFlag = true
	}
	scanner := bufio.NewScanner(os.Stdin)
	numbers := make([]int, 0)

	for scanner.Scan() {
		str := scanner.Text()
		number, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		if number < -100000 || number > 100000 {
			log.Fatalf("out of range\n")
		}
		numbers = append(numbers, int(number))
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	if len(numbers) == 0 {
		return
	}
	sort.Ints(numbers)
	if meanFlag {
		fmt.Printf("Mean: %.2f\n", mean(numbers))
	}
	if medianFlag {
		fmt.Printf("Median: %.2f\n", median(numbers))
	}
	if modeFlag {
		fmt.Printf("Mode: %d\n", mode(numbers))
	}
	if sdFlag {
		fmt.Printf("SD: %.2f\n", sd(numbers))
	}
}
