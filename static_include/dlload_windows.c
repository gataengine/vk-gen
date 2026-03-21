//go:build windows

#include <stdint.h>
#include <windows.h>
#include "dlload.h"

static void* (*getInstanceProcAddress)(void*, const char*) = NULL;

void SetProcAddr(void* getProcAddr) {
	getInstanceProcAddress = getProcAddr;
}

void* GetProcFunc(void *dispatchable, const void *name) {
	if (getInstanceProcAddress == NULL) {
		return NULL;
	}
	return (*getInstanceProcAddress)(dispatchable, (const char*)name);
}

void* OpenLibrary(const char *name) {
	return LoadLibrary(name);
}

void CloseLibrary(void *lib_handle) {
	FreeLibrary(lib_handle);
}

void* SymbolFromName(void *lib_handle, const void *name) {
	return GetProcAddress(lib_handle, name);
}

void* SymbolFromNameWithHint(void *lib_handle, const void *name, uintptr_t hint) {
	(void)hint;
	return SymbolFromName(lib_handle, name);
}
