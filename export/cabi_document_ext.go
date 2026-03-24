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
	"sync"

	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/layout"
)

//export folio_document_set_watermark
func folio_document_set_watermark(docH C.uint64_t, text *C.char) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetWatermark(C.GoString(text))
	return errOK
}

//export folio_document_set_watermark_config
func folio_document_set_watermark_config(docH C.uint64_t, text *C.char,
	fontSize, colorR, colorG, colorB, angle, opacity C.double) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetWatermarkConfig(document.WatermarkConfig{
		Text:     C.GoString(text),
		FontSize: float64(fontSize),
		ColorR:   float64(colorR),
		ColorG:   float64(colorG),
		ColorB:   float64(colorB),
		Angle:    float64(angle),
		Opacity:  float64(opacity),
	})
	return errOK
}

//export folio_document_add_outline
func folio_document_add_outline(docH C.uint64_t, title *C.char, pageIndex C.int32_t) C.uint64_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return 0
	}
	outline := doc.AddOutline(C.GoString(title), document.FitDest(int(pageIndex)))
	return C.uint64_t(ht.store(outline))
}

//export folio_document_add_outline_xyz
func folio_document_add_outline_xyz(docH C.uint64_t, title *C.char,
	pageIndex C.int32_t, left, top, zoom C.double) C.uint64_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return 0
	}
	outline := doc.AddOutline(C.GoString(title),
		document.XYZDest(int(pageIndex), float64(left), float64(top), float64(zoom)))
	return C.uint64_t(ht.store(outline))
}

//export folio_outline_add_child
func folio_outline_add_child(outlineH C.uint64_t, title *C.char, pageIndex C.int32_t) C.uint64_t {
	o, errCode := loadOutline(outlineH)
	if errCode != errOK {
		return 0
	}
	child := o.AddChild(C.GoString(title), document.FitDest(int(pageIndex)))
	return C.uint64_t(ht.store(child))
}

//export folio_outline_add_child_xyz
func folio_outline_add_child_xyz(outlineH C.uint64_t, title *C.char,
	pageIndex C.int32_t, left, top, zoom C.double) C.uint64_t {
	o, errCode := loadOutline(outlineH)
	if errCode != errOK {
		return 0
	}
	child := o.AddChild(C.GoString(title),
		document.XYZDest(int(pageIndex), float64(left), float64(top), float64(zoom)))
	return C.uint64_t(ht.store(child))
}

//export folio_document_add_named_dest
func folio_document_add_named_dest(docH C.uint64_t, name *C.char, pageIndex C.int32_t,
	fitType *C.char, top, left, zoom C.double) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.AddNamedDest(document.NamedDest{
		Name:      C.GoString(name),
		PageIndex: int(pageIndex),
		FitType:   C.GoString(fitType),
		Top:       float64(top),
		Left:      float64(left),
		Zoom:      float64(zoom),
	})
	return errOK
}

//export folio_document_set_viewer_preferences
func folio_document_set_viewer_preferences(docH C.uint64_t,
	pageLayout, pageMode *C.char,
	hideToolbar, hideMenubar, hideWindowUI, fitWindow, centerWindow, displayDocTitle C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	doc.SetViewerPreferences(document.ViewerPreferences{
		PageLayout:      document.PageLayout(C.GoString(pageLayout)),
		PageMode:        document.PageMode(C.GoString(pageMode)),
		HideToolbar:     hideToolbar != 0,
		HideMenubar:     hideMenubar != 0,
		HideWindowUI:    hideWindowUI != 0,
		FitWindow:       fitWindow != 0,
		CenterWindow:    centerWindow != 0,
		DisplayDocTitle: displayDocTitle != 0,
	})
	return errOK
}

//export folio_document_add_page_label
func folio_document_add_page_label(docH C.uint64_t, pageIndex C.int32_t,
	style, prefix *C.char, start C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	label := document.PageLabelRange{
		PageIndex: int(pageIndex),
		Style:     document.LabelStyle(C.GoString(style)),
		Prefix:    C.GoString(prefix),
		Start:     int(start),
	}
	pageLabelsMuImpl.Lock()
	labels := pageLabels[uint64(docH)]
	labels = append(labels, label)
	pageLabels[uint64(docH)] = labels
	pageLabelsMuImpl.Unlock()
	doc.SetPageLabels(labels...)
	return errOK
}

//export folio_document_remove_page
func folio_document_remove_page(docH C.uint64_t, index C.int32_t) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	if err := doc.RemovePage(int(index)); err != nil {
		return setErr(errInvalidArg, err)
	}
	return errOK
}

//export folio_document_add_absolute
func folio_document_add_absolute(docH C.uint64_t, elemH C.uint64_t, x, y, width C.double) C.int32_t {
	doc, errCode := loadDoc(docH)
	if errCode != errOK {
		return errCode
	}
	v := ht.load(uint64(elemH))
	if v == nil {
		setLastError("invalid element handle")
		return errInvalidHandle
	}
	elem, ok := v.(layout.Element)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a layout element", uint64(elemH)))
		return errTypeMismatch
	}
	doc.AddAbsolute(elem, float64(x), float64(y), float64(width))
	return errOK
}

//export folio_page_add_internal_link
func folio_page_add_internal_link(pageH C.uint64_t, x1, y1, x2, y2 C.double, destName *C.char) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	page.AddInternalLink([4]float64{float64(x1), float64(y1), float64(x2), float64(y2)}, C.GoString(destName))
	return errOK
}

//export folio_page_add_text_annotation
func folio_page_add_text_annotation(pageH C.uint64_t, x1, y1, x2, y2 C.double, text, icon *C.char) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	page.AddTextAnnotation(
		[4]float64{float64(x1), float64(y1), float64(x2), float64(y2)},
		C.GoString(text), C.GoString(icon))
	return errOK
}

//export folio_page_set_crop_box
func folio_page_set_crop_box(pageH C.uint64_t, x1, y1, x2, y2 C.double) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	page.SetCropBox([4]float64{float64(x1), float64(y1), float64(x2), float64(y2)})
	return errOK
}

//export folio_page_set_trim_box
func folio_page_set_trim_box(pageH C.uint64_t, x1, y1, x2, y2 C.double) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	page.SetTrimBox([4]float64{float64(x1), float64(y1), float64(x2), float64(y2)})
	return errOK
}

//export folio_page_set_bleed_box
func folio_page_set_bleed_box(pageH C.uint64_t, x1, y1, x2, y2 C.double) C.int32_t {
	page, errCode := loadPage(pageH)
	if errCode != errOK {
		return errCode
	}
	page.SetBleedBox([4]float64{float64(x1), float64(y1), float64(x2), float64(y2)})
	return errOK
}

// ── Helpers ────────────────────────────────────────────────────────

var (
	pageLabelsMuImpl sync.Mutex
	pageLabels       = make(map[uint64][]document.PageLabelRange)
)

func loadOutline(h C.uint64_t) (*document.Outline, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid outline handle")
		return nil, errInvalidHandle
	}
	o, ok := v.(*document.Outline)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not an outline (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return o, errOK
}
