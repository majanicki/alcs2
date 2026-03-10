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

func main() {
	rand.Seed(time.Now().UnixNano())
	perm := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 1000; i++ {
		randomPermutation(perm)
		printPermuation(perm)
	}
	// randomPermutation(perm)
	// printPermuation(perm)
	// a, b := randomPair(len(perm))
	// fmt.Println(a, b)
}
