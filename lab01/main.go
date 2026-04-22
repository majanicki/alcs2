package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

func randomSwapIndices(n int) (i, j int) {
	i = random(n)
	j = (i + random(n-1) + 1) % n
	if i > j {
		i, j = j, i
	}
	return
}

func neighborhoodSize(n int) int {
	return n * (n - 1) / 2
}

func estimateAveragePositiveDelta(permutation, A, B []int, n int, sampleCount int) float64 {
	sumPositiveDelta := 0
	countPositive := 0
	for i := 0; i < sampleCount; i++ {
		idx, jdx := randomSwapIndices(n)
		delta := swapDelta(permutation, A, B, n, idx, jdx)
		if delta > 0 {
			sumPositiveDelta += delta
			countPositive++
		}
	}
	if countPositive == 0 {
		return 1.0
	}
	return float64(sumPositiveDelta) / float64(countPositive)
}

type SimulatedAnnealingParams struct {
	LDivider int     // L = neighborhoodSize / LDivider
	P        int     // max no improvement multiplier
	Alpha    float64 // cooling rate
}

func simulatedAnnealing(permutation, A, B []int, n int, _ time.Duration, args any) ([]int, int, int, int) {
	quality := evaluate(permutation, A, B, n)
	bestQuality := quality
	bestPermutation := clonePermutation(permutation)

	iterations := 0
	deltaEvals := 0

	// parameters
	// L might be constant -> let's use a constant number
	// P it is maximum no improvement number multiplayer
	// alpha is cooling rate

	L := neighborhoodSize(n) / args.(SimulatedAnnealingParams).LDivider
	P := args.(SimulatedAnnealingParams).P
	alpha := args.(SimulatedAnnealingParams).Alpha
	// temperature := 1.0
	finalTemperature := 0.01

	tempSampleCount := L

	avgPositiveDelta := estimateAveragePositiveDelta(permutation, A, B, n, tempSampleCount)
	deltaEvals += tempSampleCount

	startAcceptance := 0.95
	temperature := -avgPositiveDelta / math.Log(startAcceptance)
	if temperature < 1e-9 {
		temperature = 1.0
	}

	noImproveCounter := 0

	for noImproveCounter <= P*L && temperature >= finalTemperature {
		for step := 0; step < L; step++ {
			i, j := randomSwapIndices(n)
			delta := swapDelta(permutation, A, B, n, i, j)
			deltaEvals++

			accept := delta <= 0
			if !accept {
				acceptanceProbability := math.Exp(-float64(delta) / temperature)
				accept = rand.Float64() < acceptanceProbability
			}

			if accept {
				swap(permutation, i, j)
				quality += delta
			}

			iterations++
			noImproveCounter++

			if quality < bestQuality {
				bestQuality = quality
				copy(bestPermutation, permutation)
				noImproveCounter = 0
			}
		}

		temperature *= alpha

	}

	return bestPermutation, bestQuality, iterations, deltaEvals
}

type moveCandidate struct {
	i     int
	j     int
	delta int
}

func selectTabuMove(permutation, A, B []int, n, quality, bestQuality, iteration int, tabuExpiry []int, eliteK int) ([]moveCandidate, int, bool) {
	deltaEvals := 0

	candidates := []moveCandidate{}

	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			delta := swapDelta(permutation, A, B, n, i, j)
		deltaEvals++

		candidateQuality := quality + delta
		tabu := iteration < tabuExpiry[i*n+j]
		aspiration := candidateQuality < bestQuality
		if tabu && !aspiration {
			continue
		} else {
			candidates = append(candidates, moveCandidate{i: i, j: j, delta: delta})
		}
		}
	}

	sort.Slice(candidates, func(a, b int) bool {
		return candidates[a].delta < candidates[b].delta
	})

	if eliteK > len(candidates) {
		eliteK = len(candidates)
	}

	if len(candidates) == 0 {
		return nil, 0, false
	}

	return candidates[:eliteK], deltaEvals, true
}

type TabuSearchParams struct {
	maxNoImprovement int // number of best candidates to consider for moves
}

func tabuSearch(permutation, A, B []int, n int, _ time.Duration, args any) ([]int, int, int, int) {
	quality := evaluate(permutation, A, B, n)
	bestQuality := quality
	bestPermutation := clonePermutation(permutation)

	iterations := 0
	deltaEvals := 0

	tabuTenure := n / 4
	
	tabuExpiry := make([]int, n*n)
	eliteK := n / 10
	
	maxNoImprovement := args.(TabuSearchParams).maxNoImprovement
	noImproveCounter := 0

	for noImproveCounter < maxNoImprovement {

		candidates, evals, found := selectTabuMove(permutation, A, B, n, quality, bestQuality, iterations, tabuExpiry, eliteK)
		deltaEvals += evals
		if found {
			for _, candidate := range candidates {

				bestI := candidate.i
				bestJ := candidate.j

				delta := swapDelta(permutation, A, B, n, bestI, bestJ)
		
				deltaEvals += 1

				tabu := iterations < tabuExpiry[bestI*n+bestJ]
				aspiration := quality > bestQuality+delta
				if tabu && !aspiration {
					continue
				}

				swap(permutation, bestI, bestJ)
				quality += delta
				iterations++
				noImproveCounter++

				
				expiry := iterations + tabuTenure
				tabuExpiry[bestI*n+bestJ] = expiry
				tabuExpiry[bestJ*n+bestI] = expiry

				if quality < bestQuality {
					bestQuality = quality
					copy(bestPermutation, permutation)
					noImproveCounter = 0
				}
			}
		}
	}

	return bestPermutation, bestQuality, iterations, deltaEvals
}

type localSearchResult struct {
	Permutation    []int
	Quality        int
	InitialQuality int
	Iterations     int
	DeltaEvals     int
	Runtime        int64
}

type OptimisationFunc func(permutation, A, B []int, n int, maxDuration time.Duration, args any) ([]int, int, int, int)

func benchmarkAlgorithm(f OptimisationFunc, A, B []int, n int, maxDuration time.Duration, args any) localSearchResult {
	inital := constructRandomPermutation(n)
	initalFitness := evaluate(inital, A, B, n)
	start := time.Now()
	permutation, quality, iterations, deltaEvals := f(inital, A, B, n, maxDuration, args)
	elapsed := time.Since(start)
	runtime := int64(elapsed.Nanoseconds())
	return localSearchResult{
		Permutation:    permutation,
		Quality:        quality,
		InitialQuality: initalFitness,
		Iterations:     iterations,
		DeltaEvals:     deltaEvals,
		Runtime:        runtime,
	}
}

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

func heuristicPermutation(permutation, A, B []int, n int, _ time.Duration, _ any) ([]int, int, int, int) {

	permutation = constructInitPermutation(A, B, n)
	quality := evaluate(permutation, A, B, n)

	return permutation, quality, 0, 1
}

func greedyLocalSearch(permutation, A, B []int, n int, _ time.Duration, _ any) ([]int, int, int, int) {
	quality := evaluate(permutation, A, B, n)
	iterations := 0
	deltaEvals := 0

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
			break
		}
	}

	return permutation, quality, iterations, deltaEvals
}

func steepestLocalSearch(permutation, A, B []int, n int, _ time.Duration, _ any) ([]int, int, int, int) {
	quality := evaluate(permutation, A, B, n)
	iterations := 0
	deltaEvals := 0

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
			break
		}

		swap(permutation, bestI, bestJ)
		quality += bestDelta
		iterations++
	}

	return permutation, quality, iterations, deltaEvals
}

func randomWalkLocalSearch(permutation, A, B []int, n int, maxDuration time.Duration, _ any) ([]int, int, int, int) {
	quality := evaluate(permutation, A, B, n)
	bestQuality := quality
	bestSolution := clonePermutation(permutation)
	iterations := 0
	deltaEvals := 0

	deadline := time.Now().Add(maxDuration)
	for time.Now().Before(deadline) {
		i := random(n)
		j := random(n)
		if i == j {
			continue
		}

		delta := swapDelta(permutation, A, B, n, i, j)
		deltaEvals++
		swap(permutation, i, j)
		quality += delta
		if quality < bestQuality {
			copy(bestSolution, permutation)
			bestQuality = quality
		}
		iterations++
	}

	return bestSolution, bestQuality, iterations, deltaEvals

}

func randomSearch(permutation, A, B []int, n int, maxDuration time.Duration, _ any) ([]int, int, int, int) {
	bestPermutation := clonePermutation(permutation)
	bestQuality := evaluate(bestPermutation, A, B, n)
	currentPermutation := clonePermutation(permutation)
	iterations := 0
	deltaEvals := 0

	deadline := time.Now().Add(maxDuration)
	for time.Now().Before(deadline) {

		randomPermutation(currentPermutation)
		currentQuality := evaluate(currentPermutation, A, B, n)
		deltaEvals++

		if currentQuality < bestQuality {
			bestQuality = currentQuality
			copy(bestPermutation, currentPermutation)
		}

		iterations++
	}

	return permutation, bestQuality, iterations, deltaEvals

}

func IntSliceToString(nums []int, sep string) string {
	strNums := make([]string, len(nums))
	for i, v := range nums {
		strNums[i] = strconv.Itoa(v)
	}
	return strings.Join(strNums, sep)
}

func produceResultsRow(filename, name string, results localSearchResult) []string {
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)
	fileWithoutExt := base[:len(base)-len(ext)]
	return []string{fileWithoutExt, name, IntSliceToString(results.Permutation, " "), strconv.Itoa(results.Quality), strconv.Itoa(results.Iterations), strconv.Itoa(results.DeltaEvals), strconv.FormatInt(results.Runtime, 10), strconv.Itoa(results.InitialQuality)}
}

func main() {
	outputFile, err := os.Create("measurements_8.csv")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	files, err := filepath.Glob("./data/*.dat")
	if err != nil {
		panic(err)
	}
		
	for _, file := range files {
		n, A, B := loadData(file)
		params := SimulatedAnnealingParams{
			LDivider: 	2,
			P: 	2*n,
			Alpha:    	0.97,
		}
		paramsT := TabuSearchParams{
			maxNoImprovement: 2*n,
		}
		for i := 0; i < 10; i++ {
			fmt.Println(file, i)

			results := benchmarkAlgorithm(steepestLocalSearch, A, B, n, time.Duration(0), nil)
			timeForRandom := time.Duration(results.Runtime)
			writer.Write(produceResultsRow(file, "steepest", results))
			results = benchmarkAlgorithm(greedyLocalSearch, A, B, n, time.Duration(0), nil)
			writer.Write(produceResultsRow(file, "greedy", results))
			results = benchmarkAlgorithm(randomWalkLocalSearch, A, B, n, timeForRandom, nil)
			writer.Write(produceResultsRow(file, "randomWalk", results))
			results = benchmarkAlgorithm(randomSearch, A, B, n, timeForRandom, nil)
			writer.Write(produceResultsRow(file, "random", results))
			results = benchmarkAlgorithm(heuristicPermutation, A, B, n, time.Duration(0), nil)
			writer.Write(produceResultsRow(file, "heuristic", results))
			results = benchmarkAlgorithm(simulatedAnnealing, A, B, n, time.Duration(0), params)
			writer.Write(produceResultsRow(file, "simulatedAnnealing", results))
			results = benchmarkAlgorithm(tabuSearch, A, B, n, time.Duration(0), paramsT)
			writer.Write(produceResultsRow(file, "tabuSearch", results))

		}
	}
}
