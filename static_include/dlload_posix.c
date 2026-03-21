//go:build !windows && !nintendosdk

#include <dlfcn.h>
#include <stdint.h>
#include <stdio.h>
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
	void* lib_handle = dlopen(name, RTLD_LOCAL | RTLD_LAZY);
	if (!lib_handle) {
		printf("Unable to load %s: %s\n", name, dlerror());
		return NULL;
	}
	return lib_handle;
}

void CloseLibrary(void *lib_handle) {
	if (dlclose(lib_handle) != 0) {
		printf("Problem closing library: %s", dlerror());
	}
}

void* SymbolFromName(void *lib_handle, const void *name) {
	return dlsym(lib_handle, (const char*)name);
}

void* SymbolFromNameWithHint(void *lib_handle, const void *name, uintptr_t hint) {
	(void)hint;
	return SymbolFromName(lib_handle, name);
}
