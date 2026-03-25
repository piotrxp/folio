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

	folioimage "github.com/carlos7ags/folio/image"
	"github.com/carlos7ags/folio/layout"
)

// ── Paragraph extensions ───────────────────────────────────────────

//export folio_paragraph_set_orphans
func folio_paragraph_set_orphans(pH C.uint64_t, n C.int32_t) C.int32_t {
	p, errCode := loadParagraph(pH)
	if errCode != errOK {
		return errCode
	}
	p.SetOrphans(int(n))
	return errOK
}

//export folio_paragraph_set_widows
func folio_paragraph_set_widows(pH C.uint64_t, n C.int32_t) C.int32_t {
	p, errCode := loadParagraph(pH)
	if errCode != errOK {
		return errCode
	}
	p.SetWidows(int(n))
	return errOK
}

//export folio_paragraph_set_ellipsis
func folio_paragraph_set_ellipsis(pH C.uint64_t, enabled C.int32_t) C.int32_t {
	p, errCode := loadParagraph(pH)
	if errCode != errOK {
		return errCode
	}
	p.SetEllipsis(enabled != 0)
	return errOK
}

//export folio_paragraph_set_word_break
func folio_paragraph_set_word_break(pH C.uint64_t, wb *C.char) C.int32_t {
	p, errCode := loadParagraph(pH)
	if errCode != errOK {
		return errCode
	}
	p.SetWordBreak(C.GoString(wb))
	return errOK
}

//export folio_paragraph_set_hyphens
func folio_paragraph_set_hyphens(pH C.uint64_t, h *C.char) C.int32_t {
	p, errCode := loadParagraph(pH)
	if errCode != errOK {
		return errCode
	}
	p.SetHyphens(C.GoString(h))
	return errOK
}

// ── Table extensions ───────────────────────────────────────────────

//export folio_table_add_footer_row
func folio_table_add_footer_row(tH C.uint64_t) C.uint64_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return 0
	}
	return C.uint64_t(ht.store(t.AddFooterRow()))
}

//export folio_table_set_cell_spacing
func folio_table_set_cell_spacing(tH C.uint64_t, h, v C.double) C.int32_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return errCode
	}
	t.SetCellSpacing(float64(h), float64(v))
	return errOK
}

//export folio_table_set_auto_column_widths
func folio_table_set_auto_column_widths(tH C.uint64_t) C.int32_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return errCode
	}
	t.SetAutoColumnWidths()
	return errOK
}

//export folio_table_set_min_width
func folio_table_set_min_width(tH C.uint64_t, pts C.double) C.int32_t {
	t, errCode := loadTable(tH)
	if errCode != errOK {
		return errCode
	}
	t.SetMinWidth(float64(pts))
	return errOK
}

// ── Cell extensions ────────────────────────────────────────────────

//export folio_cell_set_padding_sides
func folio_cell_set_padding_sides(cH C.uint64_t, top, right, bottom, left C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetPaddingSides(layout.Padding{
		Top: float64(top), Right: float64(right),
		Bottom: float64(bottom), Left: float64(left),
	})
	return errOK
}

//export folio_cell_set_valign
func folio_cell_set_valign(cH C.uint64_t, valign C.int32_t) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetVAlign(layout.VAlign(valign))
	return errOK
}

//export folio_cell_set_borders
func folio_cell_set_borders(cH C.uint64_t,
	topW, topR, topG, topB C.double,
	rightW, rightR, rightG, rightB C.double,
	bottomW, bottomR, bottomG, bottomB C.double,
	leftW, leftR, leftG, leftB C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetBorders(layout.CellBorders{
		Top:    layout.Border{Width: float64(topW), Color: layout.RGB(float64(topR), float64(topG), float64(topB))},
		Right:  layout.Border{Width: float64(rightW), Color: layout.RGB(float64(rightR), float64(rightG), float64(rightB))},
		Bottom: layout.Border{Width: float64(bottomW), Color: layout.RGB(float64(bottomR), float64(bottomG), float64(bottomB))},
		Left:   layout.Border{Width: float64(leftW), Color: layout.RGB(float64(leftR), float64(leftG), float64(leftB))},
	})
	return errOK
}

//export folio_cell_set_border
func folio_cell_set_border(cH C.uint64_t, width, r, g, b C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	b1 := layout.Border{Width: float64(width), Color: layout.RGB(float64(r), float64(g), float64(b))}
	cell.SetBorders(layout.CellBorders{Top: b1, Right: b1, Bottom: b1, Left: b1})
	return errOK
}

//export folio_cell_set_width_hint
func folio_cell_set_width_hint(cH C.uint64_t, pts C.double) C.int32_t {
	cell, errCode := loadCell(cH)
	if errCode != errOK {
		return errCode
	}
	cell.SetWidthHint(float64(pts))
	return errOK
}

// ── Div extensions ─────────────────────────────────────────────────

//export folio_div_set_border_radius
func folio_div_set_border_radius(divH C.uint64_t, r C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetBorderRadius(float64(r))
	return errOK
}

//export folio_div_set_opacity
func folio_div_set_opacity(divH C.uint64_t, opacity C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetOpacity(float64(opacity))
	return errOK
}

//export folio_div_set_overflow
func folio_div_set_overflow(divH C.uint64_t, overflow *C.char) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetOverflow(C.GoString(overflow))
	return errOK
}

//export folio_div_set_max_width
func folio_div_set_max_width(divH C.uint64_t, pts C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetMaxWidth(float64(pts))
	return errOK
}

//export folio_div_set_min_width
func folio_div_set_min_width(divH C.uint64_t, pts C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetMinWidth(float64(pts))
	return errOK
}

// ── List extensions ────────────────────────────────────────────────

//export folio_list_set_leading
func folio_list_set_leading(listH C.uint64_t, leading C.double) C.int32_t {
	l, errCode := loadList(listH)
	if errCode != errOK {
		return errCode
	}
	l.SetLeading(float64(leading))
	return errOK
}

//export folio_list_add_nested_item
func folio_list_add_nested_item(listH C.uint64_t, text *C.char) C.uint64_t {
	l, errCode := loadList(listH)
	if errCode != errOK {
		return 0
	}
	return C.uint64_t(ht.store(l.AddItemWithSubList(C.GoString(text))))
}

// ── Image element extensions ───────────────────────────────────────

//export folio_image_element_set_align
func folio_image_element_set_align(ieH C.uint64_t, align C.int32_t) C.int32_t {
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
	ie.SetAlign(layout.Align(align))
	return errOK
}

// ── Flex extensions ────────────────────────────────────────────────

//export folio_flex_set_row_gap
func folio_flex_set_row_gap(flexH C.uint64_t, gap C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetRowGap(float64(gap))
	return errOK
}

//export folio_flex_set_column_gap
func folio_flex_set_column_gap(flexH C.uint64_t, gap C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetColumnGap(float64(gap))
	return errOK
}

//export folio_flex_set_padding_all
func folio_flex_set_padding_all(flexH C.uint64_t, top, right, bottom, left C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetPaddingAll(layout.Padding{
		Top: float64(top), Right: float64(right),
		Bottom: float64(bottom), Left: float64(left),
	})
	return errOK
}

//export folio_flex_set_border
func folio_flex_set_border(flexH C.uint64_t, width, r, g, b C.double) C.int32_t {
	flex, errCode := loadFlex(flexH)
	if errCode != errOK {
		return errCode
	}
	flex.SetBorder(layout.SolidBorder(float64(width), layout.RGB(float64(r), float64(g), float64(b))))
	return errOK
}

// ── Div box shadow ─────────────────────────────────────────────────

//export folio_div_set_box_shadow
func folio_div_set_box_shadow(divH C.uint64_t,
	offsetX, offsetY, blur, spread C.double, r, g, b C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetBoxShadow(layout.BoxShadow{
		OffsetX: float64(offsetX),
		OffsetY: float64(offsetY),
		Blur:    float64(blur),
		Spread:  float64(spread),
		Color:   layout.RGB(float64(r), float64(g), float64(b)),
	})
	return errOK
}

//export folio_div_set_max_height
func folio_div_set_max_height(divH C.uint64_t, pts C.double) C.int32_t {
	div, errCode := loadDiv(divH)
	if errCode != errOK {
		return errCode
	}
	div.SetMaxHeight(float64(pts))
	return errOK
}

// ── Image TIFF support ─────────────────────────────────────────────

//export folio_image_load_tiff
func folio_image_load_tiff(path *C.char) C.uint64_t {
	img, err := folioimage.LoadTIFF(C.GoString(path))
	if err != nil {
		setLastError(err.Error())
		return 0
	}
	return C.uint64_t(ht.store(img))
}
