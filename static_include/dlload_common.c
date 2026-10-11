#include <stdint.h>
#include "dlload.h"

typedef size_t (*vkGeneric_func3)(uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func6)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkGeneric_func9)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);
typedef size_t (*vkFloat_func1p1f)(uintptr_t, float);
typedef size_t (*vkFloat_func1p2f)(uintptr_t, float, float);
typedef size_t (*vkFloat_func1p3f)(uintptr_t, float, float, float);
typedef size_t (*vkFloat_func2p1f)(uintptr_t, uintptr_t, float);
typedef size_t (*vkGeneric_func12)(uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t, uintptr_t);

size_t Trampoline3(void *symbol, size_t p0, size_t p1, size_t p2) {
	return ((vkGeneric_func3) symbol)(p0, p1, p2);
}

size_t Trampoline6(void *symbol, size_t p0, size_t p1, size_t p2, size_t p3, size_t p4, size_t p5) {
	return ((vkGeneric_func6) symbol)(p0, p1, p2, p3, p4, p5);
}

size_t Trampoline9(void *symbol, size_t p0, size_t p1, size_t p2, size_t p3, size_t p4, size_t p5, size_t p6, size_t p7, size_t p8) {
	return ((vkGeneric_func9) symbol)(p0, p1, p2, p3, p4, p5, p6, p7, p8);
}

size_t Trampoline12(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11) {
	return ((vkGeneric_func12) symbol)(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11);
}

size_t TrampolineFloat1P1F(void *symbol, uintptr_t p0, float f0) {
	return ((vkFloat_func1p1f) symbol)(p0, f0);
}

size_t TrampolineFloat1P2F(void *symbol, uintptr_t p0, float f0, float f1) {
	return ((vkFloat_func1p2f) symbol)(p0, f0, f1);
}

size_t TrampolineFloat1P3F(void *symbol, uintptr_t p0, float f0, float f1, float f2) {
	return ((vkFloat_func1p3f) symbol)(p0, f0, f1, f2);
}

size_t TrampolineFloat2P1F(void *symbol, uintptr_t p0, uintptr_t p1, float f0) {
	return ((vkFloat_func2p1f) symbol)(p0, p1, f0);
}
