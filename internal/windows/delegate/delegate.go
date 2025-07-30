//go:build windows

// MIT License

// Copyright (c) 2022 SALTO SYSTEMS, S.L

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package delegate

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// Only a limited number of callbacks may be created in a single Go process,
// and any memory allocated for these callbacks is never released.
// Between NewCallback and NewCallbackCDecl, at least 1024 callbacks can always be created.
var (
	queryInterfaceCallback = syscall.NewCallback(queryInterface)
	addRefCallback         = syscall.NewCallback(addRef)
	releaseCallback        = syscall.NewCallback(release)
	invokeCallback         = syscall.NewCallback(invoke)
)

// Delegate represents a WinRT delegate class.
type Delegate interface {
	GetIID() *ole.GUID
	Invoke(instancePtr, rawArgs0, rawArgs1, rawArgs2, rawArgs3, rawArgs4, rawArgs5, rawArgs6, rawArgs7, rawArgs8 unsafe.Pointer) uintptr
	AddRef() uintptr
	Release() uintptr
}

// Callbacks contains the syscalls registered on Windows.
type Callbacks struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Invoke         uintptr
}

var mutex = sync.RWMutex{}
var instances = make(map[uintptr]Delegate)

// RegisterCallbacks adds the given pointer and the Delegate it points to to our instances.
// This is required to redirect received callbacks to the correct object instance.
// The function returns the callbacks to use when creating a new delegate instance.
func RegisterCallbacks(ptr unsafe.Pointer, inst Delegate) *Callbacks {
	mutex.Lock()
	defer mutex.Unlock()
	instances[uintptr(ptr)] = inst

	return &Callbacks{
		QueryInterface: queryInterfaceCallback,
		AddRef:         addRefCallback,
		Release:        releaseCallback,
		Invoke:         invokeCallback,
	}
}

func getInstance(ptr unsafe.Pointer) (Delegate, bool) {
	mutex.RLock() // locks writing, allows concurrent read
	defer mutex.RUnlock()

	i, ok := instances[uintptr(ptr)]
	return i, ok
}

func removeInstance(ptr unsafe.Pointer) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(instances, uintptr(ptr))
}

func queryInterface(instancePtr unsafe.Pointer, iidPtr unsafe.Pointer, ppvObject *unsafe.Pointer) uintptr {
	instance, ok := getInstance(instancePtr)
	if !ok {
		// instance not found
		return ole.E_POINTER
	}

	// Checkout these sources for more information about the QueryInterface method.
	//   - https://docs.microsoft.com/en-us/cpp/atl/queryinterface
	//   - https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)

	if ppvObject == nil {
		// If ppvObject (the address) is nullptr, then this method returns E_POINTER.
		return ole.E_POINTER
	}

	// This function must adhere to the QueryInterface defined here:
	// https://docs.microsoft.com/en-us/windows/win32/api/unknwn/nn-unknwn-iunknown
	if iid := (*ole.GUID)(iidPtr); ole.IsEqualGUID(iid, instance.GetIID()) || ole.IsEqualGUID(iid, ole.IID_IUnknown) || ole.IsEqualGUID(iid, ole.IID_IInspectable) {
		*ppvObject = instancePtr
	} else {
		*ppvObject = nil
		// Return E_NOINTERFACE if the interface is not supported
		return ole.E_NOINTERFACE
	}

	// If the COM object implements the interface, then it returns
	// a pointer to that interface after calling IUnknown::AddRef on it.
	(*ole.IUnknown)(*ppvObject).AddRef()

	// Return S_OK if the interface is supported
	return ole.S_OK
}

func invoke(instancePtr, rawArgs0, rawArgs1, rawArgs2, rawArgs3, rawArgs4, rawArgs5, rawArgs6, rawArgs7, rawArgs8 unsafe.Pointer) uintptr {
	instance, ok := getInstance(instancePtr)
	if !ok {
		// instance not found
		return ole.E_FAIL
	}

	return instance.Invoke(instancePtr, rawArgs0, rawArgs1, rawArgs2, rawArgs3, rawArgs4, rawArgs5, rawArgs6, rawArgs7, rawArgs8)
}

func addRef(instancePtr unsafe.Pointer) uintptr {
	instance, ok := getInstance(instancePtr)
	if !ok {
		// instance not found
		return ole.E_FAIL
	}

	return instance.AddRef()
}

func release(instancePtr unsafe.Pointer) uintptr {
	instance, ok := getInstance(instancePtr)
	if !ok {
		// instance not found
		return ole.E_FAIL
	}

	rem := instance.Release()
	if rem == 0 {
		// remove this delegate
		removeInstance(instancePtr)
	}
	return rem
}
