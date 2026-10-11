#ifndef __DLLOAD_H__
#define __DLLOAD_H__

#include <stddef.h>
#include <stdint.h>

void* OpenLibrary(const char *name);
void CloseLibrary(void *lib_handle);
void SetProcAddr(void* getProcAddr);
void* GetProcFunc(void *dispatchable, const void *name);
void* SymbolFromNameWithHint(void *lib_handle, const void *name, uintptr_t hint);

size_t Trampoline3(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2);
size_t Trampoline6(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5);
size_t Trampoline9(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8);
size_t Trampoline12(void *symbol, uintptr_t p0, uintptr_t p1, uintptr_t p2, uintptr_t p3, uintptr_t p4, uintptr_t p5, uintptr_t p6, uintptr_t p7, uintptr_t p8, uintptr_t p9, uintptr_t p10, uintptr_t p11);

// TrampolineFloat* call commands taking integer/handle parameters followed only by float
// parameters (vkCmdSetDepthBias, vkSetDeviceMemoryPriorityEXT, ...). Floats travel in FP
// registers (by position on Windows x64), so each signature needs its exact prototype.
size_t TrampolineFloat1P1F(void *symbol, uintptr_t p0, float f0);
size_t TrampolineFloat1P2F(void *symbol, uintptr_t p0, float f0, float f1);
size_t TrampolineFloat1P3F(void *symbol, uintptr_t p0, float f0, float f1, float f2);
size_t TrampolineFloat2P1F(void *symbol, uintptr_t p0, uintptr_t p1, float f0);

void* SymbolFromName(void *lib_handle, const void *name);

#endif
