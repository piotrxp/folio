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

	folioimage "github.com/carlos7ags/folio/image"
	"github.com/carlos7ags/folio/layout"
)

//export folio_image_load_jpeg
func folio_image_load_jpeg(path *C.char) C.uint64_t {
	img, err := folioimage.LoadJPEG(C.GoString(path))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(img))
}

//export folio_image_load_png
func folio_image_load_png(path *C.char) C.uint64_t {
	img, err := folioimage.LoadPNG(C.GoString(path))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(img))
}

//export folio_image_parse_jpeg
func folio_image_parse_jpeg(data unsafe.Pointer, length C.int32_t) C.uint64_t {
	if data == nil || length <= 0 {
		setLastError("invalid image data")
		return 0
	}
	goData := C.GoBytes(data, C.int(length))
	img, err := folioimage.NewJPEG(goData)
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(img))
}

//export folio_image_parse_png
func folio_image_parse_png(data unsafe.Pointer, length C.int32_t) C.uint64_t {
	if data == nil || length <= 0 {
		setLastError("invalid image data")
		return 0
	}
	goData := C.GoBytes(data, C.int(length))
	img, err := folioimage.NewPNG(goData)
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(img))
}

//export folio_image_width
func folio_image_width(imgH C.uint64_t) C.int32_t {
	img, errCode := loadImage(imgH)
	if errCode != errOK {
		return 0
	}
	return C.int32_t(img.Width())
}

//export folio_image_height
func folio_image_height(imgH C.uint64_t) C.int32_t {
	img, errCode := loadImage(imgH)
	if errCode != errOK {
		return 0
	}
	return C.int32_t(img.Height())
}

//export folio_image_free
func folio_image_free(imgH C.uint64_t) {
	ht.delete(uint64(imgH))
}

//export folio_page_add_image
func folio_page_add_image(pageH C.uint64_t, imgH C.uint64_t, x, y, w, h C.double) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	img, errCode := loadImage(imgH)
	if errCode != errOK {
		return errCode
	}
	page.AddImage(img, float64(x), float64(y), float64(w), float64(h))
	return errOK
}

//export folio_image_element_new
func folio_image_element_new(imgH C.uint64_t) C.uint64_t {
	img, errCode := loadImage(imgH)
	if errCode != errOK {
		return 0
	}
	ie := layout.NewImageElement(img)
	return C.uint64_t(ht.store(ie))
}

//export folio_image_element_set_size
func folio_image_element_set_size(ieH C.uint64_t, w, h C.double) C.int32_t {
	v := ht.load(uint64(ieH))
	if v == nil {
		setLastError("invalid image element handle")
		return errInvalidHandle
	}
	ie, ok := v.(*layout.ImageElement)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not an image element", uint64(ieH)))
		return errTypeMismatch
	}
	ie.SetSize(float64(w), float64(h))
	return errOK
}

//export folio_image_element_free
func folio_image_element_free(ieH C.uint64_t) {
	ht.delete(uint64(ieH))
}

func loadImage(h C.uint64_t) (*folioimage.Image, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid image handle")
		return nil, errInvalidHandle
	}
	img, ok := v.(*folioimage.Image)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not an image (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return img, errOK
}
