#include <stdint.h>
#include <immintrin.h>
#ifndef calculator
#define calculator
void add_avx2(int64_t* a, int64_t* b, int64_t* out, int n);
void add_avx2_double(double* a, double* b, double* out, int n);
void mult_avx2_double(double* a, double* b, double* out, int n);
double sum_avx2_double(double* out, int n);
#endif