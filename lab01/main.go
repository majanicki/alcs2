package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type localSearchResult struct {
	permutation []int
	quality     int
	iterations  int
	deltaEvals  int
	runtime     float64
}

func random(n int) int {
	return rand.Intn(n)
}

func randomPermutation(permutation []int) {
	n := len(permutation)
	for i := 0; i < n; i++ {
		r := random(n - i)
		pos := n - i - 1
		swap := permutation[pos]
		permutation[pos] = permutation[r]
		permutation[r] = swap
	}
}

func constructRandomPermutation(n int) (permutation []int) {
	permutation = make([]int, n)
	for i := 0; i < n; i++ {
		permutation[i] = i
	}
	randomPermutation(permutation)
	return
}

func randomPair(n int) (a, b int) {
	a = random(n)
	b = (a + random(n-1) + 1) % n
	return
}

func printPermuation(permutation []int) {
	for _, e := range permutation {
		fmt.Print(e, " ")
	}
	fmt.Println("")
}

func swap(permutation []int, i, j int) {
	permutation[i], permutation[j] = permutation[j], permutation[i]
}

func clonePermutation(permutation []int) []int {
	copyPermutation := make([]int, len(permutation))
	copy(copyPermutation, permutation)
	return copyPermutation
}

func evaluate(permutation, A, B []int, n int) int {
	quality := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			quality += A[i*n+j] * B[permutation[i]*n+permutation[j]]
		}
	}
	return quality
}

func swapDelta(permutation, A, B []int, n, i, j int) int {
	if i == j {
		return 0
	}

	pi := permutation[i]
	pj := permutation[j]

	delta := (A[i*n+i] - A[j*n+j]) * (B[pj*n+pj] - B[pi*n+pi])
	delta += (A[i*n+j] - A[j*n+i]) * (B[pj*n+pi] - B[pi*n+pj])

	for k := 0; k < n; k++ {
		if k == i || k == j {
			continue
		}

		pk := permutation[k]
		delta += (A[k*n+i] - A[k*n+j]) * (B[pk*n+pj] - B[pk*n+pi])
		delta += (A[i*n+k] - A[j*n+k]) * (B[pj*n+pk] - B[pi*n+pk])
	}

	return delta
}

// Think about whether we want to make this more inline with what we discussed at labs
func measureTime(f func()) float64 {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	return elapsed.Seconds()
}

// We should read data with `make` allocations

// Heuristic:
// Prof: Average rows and columns and match lowest average with highest average
// Other: Pick closest with highest flow - this one is okay

// for deltas we are changing two rows and two columns
// deltas should be calculated in linear time

//// 2-OPT neighbourhoud
//// for randomizing the neighbourhood we need to start from random point
//offset := random(n) // this offset is optional (should check if this correct)
//for i := 0; i < n - 1; i++ { // skip unnecessary last iteration
//	for j := i + 1; j < n; j++ {
//		// calculate deltas
//		fmt.Println((i + offset) % n, (j + offset) % n)
//	}
//}
//// do the swap with xor if you want to be fast
//// or do the regular swap but only for accepted neighbour
//
//// outputing results:
//// - quality, permutation, running time, how many jumps from neighbour to neighbours
////    how many times we evaluated delta
//// for next week - loading data, heuristic, 2-opt neighbours

// row major format

func loadMatrix(file *os.File, n int) (A []int) {
	A = make([]int, n*n)

	for i := 0; i < n*n; i++ {
		if _, err := fmt.Fscan(file, &A[i]); err != nil {
			panic(err)
		}
	}
	return
}

func loadData(filename string) (n int, A, B []int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = fmt.Fscan(file, &n); err != nil {
		panic(err)
	}
	A = loadMatrix(file, n)
	B = loadMatrix(file, n)
	return
}

func constructInitPermutation(A, B []int, n int) (permutation []int) {
	permutation = make([]int, n)

	locationScore := make([]int, n)
	objectScore := make([]int, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			locationScore[i] += A[i*n+j] + A[j*n+i]
			objectScore[i] += B[i*n+j] + B[j*n+i]
		}
	}

	locationOrder := make([]int, n)
	objectOrder := make([]int, n)
	for i := 0; i < n; i++ {
		locationOrder[i] = i
		objectOrder[i] = i
	}

	sort.Slice(locationOrder, func(i, j int) bool {
		return locationScore[locationOrder[i]] < locationScore[locationOrder[j]]
	})
	sort.Slice(objectOrder, func(i, j int) bool {
		return objectScore[objectOrder[i]] > objectScore[objectOrder[j]]
	})

	for k := 0; k < n; k++ {
		permutation[locationOrder[k]] = objectOrder[k]
	}

	return
}

func greedyLocalSearch(permutation, A, B []int, n int) localSearchResult {
	quality := evaluate(permutation, A, B, n)
	iterations := 0
	deltaEvals := 0

	runtime := measureTime(func() {
		for {
			improved := false
			offset := random(n) // randomize the starting point in the neighborhood

			for step := 0; step < n-1 && !improved; step++ {
				i := (offset + step) % n
				for shift := 1; shift < n; shift++ {
					j := (i + shift) % n
					if i >= j {
						continue
					}

					delta := swapDelta(permutation, A, B, n, i, j)
					deltaEvals++
					if delta < 0 {
						swap(permutation, i, j)
						quality += delta
						iterations++
						improved = true
						break
					}
				}
			}

			if !improved {
				return
			}
		}
	})

	return localSearchResult{
		permutation: permutation,
		quality:     quality,
		iterations:  iterations,
		deltaEvals:  deltaEvals,
		runtime:     runtime,
	}
}

func steepestLocalSearch(permutation, A, B []int, n int) localSearchResult {
	quality := evaluate(permutation, A, B, n)
	iterations := 0
	deltaEvals := 0

	runtime := measureTime(func() {
		for {
			bestDelta := 0
			bestI := -1
			bestJ := -1

			for i := 0; i < n-1; i++ {
				for j := i + 1; j < n; j++ {
					delta := swapDelta(permutation, A, B, n, i, j)
					deltaEvals++
					if delta < bestDelta {
						bestDelta = delta
						bestI = i
						bestJ = j
					}
				}
			}

			if bestDelta >= 0 {
				return
			}

			swap(permutation, bestI, bestJ)
			quality += bestDelta
			iterations++
		}
	})

	return localSearchResult{
		permutation: permutation,
		quality:     quality,
		iterations:  iterations,
		deltaEvals:  deltaEvals,
		runtime:     runtime,
	}
}

func randomWalkLocalSearch(permutation, A, B []int, n int, maxIterations int) localSearchResult {
	quality := evaluate(permutation, A, B, n)
	iterations := 0
	deltaEvals := 0

	runtime := measureTime(func() {
		for iterations < maxIterations {
			i := random(n)
			j := random(n)
			if i == j {
				continue
			}

			delta := swapDelta(permutation, A, B, n, i, j)
			deltaEvals++

			swap(permutation, i, j)
			quality += delta
			iterations++
		}
	})

	return localSearchResult{
		permutation: permutation,
		quality:     quality,
		iterations:  iterations,
		deltaEvals:  deltaEvals,
		runtime:     runtime,
	}
}

func randomSearch(permutation, A, B []int, n int, maxIterations int) localSearchResult {
	bestPermutation := clonePermutation(permutation)
	bestQuality := evaluate(bestPermutation, A, B, n)
	currentPermutation := clonePermutation(permutation)
	iterations := 0
	deltaEvals := 0

	runtime := measureTime(func() {
		for iterations < maxIterations {

			randomPermutation(currentPermutation)
			currentQuality := evaluate(currentPermutation, A, B, n)
			deltaEvals++


			if currentQuality < bestQuality {
				bestQuality = currentQuality
				copy(bestPermutation, currentPermutation)
			}

			iterations++
		}
	})

	return localSearchResult{
		permutation: bestPermutation,
		quality:     bestQuality,
		iterations:  iterations,
		deltaEvals:  deltaEvals,
		runtime:     runtime,
	}
}

func appendResultToCSV(filename string, result localSearchResult, algorithm string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "%s,%d,%d,%d,%v\n", algorithm, result.quality, result.iterations, result.deltaEvals, result.runtime)
	if err != nil {
		return fmt.Errorf("failed to write result: %w", err)
	}
	return nil
}

func saveResultCSV(filename string, results []localSearchResult, algorithm string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Algorithm,Quality,Iterations,DeltaEvals,Runtime\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, result := range results {
		_, err := fmt.Fprintf(file, "%s,%d,%d,%d,%v\n", algorithm, result.quality, result.iterations, result.deltaEvals, result.runtime)
		if err != nil {
			return fmt.Errorf("failed to write result: %w", err)
		}
	}

	return nil
}


func printResult(name string, result localSearchResult) {
	fmt.Println(name)
	fmt.Println("quality:", result.quality)
	fmt.Println("permutation:")
	printPermuation(result.permutation)
	fmt.Println("iterations:", result.iterations)
	fmt.Println("delta evaluations:", result.deltaEvals)
	fmt.Println("time:", result.runtime)
	fmt.Println("")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("Usage: program <data file>")
		return
	}

	n, A, B := loadData(os.Args[1])

	randomPermutation := constructRandomPermutation(n)

	greedyResult := greedyLocalSearch(clonePermutation(randomPermutation), A, B, n)
	steepestResult := steepestLocalSearch(clonePermutation(randomPermutation), A, B, n)

	appendResultToCSV(".\\result\\result.csv", greedyResult, "greedy/" + os.Args[1])
	appendResultToCSV(".\\result\\result.csv", steepestResult, "steepest/" + os.Args[1])


}
