#include <immintrin.h>

// add_avx2 adds the i64
void add_avx2(int64_t* a, int64_t* b, int64_t* out, int n) {
    for (int i = 0; i < n; i += 4) {
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
