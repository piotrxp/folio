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

	"github.com/carlos7ags/folio/reader"
)

//export folio_reader_open
func folio_reader_open(path *C.char) C.uint64_t {
	r, err := reader.Open(C.GoString(path))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(r))
}

//export folio_reader_parse
func folio_reader_parse(data unsafe.Pointer, length C.int32_t) C.uint64_t {
	if data == nil || length <= 0 {
		setLastError("invalid PDF data")
		return 0
	}
	goData := C.GoBytes(data, C.int(length))
	r, err := reader.Parse(goData)
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(r))
}

//export folio_reader_page_count
func folio_reader_page_count(rH C.uint64_t) C.int32_t {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	return C.int32_t(r.PageCount())
}

//export folio_reader_version
func folio_reader_version(rH C.uint64_t) C.uint64_t {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	return C.uint64_t(ht.store(newCBuffer([]byte(r.Version()))))
}

//export folio_reader_info_title
func folio_reader_info_title(rH C.uint64_t) C.uint64_t {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	title, _, _, _, _ := r.Info()
	return C.uint64_t(ht.store(newCBuffer([]byte(title))))
}

//export folio_reader_info_author
func folio_reader_info_author(rH C.uint64_t) C.uint64_t {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	_, author, _, _, _ := r.Info()
	return C.uint64_t(ht.store(newCBuffer([]byte(author))))
}

//export folio_reader_extract_text
func folio_reader_extract_text(rH C.uint64_t, pageIndex C.int32_t) C.uint64_t {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	page, err := r.Page(int(pageIndex))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	text, err := page.ExtractText()
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(newCBuffer([]byte(text))))
}

//export folio_reader_page_width
func folio_reader_page_width(rH C.uint64_t, pageIndex C.int32_t) C.double {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	page, err := r.Page(int(pageIndex))
	if err != nil {
		return 0
	}
	return C.double(page.Width)
}

//export folio_reader_page_height
func folio_reader_page_height(rH C.uint64_t, pageIndex C.int32_t) C.double {
	r, errCode := loadReader(rH)
	if errCode != errOK {
		return 0
	}
	page, err := r.Page(int(pageIndex))
	if err != nil {
		return 0
	}
	return C.double(page.Height)
}

//export folio_reader_free
func folio_reader_free(rH C.uint64_t) {
	ht.delete(uint64(rH))
}

func loadReader(h C.uint64_t) (*reader.PdfReader, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid reader handle")
		return nil, errInvalidHandle
	}
	r, ok := v.(*reader.PdfReader)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a reader (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return r, errOK
}
