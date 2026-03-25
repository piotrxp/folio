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
	"unsafe"

	"github.com/carlos7ags/folio/layout"
	"github.com/carlos7ags/folio/svg"
)

//export folio_svg_parse
func folio_svg_parse(svgXML *C.char) C.uint64_t {
	s, err := svg.Parse(C.GoString(svgXML))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(s))
}

//export folio_svg_parse_bytes
func folio_svg_parse_bytes(data unsafe.Pointer, length C.int32_t) C.uint64_t {
	if data == nil || length <= 0 {
		setLastError("invalid SVG data")
		return 0
	}
	goData := C.GoBytes(data, C.int(length))
	s, err := svg.ParseBytes(goData)
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(s))
}

//export folio_svg_width
func folio_svg_width(svgH C.uint64_t) C.double {
	s, errCode := loadSVG(svgH)
	if errCode != errOK {
		return 0
	}
	return C.double(s.Width())
}

//export folio_svg_height
func folio_svg_height(svgH C.uint64_t) C.double {
	s, errCode := loadSVG(svgH)
	if errCode != errOK {
		return 0
	}
	return C.double(s.Height())
}

//export folio_svg_element_new
func folio_svg_element_new(svgH C.uint64_t) C.uint64_t {
	s, errCode := loadSVG(svgH)
	if errCode != errOK {
		return 0
	}
	se := layout.NewSVGElement(s)
	return C.uint64_t(ht.store(se))
}

//export folio_svg_element_set_size
func folio_svg_element_set_size(seH C.uint64_t, w, h C.double) C.int32_t {
	se, errCode := loadSVGElement(seH)
	if errCode != errOK {
		return errCode
	}
	se.SetSize(float64(w), float64(h))
	return errOK
}

//export folio_svg_element_set_align
func folio_svg_element_set_align(seH C.uint64_t, align C.int32_t) C.int32_t {
	se, errCode := loadSVGElement(seH)
	if errCode != errOK {
		return errCode
	}
	se.SetAlign(layout.Align(align))
	return errOK
}

//export folio_svg_free
func folio_svg_free(svgH C.uint64_t) {
	ht.delete(uint64(svgH))
}

//export folio_svg_element_free
func folio_svg_element_free(seH C.uint64_t) {
	ht.delete(uint64(seH))
}

func loadSVG(h C.uint64_t) (*svg.SVG, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid SVG handle")
		return nil, errInvalidHandle
	}
	s, ok := v.(*svg.SVG)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not an SVG (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return s, errOK
}

func loadSVGElement(h C.uint64_t) (*layout.SVGElement, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid SVG element handle")
		return nil, errInvalidHandle
	}
	se, ok := v.(*layout.SVGElement)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not an SVG element (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return se, errOK
}
