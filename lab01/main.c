#include <stdlib.h>
#include <stdio.h>
#include <time.h>

typedef int** Matrix;

void print_matrix(Matrix matrix, int n) {
    for(int i = 0; i < n; i++) {
        for(int j = 0; j < n; j++) {
            printf("%5d ", matrix[i][j]);
        }
        printf("\n");
    }
}
static inline int random(int n) {
    return rand() % n;
}

void random_permutation(int *permutation, int n) {
    for(int i = 0; i < n; i++) {
        int r = random(n - i);
        int swap = permutation[n - i - 1];
        permutation[n - i - 1] = permutation[r];
        permutation[r] = swap;
    }
}

void random_pair(int n, int *a, int *b) {
    *a = random(n);
    *b = (*a + random(n-1) + 1) % n;
}
/*
int measure_ns(int threshold) {
    int time_0 = time();
    int c = 0;
    do {
        f();
        c++;
    } while((time() - time_0) < threshold);
    return (time() - time_0) / e;
}

*/
Matrix read_matrix(FILE* file, int n) {
    int res;
    // preserve some data locality
    int *matrix_mem = malloc(sizeof(int) * n * n);
    int **matrix = malloc(sizeof(int*) * n);
    for(int i = 0; i < n; i++) {
        matrix[i] = matrix_mem + i * n;
    }
    for(int i = 0; i < n; i++) {
        for(int j = 0; j < n - 1; j++) {
            res = fscanf(file, "%d ", &matrix[i][j]);
            if(res != 1) {
                fprintf(stderr, "Failed to read matrix\n");
                exit(1);
            }
        }
        res = fscanf(file, "%d\n", &matrix[i][n-1]);
        if(res != 1) {
            fprintf(stderr, "Failed to read matrix\n");
            exit(1);
        }
    }
    fscanf(file, "\n");
    return matrix;
}

void read_instance(char* filename, int *rn, Matrix *rA, Matrix *rB) {
    int res;
    FILE *file = fopen(filename,"r");
    if(file == NULL) {
        fprintf(stderr, "Failed to open file %s\n", filename);
        exit(1);
    }
    int n;
    res = fscanf(file, "%d\n\n", &n);
    if(res != 1) {
        fprintf(stderr, "Failed to read N from file %s\n", filename);
        exit(1);
    }
    printf("N = %d\n", n);
    printf("Matrix A\n");
    Matrix A = read_matrix(file, n);
    print_matrix(A, n);
    printf("\n");
    printf("Matrix B\n");
    Matrix B = read_matrix(file, n);
    print_matrix(B, n);
    *rn = n;
    *rA = A;
    *rB = B;
}

int evaluate_solution(int *p, int n, Matrix A, Matrix B) {
    int cost = 0;
    for(int i = 0; i < n; i++) {
        for(int j = 0; j < n; j++) {
            cost += A[i][j] * B[p[i]][p[j]];
        }
    }
    return cost;
}

int main (int argc, char**argv) {
    if(argc != 2) {
        fprintf(stderr, "Usage: %s [DATA_FILE]\n", argv[0]);
        return 1;
    }
    Matrix A, B;
    int n;
    read_instance(argv[1], &n, &A, &B);
    int test_soluton[] = {
        10,
        14,
        25,
        6,
        3,
        12,
        11,
        1,
        5,
        17,
        8,
        4,
        0,
        20,
        7,
        13,
        2,
        19,
        18,
        9,
        16,
        24,
        15,
        23,
        21,
        22,
    };
    printf("Solution cost: %d\n", evaluate_solution(test_soluton, n, A, B));
    return 0;
}
