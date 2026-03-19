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
)

//export folio_table_new
func folio_table_new() C.uint64_t {
	t := layout.NewTable()
	return C.uint64_t(ht.store(t))
}

//export folio_table_set_column_widths
func folio_table_set_column_widths(tH C.uint64_t, widths *C.double, count C.int32_t) C.int32_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return errCode
	}
	n := int(count)
	goWidths := make([]float64, n)
	ptr := (*[1 << 20]C.double)(unsafe.Pointer(widths))[:n:n]
	for i := 0; i < n; i++ {
		goWidths[i] = float64(ptr[i])
	}
	t.SetColumnWidths(goWidths)
	return errOK
}

//export folio_table_set_border_collapse
func folio_table_set_border_collapse(tH C.uint64_t, enabled C.int32_t) C.int32_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return errCode
	}
	t.SetBorderCollapse(enabled != 0)
	return errOK
}

//export folio_table_add_row
func folio_table_add_row(tH C.uint64_t) C.uint64_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return 0
	}
	row := t.AddRow()
	return C.uint64_t(ht.store(row))
}

//export folio_table_add_header_row
func folio_table_add_header_row(tH C.uint64_t) C.uint64_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return 0
	}
	row := t.AddHeaderRow()
	return C.uint64_t(ht.store(row))
}

//export folio_row_add_cell
func folio_row_add_cell(rowH C.uint64_t, text *C.char, fontH C.uint64_t, fontSize C.double) C.uint64_t {
	row, errCode := loadRow(rowH)
	if errCode != errOK {
		return 0
	}
	f, errCode := loadStandardFont(fontH)
	if errCode != errOK {
		return 0
	}
	cell := row.AddCell(C.GoString(text), f, float64(fontSize))
	return C.uint64_t(ht.store(cell))
}

//export folio_row_add_cell_embedded
func folio_row_add_cell_embedded(rowH C.uint64_t, text *C.char, fontH C.uint64_t, fontSize C.double) C.uint64_t {
	row, errCode := loadRow(rowH)
	if errCode != errOK {
		return 0
	}
	ef, errCode := loadEmbeddedFont(fontH)
	if errCode != errOK {
		return 0
	}
	cell := row.AddCellEmbedded(C.GoString(text), ef, float64(fontSize))
	return C.uint64_t(ht.store(cell))
}

//export folio_row_add_cell_element
func folio_row_add_cell_element(rowH C.uint64_t, elemH C.uint64_t) C.uint64_t {
	row, errCode := loadRow(rowH)
	if errCode != errOK {
		return 0
	}
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
	cell := row.AddCellElement(elem)
	return C.uint64_t(ht.store(cell))
}

//export folio_cell_set_align
func folio_cell_set_align(cH C.uint64_t, align C.int32_t) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetAlign(layout.Align(align))
	return errOK
}

//export folio_cell_set_padding
func folio_cell_set_padding(cH C.uint64_t, padding C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetPadding(float64(padding))
	return errOK
}

//export folio_cell_set_background
func folio_cell_set_background(cH C.uint64_t, r, g, b C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetBackground(layout.RGB(float64(r), float64(g), float64(b)))
	return errOK
}

//export folio_cell_set_colspan
func folio_cell_set_colspan(cH C.uint64_t, n C.int32_t) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetColspan(int(n))
	return errOK
}

//export folio_cell_set_rowspan
func folio_cell_set_rowspan(cH C.uint64_t, n C.int32_t) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetRowspan(int(n))
	return errOK
}

//export folio_row_free
func folio_row_free(rowH C.uint64_t) {
	ht.delete(uint64(rowH))
}

//export folio_cell_free
func folio_cell_free(cH C.uint64_t) {
	ht.delete(uint64(cH))
}

//export folio_table_free
func folio_table_free(tH C.uint64_t) {
	ht.delete(uint64(tH))
}

func loadTable(h C.uint64_t) (*layout.Table, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid table handle")
		return nil, errInvalidHandle
	}
	t, ok := v.(*layout.Table)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a table (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return t, errOK
}

func loadRow(h C.uint64_t) (*layout.Row, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid row handle")
		return nil, errInvalidHandle
	}
	r, ok := v.(*layout.Row)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a row (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return r, errOK
}

func loadCell(h C.uint64_t) (*layout.Cell, C.int32_t) {
	v := ht.load(uint64(h))
	if v == nil {
		setLastError("invalid cell handle")
		return nil, errInvalidHandle
	}
	c, ok := v.(*layout.Cell)
	if !ok {
		setLastError(fmt.Sprintf("handle %d is not a cell (type %T)", uint64(h), v))
		return nil, errTypeMismatch
	}
	return c, errOK
}
