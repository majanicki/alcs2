package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

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

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))


	if len(os.Args) < 2 {
		fmt.Println("Usage: program <data file>")
		return
	}

	n, A, B := loadData(os.Args[1])

	permutation := constructInitPermutation(A, B, n)
	printPermuation(permutation)

}
