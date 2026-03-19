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

	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/forms"
)

// --- AcroForm ---

//export folio_form_new
func folio_form_new() C.uint64_t {
	return C.uint64_t(ht.store(forms.NewAcroForm()))
}

//export folio_form_add_text_field
func folio_form_add_text_field(formH C.uint64_t, name *C.char, x1, y1, x2, y2 C.double, pageIndex C.int32_t) C.int32_t {
	af, errCode := loadForm(formH)
	if errCode != errOK {
		return errCode
	}
	rect := [4]float64{float64(x1), float64(y1), float64(x2), float64(y2)}
	af.Add(forms.TextField(C.GoString(name), rect, int(pageIndex)))
	return errOK
}

//export folio_form_add_checkbox
func folio_form_add_checkbox(formH C.uint64_t, name *C.char, x1, y1, x2, y2 C.double, pageIndex C.int32_t, checked C.int32_t) C.int32_t {
	af, errCode := loadForm(formH)
	if errCode != errOK {
		return errCode
	}
	rect := [4]float64{float64(x1), float64(y1), float64(x2), float64(y2)}
	af.Add(forms.Checkbox(C.GoString(name), rect, int(pageIndex), checked != 0))
	return errOK
}

//export folio_form_add_dropdown
func folio_form_add_dropdown(formH C.uint64_t, name *C.char, x1, y1, x2, y2 C.double, pageIndex C.int32_t, options **C.char, optCount C.int32_t) C.int32_t {
	af, errCode := loadForm(formH)
	if errCode != errOK {
		return errCode
	}
	rect := [4]float64{float64(x1), float64(y1), float64(x2), float64(y2)}
	n := int(optCount)
	goOpts := make([]string, n)
	if n > 0 && options != nil {
		cArray := (*[1 << 20]*C.char)(unsafe.Pointer(options))[:n:n]
		for i := 0; i < n; i++ {
			goOpts[i] = C.GoString(cArray[i])
		}
	}
	af.Add(forms.Dropdown(C.GoString(name), rect, int(pageIndex), goOpts))
	return errOK
}

//export folio_form_add_signature
func folio_form_add_signature(formH C.uint64_t, name *C.char, x1, y1, x2, y2 C.double, pageIndex C.int32_t) C.int32_t {
	af, errCode := loadForm(formH)
	if errCode != errOK {
		return errCode
	}
	rect := [4]float64{float64(x1), float64(y1), float64(x2), float64(y2)}
	af.Add(forms.SignatureField(C.GoString(name), rect, int(pageIndex)))
	return errOK
}

//export folio_document_set_form
func folio_document_set_form(docH C.uint64_t, formH C.uint64_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	af, errCode := loadForm(formH)
	if errCode != errOK {
		return errCode
	}
	doc.SetAcroForm(af)
	return errOK
}

//export folio_form_free
func folio_form_free(formH C.uint64_t) {
	ht.delete(uint64(formH))
}

// --- Document feature flags ---

//export folio_document_set_tagged
func folio_document_set_tagged(docH C.uint64_t, enabled C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetTagged(enabled != 0)
	return errOK
}

//export folio_document_set_pdfa
func folio_document_set_pdfa(docH C.uint64_t, level C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetPdfA(document.PdfAConfig{Level: document.PdfALevel(level)})
	return errOK
}

//export folio_document_set_encryption
func folio_document_set_encryption(docH C.uint64_t, userPw, ownerPw *C.char, algorithm C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetEncryption(document.EncryptionConfig{
		Algorithm:     document.EncryptionAlgorithm(algorithm),
		UserPassword:  C.GoString(userPw),
		OwnerPassword: C.GoString(ownerPw),
	})
	return errOK
}

//export folio_document_set_auto_bookmarks
func folio_document_set_auto_bookmarks(docH C.uint64_t, enabled C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetAutoBookmarks(enabled != 0)
	return errOK
}

func loadForm(h C.uint64_t) (*forms.AcroForm, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid form handle")
		return nil, errInvalidHandle
	}
	af, ok := v.(*forms.AcroForm)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a form (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return af, errOK
}
