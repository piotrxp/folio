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

//export folio_flex_new
func folio_flex_new() C.uint64_t {
	return C.uint64_t(ht.store(layout.NewFlex()))
}

//export folio_flex_add
func folio_flex_add(flexH C.uint64_t, elemH C.uint64_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
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
	flex.Add(elem)
	return errOK
}

//export folio_flex_add_item
func folio_flex_add_item(flexH C.uint64_t, itemH C.uint64_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	flex.AddItem(item)
	return errOK
}

//export folio_flex_set_direction
func folio_flex_set_direction(flexH C.uint64_t, direction C.int32_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetDirection(layout.FlexDirection(direction))
	return errOK
}

//export folio_flex_set_justify_content
func folio_flex_set_justify_content(flexH C.uint64_t, justify C.int32_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetJustifyContent(layout.JustifyContent(justify))
	return errOK
}

//export folio_flex_set_align_items
func folio_flex_set_align_items(flexH C.uint64_t, align C.int32_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetAlignItems(layout.AlignItems(align))
	return errOK
}

//export folio_flex_set_wrap
func folio_flex_set_wrap(flexH C.uint64_t, wrap C.int32_t) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetWrap(layout.FlexWrap(wrap))
	return errOK
}

//export folio_flex_set_gap
func folio_flex_set_gap(flexH C.uint64_t, gap C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetGap(float64(gap))
	return errOK
}

//export folio_flex_set_padding
func folio_flex_set_padding(flexH C.uint64_t, padding C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetPadding(float64(padding))
	return errOK
}

//export folio_flex_set_background
func folio_flex_set_background(flexH C.uint64_t, r, g, b C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetBackground(layout.RGB(float64(r), float64(g), float64(b)))
	return errOK
}

//export folio_flex_set_space_before
func folio_flex_set_space_before(flexH C.uint64_t, pts C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetSpaceBefore(float64(pts))
	return errOK
}

//export folio_flex_set_space_after
func folio_flex_set_space_after(flexH C.uint64_t, pts C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetSpaceAfter(float64(pts))
	return errOK
}

//export folio_flex_free
func folio_flex_free(flexH C.uint64_t) {
	ht.delete(uint64(flexH))
}

//export folio_flex_item_new
func folio_flex_item_new(elemH C.uint64_t) C.uint64_t {
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
	return C.uint64_t(ht.store(layout.NewFlexItem(elem)))
}

//export folio_flex_item_set_grow
func folio_flex_item_set_grow(itemH C.uint64_t, grow C.double) C.int32_t {
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	item.SetGrow(float64(grow))
	return errOK
}

//export folio_flex_item_set_shrink
func folio_flex_item_set_shrink(itemH C.uint64_t, shrink C.double) C.int32_t {
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	item.SetShrink(float64(shrink))
	return errOK
}

//export folio_flex_item_set_basis
func folio_flex_item_set_basis(itemH C.uint64_t, basis C.double) C.int32_t {
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	item.SetBasis(float64(basis))
	return errOK
}

//export folio_flex_item_set_align_self
func folio_flex_item_set_align_self(itemH C.uint64_t, align C.int32_t) C.int32_t {
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	item.SetAlignSelf(layout.AlignItems(align))
	return errOK
}

//export folio_flex_item_set_margins
func folio_flex_item_set_margins(itemH C.uint64_t, top, right, bottom, left C.double) C.int32_t {
	item, errCode := loadFlexItem(itemH)
	if errCode != errOK {
		return errCode
	}
	item.SetMargins(float64(top), float64(right), float64(bottom), float64(left))
	return errOK
}

//export folio_flex_item_free
func folio_flex_item_free(itemH C.uint64_t) {
	ht.delete(uint64(itemH))
}

func loadFlex(h C.uint64_t) (*layout.Flex, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid flex handle")
		return nil, errInvalidHandle
	}
	f, ok := v.(*layout.Flex)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a flex container (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return f, errOK
}

func loadFlexItem(h C.uint64_t) (*layout.FlexItem, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid flex item handle")
		return nil, errInvalidHandle
	}
	item, ok := v.(*layout.FlexItem)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a flex item (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return item, errOK
}
