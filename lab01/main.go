package main
import (
	"fmt"
	"math/rand"
	"time"
)

func random(n int) (int) {
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
	b = (a + random(n - 1) + 1) % n
	return
}

func printPermuation(permutation []int) {
	for _, e := range permutation {
		fmt.Print(e, " ")
	}
	fmt.Println("")
}

// Think about whether we want to make this more inline with what we discussed at labs
func measureTime(f func()) (float64) {
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


func main() {
	rand.Seed(time.Now().UnixNano())
	// perm := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// randomPermutation(perm)
	// printPermuation(perm)
	// a, b := randomPair(len(perm))
	// fmt.Println(a, b)
	// elapsed := measureTime(func() {time.Sleep(2 * time.Millisecond)})
	// fmt.Println(elapsed)
	// 2-OPT neighbourhoud
	n := 5
	// for randomizing the neighbourhood we need to start from random point
	offset := random(n) // this offset is optional (should check if this correct)
	for i := 0; i < n - 1; i++ { // skip unnecessary last iteration
		for j := i + 1; j < n; j++ {
			// calculate deltas
			fmt.Println((i + offset) % n, (j + offset) % n)
		}
	}
	// do the swap with xor if you want to be fast
	// or do the regular swap but only for accepted neighbour

	// outputing results:
	// - quality, permutation, running time, how many jumps from neighbour to neighbours
	//    how many times we evaluated delta
	// for next week - loading data, heuristic, 2-opt neighbours
}
