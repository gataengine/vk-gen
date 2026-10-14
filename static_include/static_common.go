package vk

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

// #include <stdlib.h>
// #include "dlload.h"
// #cgo noescape SymbolFromNameWithHint
// #cgo nocallback SymbolFromNameWithHint
// #cgo nocallback OpenLibrary
// #cgo noescape Trampoline3
// #cgo noescape Trampoline6
// #cgo noescape Trampoline9
// #cgo noescape Trampoline12
// #cgo nocallback Trampoline3
// #cgo nocallback Trampoline6
// #cgo nocallback Trampoline9
// #cgo nocallback Trampoline12
// #cgo nocallback TrampolineFloat1P1F
// #cgo nocallback TrampolineFloat1P2F
// #cgo nocallback TrampolineFloat1P3F
// #cgo nocallback TrampolineFloat2P1F
import "C"

// Vulkanizer allows conversion from go-vk style structs to Vulkan-native structs. This
// includes setting the structure type flag, converting slices to pointers, etc.
type Vulkanizer interface {
	Vulkanize() unsafe.Pointer
}

// Goifier converts Vulkan-native structs back into go-vk style structs
type Goifier interface {
	Goify() Vulkanizer
}

// max is an internal utility function, used in processing struct member slice/array lengths
func max(nums ...int) int {
	rval := 0
	for _, v := range nums {
		if rval < v {
			rval = v
		}
	}
	return rval
}

// Error implements the error interface
// TODO: A way for commands to indicate if the Result code is an error for that command, or an unexpected return value?
func (r Result) Error() string {
	return r.String()
}

type vkCommand struct {
	protoName        string
	argCount         int
	hasReturn        bool
	dispatchArgIndex int
	fnHandle         unsafe.Pointer
}

var dlHandle unsafe.Pointer
var procAddressEnabled bool
var allCommands []*vkCommand

var overrideLibName string
var loadedLibName string

// LoadedVulkanLibrary returns the name or path go-vk opened, or "" before the first
// Vulkan call (or when SetProcAddress supplies the entry points).
func LoadedVulkanLibrary() string {
	return loadedLibName
}

// OverrideDefaultVulkanLibrary allows you to set a specific Vulkan library name to be used in your program. For
// example, if you want to enable the validation layers, those layers are only available in the Vulkan SDK libary. go-vk
// passes the name to the host operating system's library opening/search method, so you must provide a relative or
// absolute path if your Vulkan library is not in the default search path for the platform.
func OverrideDefaultVulkanLibrary(nameOrPath string) {
	overrideLibName = nameOrPath
}

func newVkCommand(protoName string, argCount int, hasReturn bool, dispatchArgIndex int) *vkCommand {
	cmd := &vkCommand{
		protoName:        protoName,
		argCount:         argCount,
		hasReturn:        hasReturn,
		dispatchArgIndex: dispatchArgIndex,
	}
	allCommands = append(allCommands, cmd)
	return cmd
}

func resetCommandHandles() {
	for _, cmd := range allCommands {
		cmd.fnHandle = nil
	}
}

func SetProcAddress(addr unsafe.Pointer) {
	C.SetProcAddr(addr)
	procAddressEnabled = addr != nil
	resetCommandHandles()
}

func resolveCommandHandle(cmd *vkCommand, args []uintptr) unsafe.Pointer {
	var hint C.uintptr_t
	if len(args) > 0 {
		hint = C.uintptr_t(args[0])
	}

	if !procAddressEnabled {
		return C.SymbolFromNameWithHint(dlHandle, unsafe.Pointer(sys_stringToBytePointer(cmd.protoName)), hint)
	}

	var dispatchable unsafe.Pointer
	if cmd.dispatchArgIndex >= 0 && cmd.dispatchArgIndex < len(args) && args[cmd.dispatchArgIndex] != 0 {
		dispatchable = unsafe.Pointer(args[cmd.dispatchArgIndex])
	}

	return C.GetProcFunc(dispatchable, unsafe.Pointer(sys_stringToBytePointer(cmd.protoName)))
}

// execTrampoline is marked //go:uintptrescapes because bindings pass Go pointers
// (notably output parameters such as &instance) as uintptr(unsafe.Pointer(p)). A plain
// uintptr is invisible to the runtime: if the goroutine stack grows and is copied before
// the cgo call (resolveCommandHandle, the prologue here, ...), the pointer is not
// adjusted and Vulkan writes the result into the stale stack, leaving the caller with a
// zero handle. The directive makes every such conversion at a call site heap-allocate
// the pointee and keep it alive for the duration of the call, as syscall.SyscallN does.
//
//go:uintptrescapes
func execTrampoline(cmd *vkCommand, args ...uintptr) uintptr {
	if len(args) != cmd.argCount {
		panic("Wrong number of arguments passed for cmd " + cmd.protoName)
	}
	prepareCommand(cmd, args)

	var result C.uintptr_t

	switch cmd.argCount {
	case 1:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), 0, 0)
	case 2:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), 0)
	case 3:
		result = C.Trampoline3(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]))
	case 4:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), 0, 0)
	case 5:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), 0)
	case 6:
		result = C.Trampoline6(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]))
	case 7:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), 0, 0)
	case 8:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), 0)
	case 9:
		result = C.Trampoline9(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]))
	case 10:
		result = C.Trampoline12(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]), C.uintptr_t(args[9]), 0, 0)
	case 11:
		result = C.Trampoline12(cmd.fnHandle, C.uintptr_t(args[0]), C.uintptr_t(args[1]), C.uintptr_t(args[2]), C.uintptr_t(args[3]), C.uintptr_t(args[4]), C.uintptr_t(args[5]), C.uintptr_t(args[6]), C.uintptr_t(args[7]), C.uintptr_t(args[8]), C.uintptr_t(args[9]), C.uintptr_t(args[10]), 0)
	default:
		// There are no commands with 0 or 12+ arguments as of Vulkan 1.3.204
		panic("Unhandled number of arguments passed for cmd " + cmd.protoName)
	}

	return uintptr(result)
}

// defaultVulkanLibraries lists the libraries tried, in order, when no override is set.
// On macOS the Vulkan loader comes first (it provides the validation layers), then
// MoltenVK alone; each is looked up in $VULKAN_SDK/lib, /usr/local/lib (the SDK's system
// install), /opt/homebrew/lib and finally by bare name. The explicit directories are
// needed because dyld's default fallback path for Go binaries is only /usr/lib.
func defaultVulkanLibraries() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"vulkan-1.dll"}
	case "darwin":
		dirs := []string{"/usr/local/lib", "/opt/homebrew/lib", ""}
		if sdk := os.Getenv("VULKAN_SDK"); sdk != "" {
			dirs = append([]string{filepath.Join(sdk, "lib")}, dirs...)
		}
		var names []string
		for _, lib := range []string{"libvulkan.1.dylib", "libMoltenVK.dylib"} {
			for _, dir := range dirs {
				names = append(names, filepath.Join(dir, lib))
			}
		}
		return names
	case "linux":
		// The runtime soname; plain libvulkan.so exists only with dev packages.
		return []string{"libvulkan.so.1"}
	default:
		panic("Unsupported GOOS at OpenLibrary: " + runtime.GOOS)
	}
}

// prepareCommand loads the Vulkan library on first use and resolves cmd's entry point.
func prepareCommand(cmd *vkCommand, args []uintptr) {
	if dlHandle == nil && !procAddressEnabled && !disableDefaultVulkanLibraryFallback {
		libNames := defaultVulkanLibraries()
		if overrideLibName != "" {
			libNames = []string{overrideLibName}
		}

		for _, libName := range libNames {
			cstr := C.CString(libName)
			dlHandle = C.OpenLibrary(cstr)
			C.free(unsafe.Pointer(cstr))
			if dlHandle != nil {
				loadedLibName = libName
				break
			}
		}
		if dlHandle == nil {
			// Fail here: a nil handle otherwise surfaces as a misleading
			// "Unable to resolve command" for the first command called.
			panic("Unable to load Vulkan library, tried " + fmt.Sprint(libNames))
		}
	}

	if cmd.fnHandle == nil {
		cmd.fnHandle = resolveCommandHandle(cmd, args)
		if cmd.fnHandle == nil {
			panic("Unable to resolve command " + cmd.protoName)
		}
	}
}

// execTrampolineFloat calls a command whose parameters are up to two integers or handles
// followed only by up to three floats (vkCmdSetDepthBias, vkSetDeviceMemoryPriorityEXT,
// ...). The C ABI passes floats in floating-point registers, so routing them through
// execTrampoline as uintptr would make the callee read stale register contents instead
// of the caller's values.
func execTrampolineFloat(cmd *vkCommand, ints []uintptr, floats ...float32) uintptr {
	if len(ints)+len(floats) != cmd.argCount || len(ints) > 2 || len(floats) > 3 {
		panic("Wrong number of arguments passed for cmd " + cmd.protoName)
	}
	prepareCommand(cmd, ints)

	// Exact prototypes per signature: Windows x64 assigns registers by argument position,
	// so a generic (int, int, float, float, float) shim would misplace the floats there.
	var result C.size_t
	switch {
	case len(ints) == 1 && len(floats) == 1:
		result = C.TrampolineFloat1P1F(cmd.fnHandle, C.uintptr_t(ints[0]), C.float(floats[0]))
	case len(ints) == 1 && len(floats) == 2:
		result = C.TrampolineFloat1P2F(cmd.fnHandle, C.uintptr_t(ints[0]), C.float(floats[0]), C.float(floats[1]))
	case len(ints) == 1 && len(floats) == 3:
		result = C.TrampolineFloat1P3F(cmd.fnHandle, C.uintptr_t(ints[0]), C.float(floats[0]), C.float(floats[1]), C.float(floats[2]))
	case len(ints) == 2 && len(floats) == 1:
		result = C.TrampolineFloat2P1F(cmd.fnHandle, C.uintptr_t(ints[0]), C.uintptr_t(ints[1]), C.float(floats[0]))
	default:
		panic("Unhandled float signature for cmd " + cmd.protoName)
	}
	return uintptr(result)
}

// If this package grows real Go callback bridges that can be invoked while inside
// the trampoline helpers, the nocallback directives above must be revisited.

func stringToNullTermBytes(s string) *byte {
	b := []byte(s)
	b = append(b, 0)
	return &b[0]
}

func nullTermBytesToString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	return string(b[:n])
}
