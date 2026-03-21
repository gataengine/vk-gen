//go:build nintendosdk

#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <vulkan/vulkan.h>
#include "dlload.h"

static PFN_vkVoidFunction (*getInstanceProcAddress)(VkInstance, const char*) = NULL;
static VkInstance gLastInstance = VK_NULL_HANDLE;

static int str_starts_with(const char* s, const char* prefix) {
	size_t n = strlen(prefix);
	return strncmp(s, prefix, n) == 0;
}

static int is_surface_create_command(const char* procName) {
	return str_starts_with(procName, "vkCreate") && strstr(procName, "Surface") != NULL;
}

static int is_instance_hint_command(const char* procName) {
	return strcmp(procName, "vkEnumeratePhysicalDevices") == 0 ||
		strcmp(procName, "vkEnumeratePhysicalDeviceGroups") == 0 ||
		strcmp(procName, "vkDestroyInstance") == 0 ||
		strcmp(procName, "vkCreateHeadlessSurfaceEXT") == 0 ||
		strcmp(procName, "vkDestroySurfaceKHR") == 0 ||
		is_surface_create_command(procName);
}

static int is_physical_device_command(const char* procName) {
	return str_starts_with(procName, "vkGetPhysicalDevice") ||
		strcmp(procName, "vkEnumerateDeviceExtensionProperties") == 0 ||
		strcmp(procName, "vkEnumerateDeviceLayerProperties") == 0;
}

static int is_global_or_instance_command(const char* procName) {
	return strcmp(procName, "vkCreateInstance") == 0 ||
		strcmp(procName, "vkDestroyInstance") == 0 ||
		strcmp(procName, "vkEnumerateInstanceExtensionProperties") == 0 ||
		strcmp(procName, "vkEnumerateInstanceLayerProperties") == 0 ||
		strcmp(procName, "vkEnumerateInstanceVersion") == 0 ||
		strcmp(procName, "vkGetInstanceProcAddr") == 0 ||
		strcmp(procName, "vkCreateDevice") == 0 ||
		strcmp(procName, "vkEnumeratePhysicalDevices") == 0 ||
		strcmp(procName, "vkEnumeratePhysicalDeviceGroups") == 0 ||
		strcmp(procName, "vkEnumerateDeviceExtensionProperties") == 0 ||
		strcmp(procName, "vkEnumerateDeviceLayerProperties") == 0 ||
		strcmp(procName, "vkCreateHeadlessSurfaceEXT") == 0 ||
		strcmp(procName, "vkDestroySurfaceKHR") == 0 ||
		is_surface_create_command(procName) ||
		str_starts_with(procName, "vkGetPhysicalDevice");
}

void SetProcAddr(void* getProcAddr) {
	getInstanceProcAddress = (PFN_vkVoidFunction (*)(VkInstance, const char*))getProcAddr;
}

void* GetProcFunc(void *instance, const void *name) {
	if (getInstanceProcAddress) {
		return (void*)getInstanceProcAddress((VkInstance)instance, (const char*)name);
	}
	return (void*)vkGetInstanceProcAddr((VkInstance)instance, (const char*)name);
}

void* OpenLibrary(const char *name) {
	(void)name;
	return (void*)0x1;
}

void CloseLibrary(void *lib_handle) {
	(void)lib_handle;
}

void* SymbolFromNameWithHint(void *lib_handle, const void *name, uintptr_t hint) {
	(void)lib_handle;
	const char* procName = (const char*)name;
	PFN_vkVoidFunction fn = vkGetInstanceProcAddr(VK_NULL_HANDLE, procName);
	if (fn) {
		return (void*)fn;
	}

	if (hint != 0 && is_instance_hint_command(procName)) {
		VkInstance inst = (VkInstance)hint;
		gLastInstance = inst;
		fn = vkGetInstanceProcAddr(inst, procName);
		if (fn) {
			return (void*)fn;
		}
	}

	if (is_physical_device_command(procName) && gLastInstance != VK_NULL_HANDLE) {
		fn = vkGetInstanceProcAddr(gLastInstance, procName);
		if (fn) {
			return (void*)fn;
		}
	}

	if (strcmp(procName, "vkCreateDevice") == 0 && gLastInstance != VK_NULL_HANDLE) {
		fn = vkGetInstanceProcAddr(gLastInstance, procName);
		if (fn) {
			return (void*)fn;
		}
	}

	if (hint != 0 && !is_global_or_instance_command(procName)) {
		VkInstance instForLookup = gLastInstance;
		PFN_vkGetDeviceProcAddr getDeviceProcAddr = NULL;
		if (instForLookup != VK_NULL_HANDLE) {
			getDeviceProcAddr = (PFN_vkGetDeviceProcAddr)vkGetInstanceProcAddr(instForLookup, "vkGetDeviceProcAddr");
		}
		if (!getDeviceProcAddr) {
			getDeviceProcAddr = (PFN_vkGetDeviceProcAddr)vkGetInstanceProcAddr(VK_NULL_HANDLE, "vkGetDeviceProcAddr");
		}
		if (getDeviceProcAddr) {
			fn = getDeviceProcAddr((VkDevice)hint, procName);
			if (fn) {
				return (void*)fn;
			}
		}
	}

	return (void*)fn;
}

void* SymbolFromName(void *lib_handle, const void *name) {
	return SymbolFromNameWithHint(lib_handle, name, 0);
}
