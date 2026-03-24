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

//export folio_link_new
func folio_link_new(text, uri *C.char, fontH C.uint64_t, fontSize C.double) C.uint64_t {
	f, errCode := loadStandardFont(fontH)
	if errCode != errOK {
		return 0
	}
	lnk := layout.NewLink(C.GoString(text), C.GoString(uri), f, float64(fontSize))
	return C.uint64_t(ht.store(lnk))
}

//export folio_link_new_embedded
func folio_link_new_embedded(text, uri *C.char, fontH C.uint64_t, fontSize C.double) C.uint64_t {
	ef, errCode := loadEmbeddedFont(fontH)
	if errCode != errOK {
		return 0
	}
	lnk := layout.NewLinkEmbedded(C.GoString(text), C.GoString(uri), ef, float64(fontSize))
	return C.uint64_t(ht.store(lnk))
}

//export folio_link_new_internal
func folio_link_new_internal(text, destName *C.char, fontH C.uint64_t, fontSize C.double) C.uint64_t {
	f, errCode := loadStandardFont(fontH)
	if errCode != errOK {
		return 0
	}
	lnk := layout.NewInternalLink(C.GoString(text), C.GoString(destName), f, float64(fontSize))
	return C.uint64_t(ht.store(lnk))
}

//export folio_link_set_color
func folio_link_set_color(linkH C.uint64_t, r, g, b C.double) C.int32_t {
	lnk, errCode := loadLink(linkH)
	if errCode != errOK {
		return errCode
	}
	lnk.SetColor(layout.RGB(float64(r), float64(g), float64(b)))
	return errOK
}

//export folio_link_set_underline
func folio_link_set_underline(linkH C.uint64_t) C.int32_t {
	lnk, errCode := loadLink(linkH)
	if errCode != errOK {
		return errCode
	}
	lnk.SetUnderline()
	return errOK
}

//export folio_link_set_align
func folio_link_set_align(linkH C.uint64_t, align C.int32_t) C.int32_t {
	lnk, errCode := loadLink(linkH)
	if errCode != errOK {
		return errCode
	}
	lnk.SetAlign(layout.Align(align))
	return errOK
}

//export folio_link_free
func folio_link_free(linkH C.uint64_t) {
	ht.delete(uint64(linkH))
}

func loadLink(h C.uint64_t) (*layout.Link, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid link handle")
		return nil, errInvalidHandle
	}
	lnk, ok := v.(*layout.Link)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a link (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return lnk, errOK
}
