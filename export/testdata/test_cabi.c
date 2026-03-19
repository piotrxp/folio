/*
 * C integration test for the Folio C ABI.
 * Compile and run:
 *   cc -o test_cabi test_cabi.c -L../.. -lfolio -Wl,-rpath,../..
 *   ./test_cabi
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>

/* Use the auto-generated header from the build */
#include "../../libfolio.h"

#define ASSERT(cond, msg) do { \
    if (!(cond)) { \
        fprintf(stderr, "FAIL: %s (line %d)\n", msg, __LINE__); \
        const char* err = folio_last_error(); \
        if (err) fprintf(stderr, "  last_error: %s\n", err); \
        failures++; \
    } else { \
        passes++; \
    } \
} while(0)

int main(void) {
    int passes = 0, failures = 0;

    /* Test 1: Version string — persistent pointer, do NOT free */
    const char* ver = folio_version();
    ASSERT(ver != NULL, "folio_version returns non-null");
    ASSERT(strlen(ver) > 0, "version string is non-empty");
    printf("folio version: %s\n", ver);

    /* Test 2: Create blank document and save to buffer */
    uint64_t doc = folio_document_new_letter();
    ASSERT(doc != 0, "document_new_letter returns handle");

    int32_t rc = folio_document_set_title(doc, "C ABI Test");
    ASSERT(rc == 0, "set_title succeeds");

    rc = folio_document_set_author(doc, "Test Author");
    ASSERT(rc == 0, "set_author succeeds");

    rc = folio_document_set_margins(doc, 36, 36, 36, 36);
    ASSERT(rc == 0, "set_margins succeeds");

    uint64_t page = folio_document_add_page(doc);
    ASSERT(page != 0, "add_page returns handle");

    int32_t count = folio_document_page_count(doc);
    ASSERT(count == 1, "page_count is 1");

    uint64_t buf = folio_document_write_to_buffer(doc);
    ASSERT(buf != 0, "write_to_buffer returns handle");

    int32_t len = folio_buffer_len(buf);
    ASSERT(len > 0, "buffer has data");

    void* data = folio_buffer_data(buf);
    ASSERT(data != NULL, "buffer data is non-null");
    ASSERT(memcmp(data, "%PDF-1.7", 8) == 0, "buffer starts with PDF header");

    folio_buffer_free(buf);
    folio_document_free(doc);

    /* Test 3: Text on page with standard font */
    doc = folio_document_new(595.28, 841.89); /* A4 */
    ASSERT(doc != 0, "document_new with custom size");

    folio_document_set_title(doc, "Font Test");

    page = folio_document_add_page(doc);
    ASSERT(page != 0, "add_page for font test");

    uint64_t helv = folio_font_helvetica();
    ASSERT(helv != 0, "font_helvetica returns handle");

    rc = folio_page_add_text(page, "Hello from C!", helv, 24.0, 72.0, 700.0);
    ASSERT(rc == 0, "page_add_text succeeds");

    uint64_t times = folio_font_times_roman();
    ASSERT(times != 0, "font_times_roman returns handle");

    rc = folio_page_add_text(page, "Second line in Times.", times, 12.0, 72.0, 660.0);
    ASSERT(rc == 0, "page_add_text with Times succeeds");

    rc = folio_page_add_link(page, 72.0, 640.0, 200.0, 655.0, "https://folio.dev");
    ASSERT(rc == 0, "page_add_link succeeds");

    /* Save to file */
    rc = folio_document_save(doc, "/tmp/folio_cabi_test.pdf");
    ASSERT(rc == 0, "document_save succeeds");

    folio_document_free(doc);

    /* Test 4: Invalid handle */
    rc = folio_document_set_title(99999, "bad");
    ASSERT(rc != 0, "invalid handle returns error");
    ASSERT(folio_last_error() != NULL, "last_error set for invalid handle");

    /* Test 5: Font lookup by name */
    uint64_t courier = folio_font_standard("Courier");
    ASSERT(courier != 0, "font_standard Courier");

    uint64_t bad_font = folio_font_standard("NotAFont");
    ASSERT(bad_font == 0, "unknown font returns 0");

    /* Test 6: Layout engine — paragraphs with word wrapping */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Layout Test");

    helv = folio_font_helvetica();
    uint64_t para = folio_paragraph_new("This is a paragraph that should wrap automatically when it exceeds the page width. The layout engine handles word wrapping and page breaks.", helv, 12.0);
    ASSERT(para != 0, "paragraph_new returns handle");

    rc = folio_paragraph_set_align(para, 0); /* AlignLeft */
    ASSERT(rc == 0, "paragraph_set_align succeeds");

    rc = folio_paragraph_set_leading(para, 1.5);
    ASSERT(rc == 0, "paragraph_set_leading succeeds");

    rc = folio_paragraph_set_space_after(para, 12.0);
    ASSERT(rc == 0, "paragraph_set_space_after succeeds");

    rc = folio_paragraph_set_background(para, 0.95, 0.95, 0.95);
    ASSERT(rc == 0, "paragraph_set_background succeeds");

    rc = folio_paragraph_set_first_indent(para, 24.0);
    ASSERT(rc == 0, "paragraph_set_first_indent succeeds");

    rc = folio_document_add(doc, para);
    ASSERT(rc == 0, "document_add paragraph succeeds");

    /* Test 7: Heading */
    uint64_t h1 = folio_heading_new("Chapter 1: Introduction", 1);
    ASSERT(h1 != 0, "heading_new returns handle");

    rc = folio_heading_set_align(h1, 1); /* AlignCenter */
    ASSERT(rc == 0, "heading_set_align succeeds");

    rc = folio_document_add(doc, h1);
    ASSERT(rc == 0, "document_add heading succeeds");

    /* Test 8: Heading with specific font */
    uint64_t h2 = folio_heading_new_with_font("Section 1.1", 2, folio_font_helvetica_bold(), 18.0);
    ASSERT(h2 != 0, "heading_new_with_font returns handle");

    rc = folio_document_add(doc, h2);
    ASSERT(rc == 0, "document_add heading_with_font succeeds");

    /* Test 9: Paragraph with mixed runs */
    uint64_t styled = folio_paragraph_new("Bold start: ", folio_font_helvetica_bold(), 12.0);
    ASSERT(styled != 0, "styled paragraph created");

    rc = folio_paragraph_add_run(styled, "normal continuation.", helv, 12.0, 0.0, 0.0, 0.0);
    ASSERT(rc == 0, "paragraph_add_run succeeds");

    rc = folio_paragraph_add_run(styled, " Red text.", helv, 12.0, 1.0, 0.0, 0.0);
    ASSERT(rc == 0, "paragraph_add_run with color succeeds");

    rc = folio_document_add(doc, styled);
    ASSERT(rc == 0, "document_add styled paragraph succeeds");

    /* Save layout document */
    rc = folio_document_save(doc, "/tmp/folio_cabi_layout.pdf");
    ASSERT(rc == 0, "document_save layout succeeds");

    folio_paragraph_free(para);
    folio_heading_free(h1);
    folio_heading_free(h2);
    folio_paragraph_free(styled);
    folio_document_free(doc);

    /* Test 10: Font free — standard fonts are no-op */
    folio_font_free(helv); /* should not crash */
    uint64_t helv_again = folio_font_helvetica();
    ASSERT(helv_again != 0, "standard font still available after free");

    /* ===== Stage 5: Tables ===== */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Table Test");
    helv = folio_font_helvetica();

    uint64_t tbl = folio_table_new();
    ASSERT(tbl != 0, "table_new returns handle");

    rc = folio_table_set_border_collapse(tbl, 1);
    ASSERT(rc == 0, "table_set_border_collapse succeeds");

    /* Header row */
    uint64_t hrow = folio_table_add_header_row(tbl);
    ASSERT(hrow != 0, "table_add_header_row returns handle");

    uint64_t c1 = folio_row_add_cell(hrow, "Name", helv, 12.0);
    ASSERT(c1 != 0, "row_add_cell returns handle");
    folio_cell_set_background(c1, 0.9, 0.9, 0.9);

    uint64_t c2 = folio_row_add_cell(hrow, "Value", helv, 12.0);
    ASSERT(c2 != 0, "row_add_cell 2 returns handle");

    /* Data row */
    uint64_t drow = folio_table_add_row(tbl);
    ASSERT(drow != 0, "table_add_row returns handle");
    folio_row_add_cell(drow, "Folio", helv, 12.0);
    folio_row_add_cell(drow, "PDF Library", helv, 12.0);

    rc = folio_document_add(doc, tbl);
    ASSERT(rc == 0, "document_add table succeeds");

    rc = folio_document_save(doc, "/tmp/folio_cabi_table.pdf");
    ASSERT(rc == 0, "table document save succeeds");
    folio_table_free(tbl);
    folio_document_free(doc);

    /* ===== Stage 7: Containers (Div, List, AreaBreak) ===== */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Container Test");
    helv = folio_font_helvetica();

    uint64_t div = folio_div_new();
    ASSERT(div != 0, "div_new returns handle");

    rc = folio_div_set_padding(div, 10, 10, 10, 10);
    ASSERT(rc == 0, "div_set_padding succeeds");

    rc = folio_div_set_background(div, 0.95, 0.95, 1.0);
    ASSERT(rc == 0, "div_set_background succeeds");

    rc = folio_div_set_border(div, 1.0, 0.0, 0.0, 0.5);
    ASSERT(rc == 0, "div_set_border succeeds");

    /* Add a paragraph inside the div */
    para = folio_paragraph_new("Content inside a div.", helv, 12.0);
    rc = folio_div_add(div, para);
    ASSERT(rc == 0, "div_add paragraph succeeds");

    rc = folio_document_add(doc, div);
    ASSERT(rc == 0, "document_add div succeeds");

    /* List */
    uint64_t list = folio_list_new(helv, 12.0);
    ASSERT(list != 0, "list_new returns handle");

    folio_list_set_style(list, 1); /* ListOrdered */
    folio_list_add_item(list, "First item");
    folio_list_add_item(list, "Second item");
    folio_list_add_item(list, "Third item");

    rc = folio_document_add(doc, list);
    ASSERT(rc == 0, "document_add list succeeds");

    /* Area break */
    uint64_t brk = folio_area_break_new();
    ASSERT(brk != 0, "area_break_new returns handle");
    rc = folio_document_add(doc, brk);
    ASSERT(rc == 0, "document_add area_break succeeds");

    /* Line separator */
    uint64_t sep = folio_line_separator_new();
    ASSERT(sep != 0, "line_separator_new returns handle");
    rc = folio_document_add(doc, sep);
    ASSERT(rc == 0, "document_add line_separator succeeds");

    rc = folio_document_save(doc, "/tmp/folio_cabi_containers.pdf");
    ASSERT(rc == 0, "containers document save succeeds");
    folio_div_free(div);
    folio_list_free(list);
    folio_document_free(doc);

    /* ===== Stage 8: HTML to PDF ===== */
    rc = folio_html_to_pdf("<h1>Hello from C</h1><p>This PDF was generated from HTML via the C ABI.</p>", "/tmp/folio_cabi_html.pdf");
    ASSERT(rc == 0, "html_to_pdf succeeds");

    uint64_t htmlBuf = folio_html_to_buffer("<h1>Buffer Test</h1>", 612, 792);
    ASSERT(htmlBuf != 0, "html_to_buffer returns handle");
    ASSERT(folio_buffer_len(htmlBuf) > 0, "html buffer has data");
    folio_buffer_free(htmlBuf);

    uint64_t htmlDoc = folio_html_convert("<h1>Convert Test</h1><p>Paragraph</p>", 612, 792);
    ASSERT(htmlDoc != 0, "html_convert returns doc handle");
    folio_document_set_title(htmlDoc, "HTML Convert");
    rc = folio_document_save(htmlDoc, "/tmp/folio_cabi_html_convert.pdf");
    ASSERT(rc == 0, "html_convert doc save succeeds");
    folio_document_free(htmlDoc);

    /* ===== Stage 9: Reader ===== */
    /* Read back the HTML PDF we just created */
    uint64_t rdr = folio_reader_open("/tmp/folio_cabi_html.pdf");
    ASSERT(rdr != 0, "reader_open succeeds");

    int32_t pageCount = folio_reader_page_count(rdr);
    ASSERT(pageCount >= 1, "reader page_count >= 1");

    double pw = folio_reader_page_width(rdr, 0);
    ASSERT(pw > 0, "reader page_width > 0");

    double ph = folio_reader_page_height(rdr, 0);
    ASSERT(ph > 0, "reader page_height > 0");

    uint64_t textBuf = folio_reader_extract_text(rdr, 0);
    ASSERT(textBuf != 0, "reader extract_text returns handle");
    ASSERT(folio_buffer_len(textBuf) > 0, "extracted text is non-empty");
    folio_buffer_free(textBuf);

    folio_reader_free(rdr);

    /* ===== Stage 10: Forms & Document Features ===== */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Forms Test");
    folio_document_add_page(doc);

    uint64_t form = folio_form_new();
    ASSERT(form != 0, "form_new returns handle");

    rc = folio_form_add_text_field(form, "name", 72, 700, 300, 720, 0);
    ASSERT(rc == 0, "form_add_text_field succeeds");

    rc = folio_form_add_checkbox(form, "agree", 72, 650, 90, 668, 0, 1);
    ASSERT(rc == 0, "form_add_checkbox succeeds");

    rc = folio_document_set_form(doc, form);
    ASSERT(rc == 0, "document_set_form succeeds");

    rc = folio_document_save(doc, "/tmp/folio_cabi_forms.pdf");
    ASSERT(rc == 0, "forms document save succeeds");
    folio_form_free(form);
    folio_document_free(doc);

    /* Document features */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Features Test");

    rc = folio_document_set_tagged(doc, 1);
    ASSERT(rc == 0, "set_tagged succeeds");

    rc = folio_document_set_auto_bookmarks(doc, 1);
    ASSERT(rc == 0, "set_auto_bookmarks succeeds");

    folio_document_free(doc);

    /* ===== Stage 11: Callbacks ===== */
    doc = folio_document_new_letter();
    folio_document_set_title(doc, "Callback Test");
    folio_document_add_page(doc);

    /* NULL callback must be rejected */
    rc = folio_document_set_header(doc, NULL, NULL);
    ASSERT(rc != 0, "NULL header callback rejected");

    rc = folio_document_set_footer(doc, NULL, NULL);
    ASSERT(rc != 0, "NULL footer callback rejected");

    folio_document_free(doc);

    /* Summary */
    printf("\n%d passed, %d failed\n", passes, failures);
    return failures > 0 ? 1 : 0;
}
