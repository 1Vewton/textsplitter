#include <stdint.h>
#include <immintrin.h>

// add_avx2 adds the i64
void add_avx2(int64_t* a, int64_t* b, int64_t* out, int n) {
    int i = 0;
    for (; i < n; i += 4) {
        __m256i va = _mm256_loadu_si256((__m256i*)&a[i]);
        __m256i vb = _mm256_loadu_si256((__m256i*)&b[i]);
        __m256i sum = _mm256_add_epi64(va, vb);
        _mm256_storeu_si256((__m256i*)&out[i], sum);
    }
    // Process the rest of the numbers
    for (; i < n; i++) {
        out[i] = a[i] + b[i];
    }
}

// add_avx2_double adds the f64
void add_avx2_double(double* a, double* b, double* out, int n) {
    int i = 0;
    for (; i + 4 <= n; i += 4) {
        __m256d va = _mm256_loadu_pd(&a[i]);
        __m256d vb = _mm256_loadu_pd(&b[i]);
        __m256d sum = _mm256_add_pd(va, vb);
        _mm256_storeu_pd(&out[i], sum);
    }
    for (; i < n; i++) {
        out[i] = a[i] + b[i];
    }
}

// mult_avx2_double adds the f64
void mult_avx2_double(double* a, double* b, double* out, int n) {
    int i = 0;
    for (; i + 4 <= n; i += 4) {
        __m256d va = _mm256_loadu_pd(&a[i]);
        __m256d vb = _mm256_loadu_pd(&b[i]);
        __m256d result = _mm256_mul_pd(va, vb);
        _mm256_storeu_pd(&out[i], result);
    }
    for (; i < n; i++) {
        out[i] = a[i] + b[i];
    }
}

// Calculate the sum of the avx2
double sum_avx2_double(double* out, int n) {
    int i = 0;
    __m256d acc = _mm256_setzero_pd();
    for (; i + 4 <= n; i += 4) {
        __m256d v = _mm256_loadu_pd(&out[i]);
        acc = _mm256_add_pd(acc, v);
    }
    __m128d lo = _mm256_castpd256_pd128(acc);
    __m128d hi = _mm256_extractf128_pd(acc, 1);
    __m128d sum128 = _mm_add_pd(lo, hi);
    __m128d sum64 = _mm_hadd_pd(sum128, sum128);
    double result = _mm_cvtsd_f64(sum64);
    for (; i < n; i++) {
        result += out[i];
    }
    return result;
}
