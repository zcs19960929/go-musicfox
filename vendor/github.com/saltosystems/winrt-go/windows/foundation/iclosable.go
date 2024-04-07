// Code generated by winrt-go-gen. DO NOT EDIT.

//go:build windows

//nolint:all
package foundation

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

const GUIDIClosable string = "30d5a829-7fa4-4026-83bb-d75bae4ea99e"
const SignatureIClosable string = "{30d5a829-7fa4-4026-83bb-d75bae4ea99e}"

type IClosable struct {
	ole.IInspectable
}

type IClosableVtbl struct {
	ole.IInspectableVtbl

	Close uintptr
}

func (v *IClosable) VTable() *IClosableVtbl {
	return (*IClosableVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IClosable) Close() error {
	hr, _, _ := syscall.SyscallN(
		v.VTable().Close,
		uintptr(unsafe.Pointer(v)), // this
	)

	if hr != 0 {
		return ole.NewError(hr)
	}

	return nil
}