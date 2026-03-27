package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type SolutionData struct {
	Filename    string
	N           int
	Objective   int
	Permutation []int
	A		   []int
	B		   []int
}

var dataFolder = "data"

func importData() ([]SolutionData, error) {
	var solutions []SolutionData
	

	// Read all files in the data folder
	files, err := os.ReadDir(dataFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read data folder: %w", err)
	}

	// Process each .sln file
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sln") {
			continue
		}

		filepath := filepath.Join(dataFolder, file.Name())
		data, err := parseSolutionFile(filepath)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file.Name(), err)
			continue
		}

		solutions = append(solutions, data)
	}

	return solutions, nil
}

func parseSolutionFile(filename string) (SolutionData, error) {
	data := SolutionData{Filename: filename}

	file, err := os.Open(filename)
	if err != nil {
		return data, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if lineNum == 0 {
			// Parse first line: n and objective value
			parts := strings.Fields(line)
			if len(parts) < 2 {
				return data, fmt.Errorf("invalid format in first line")
			}

			n, err := strconv.Atoi(parts[0])
			if err != nil {
				return data, fmt.Errorf("failed to parse n: %w", err)
			}

			obj, err := strconv.Atoi(parts[1])
			if err != nil {
				return data, fmt.Errorf("failed to parse objective: %w", err)
			}

			data.N = n
			data.Objective = obj
			data.Permutation = make([]int, 0, n)
		} else {
			// Parse permutation lines
			parts := strings.Fields(line)
			for _, part := range parts {
				val, err := strconv.Atoi(part)
				if err != nil {
					return data, fmt.Errorf("failed to parse permutation value: %w", err)
				}
				data.Permutation = append(data.Permutation, val)
			}
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return data, fmt.Errorf("scanner error: %w", err)
	}

	if len(data.Permutation) != data.N {
		return data, fmt.Errorf("permutation size mismatch: expected %d, got %d", data.N, len(data.Permutation))
	}

	for i := 0; i < data.N; i++ {
		data.Permutation[i]-- // Convert to 0-based index
	}

	return data, nil
}

func loadMatrixData(solutions []SolutionData)  {


	// Read all files in the data folder
	files, err := os.ReadDir(dataFolder)
	if err != nil {
		// return nil, fmt.Errorf("failed to read data folder: %w", err)
	}

	// Process each .sln file
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".dat") {
			continue
		}

		filepath := filepath.Join(dataFolder, file.Name())
		n, A, B := loadData(filepath)
		for i := range solutions {
			if solutions[i].N == n && filepath == solutions[i].Filename[:len(solutions[i].Filename)-4]+".dat" {
				solutions[i].A = A
				solutions[i].B = B
				break
			}

		}
	}
}


func TestEvaluate(t *testing.T) {
	solutions, err := importData()
	loadMatrixData(solutions)

	for i, sol := range solutions {
		if len(sol.A) != sol.N*sol.N || len(sol.B) != sol.N*sol.N || len(sol.Permutation) != sol.N {
			t.Errorf("Solution %d: matrix size mismatch for N=%d", i, sol.N)
		}
	}

	if err != nil {
		t.Fatalf("Failed to import data: %v", err)
	}

	if len(solutions) == 0 {
		t.Fatal("No solution files were loaded")
	}

	for i, sol := range solutions {
		eval := evaluate(sol.Permutation, sol.A, sol.B, sol.N)
		// fmt.Println("iteration", sol.Filename, "eval", eval, "objective", sol.Objective)
		if eval != sol.Objective {
			t.Errorf("Solution %d, %s: expected objective %d, got %d", i, sol.Filename, sol.Objective, eval)
		}
	}
}

func TestSwapDelta(t *testing.T) {
	solutions, err := importData()
	loadMatrixData(solutions)

	for i, sol := range solutions {
		if len(sol.A) != sol.N*sol.N || len(sol.B) != sol.N*sol.N || len(sol.Permutation) != sol.N {
			t.Errorf("Solution %d: matrix size mismatch for N=%d", i, sol.N)
		}
	}

	if err != nil {
		t.Fatalf("Failed to import data: %v", err)
	}

	if len(solutions) == 0 {
		t.Fatal("No solution files were loaded")
	}

	for i, sol := range solutions {
		eval := evaluate(sol.Permutation, sol.A, sol.B, sol.N)
		for x := 0; x < sol.N; x++ {
			for y := 0; y < sol.N; y++ {
				delta := swapDelta(sol.Permutation, sol.A, sol.B, sol.N, x, y)
				// Create a copy of the permutation and swap x and y
				permutationCopy := clonePermutation(sol.Permutation)
				swap(permutationCopy, x, y)
				evalAfterSwap := evaluate(permutationCopy, sol.A, sol.B, sol.N)
				expectedDelta := evalAfterSwap - eval
				if delta != expectedDelta {
					t.Errorf("Solution %d, %s: swapDelta mismatch for swap (%d, %d): expected %d, got %d", i, sol.Filename, x, y, expectedDelta, delta)
				}
			}
		
		}
	}
}
