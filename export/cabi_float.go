// Copyright 2026 Carlos Munoz and the Folio Authors
// SPDX-License-Identifier: Apache-2.0

//go:build cgo && !js && !wasm

package main

/*
#include <stdint.h>
*/
import "C"
import (
	"fmt"

	"github.com/carlos7ags/folio/layout"
)

//export folio_float_new
func folio_float_new(side C.int32_t, elemH C.uint64_t) C.uint64_t {
	v := ht.load(uint64(elemH))
	if v == nil {
		setLastError("invalid element handle")
		return 0
	}
	elem, ok := v.(layout.Element)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a layout element", uint64(elemH)))
		return 0
	}
	f := layout.NewFloat(layout.FloatSide(side), elem)
	return C.uint64_t(ht.store(f))
}

//export folio_float_set_margin
func folio_float_set_margin(floatH C.uint64_t, margin C.double) C.int32_t {
	f, errCode := loadFloat(floatH)
	if errCode != errOK {
		return errCode
	}
	f.SetMargin(float64(margin))
	return errOK
}

//export folio_float_free
func folio_float_free(floatH C.uint64_t) {
	ht.delete(uint64(floatH))
}

func loadFloat(h C.uint64_t) (*layout.Float, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid float handle")
		return nil, errInvalidHandle
	}
	f, ok := v.(*layout.Float)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a float (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return f, errOK
}
