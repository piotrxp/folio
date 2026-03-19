/*
 * Folio C ABI — PDF generation library
 * https://github.com/carlos7ags/folio
 * Apache 2.0 License
 *
 * MEMORY OWNERSHIP RULES
 * ----------------------
 * 1. All objects are opaque uint64_t handles. The library owns the memory.
 *    Every folio_*_new() / folio_*_load() has a matching folio_*_free().
 *
 * 2. Strings passed TO the library (const char*) are copied immediately.
 *    The caller retains ownership and may free after the call returns.
 *
 * 3. Strings returned FROM the library:
 *    - folio_version(): persistent pointer, do NOT free.
 *    - folio_last_error(): library-owned, valid until the next C ABI call.
 *    - All other string data is returned as buffer handles (see below).
 *
 * 4. Buffer handles (folio_buffer_data / folio_buffer_len / folio_buffer_free):
 *    The library allocates the buffer; the caller MUST call folio_buffer_free()
 *    when done. The data pointer is valid until folio_buffer_free() is called.
 *
 * ERROR CONVENTION
 * ----------------
 * Functions returning int32_t: 0 = success, negative = error.
 * Call folio_last_error() for the human-readable message.
 *
 * THREAD SAFETY
 * -------------
 * The library is NOT thread-safe. All calls must be serialized.
 * A single-threaded wrapper is recommended for multi-threaded applications.
 */

#ifndef FOLIO_H
#define FOLIO_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/* ── Error codes ───────────────────────────────────────────────────── */

#define FOLIO_OK               0
#define FOLIO_ERR_HANDLE      -1   /* invalid or expired handle */
#define FOLIO_ERR_ARG         -2   /* invalid argument (NULL, out of range) */
#define FOLIO_ERR_IO          -3   /* file I/O error */
#define FOLIO_ERR_PDF         -4   /* PDF generation/parsing error */
#define FOLIO_ERR_TYPE        -5   /* handle type mismatch */
#define FOLIO_ERR_INTERNAL    -6   /* unexpected internal error */

/* ── Enums ─────────────────────────────────────────────────────────── */

/* Text alignment (layout.Align) */
#define FOLIO_ALIGN_LEFT       0
#define FOLIO_ALIGN_CENTER     1
#define FOLIO_ALIGN_RIGHT      2
#define FOLIO_ALIGN_JUSTIFY    3

/* Heading levels (layout.HeadingLevel) */
#define FOLIO_H1               1
#define FOLIO_H2               2
#define FOLIO_H3               3
#define FOLIO_H4               4
#define FOLIO_H5               5
#define FOLIO_H6               6

/* List styles (layout.ListStyle) */
#define FOLIO_LIST_BULLET      0
#define FOLIO_LIST_DECIMAL     1
#define FOLIO_LIST_LOWER_ALPHA 2
#define FOLIO_LIST_UPPER_ALPHA 3
#define FOLIO_LIST_LOWER_ROMAN 4
#define FOLIO_LIST_UPPER_ROMAN 5

/* PDF/A levels (document.PdfALevel) */
#define FOLIO_PDFA_2B          0
#define FOLIO_PDFA_2U          1
#define FOLIO_PDFA_2A          2
#define FOLIO_PDFA_3B          3
#define FOLIO_PDFA_1B          4
#define FOLIO_PDFA_1A          5

/* Encryption algorithms (document.EncryptionAlgorithm) */
#define FOLIO_ENCRYPT_RC4_128  0
#define FOLIO_ENCRYPT_AES_128  1
#define FOLIO_ENCRYPT_AES_256  2

/* ── Callback types ────────────────────────────────────────────────── */

/*
 * Page decorator callback. Invoked for each page during rendering.
 * page_handle is a temporary handle valid only inside the callback.
 */
typedef void (*folio_page_decorator_fn)(
    int32_t   page_index,
    int32_t   total_pages,
    uint64_t  page_handle,
    void     *user_data
);

/* ── Core ──────────────────────────────────────────────────────────── */

/* Returns the library version string. Persistent — do NOT free. */
const char *folio_version(void);

/* Returns the last error message. Library-owned — do NOT free.
 * Valid until the next folio_* call. NULL if no error. */
const char *folio_last_error(void);

/* ── Buffer ────────────────────────────────────────────────────────── */

/* Returns a pointer to the buffer's data. Valid until folio_buffer_free(). */
void    *folio_buffer_data(uint64_t buf);
int32_t  folio_buffer_len(uint64_t buf);
void     folio_buffer_free(uint64_t buf);

/* ── Document ──────────────────────────────────────────────────────── */

uint64_t folio_document_new(double width, double height);
uint64_t folio_document_new_letter(void);
uint64_t folio_document_new_a4(void);
void     folio_document_free(uint64_t doc);

int32_t  folio_document_set_title(uint64_t doc, const char *title);
int32_t  folio_document_set_author(uint64_t doc, const char *author);
int32_t  folio_document_set_margins(uint64_t doc, double top, double right, double bottom, double left);

uint64_t folio_document_add_page(uint64_t doc);  /* returns page handle */
int32_t  folio_document_page_count(uint64_t doc);
int32_t  folio_document_add(uint64_t doc, uint64_t element);  /* any layout element */

int32_t  folio_document_save(uint64_t doc, const char *path);
uint64_t folio_document_write_to_buffer(uint64_t doc);  /* returns buffer handle */

/* Document features */
int32_t  folio_document_set_tagged(uint64_t doc, int32_t enabled);
int32_t  folio_document_set_pdfa(uint64_t doc, int32_t level);  /* FOLIO_PDFA_* */
int32_t  folio_document_set_encryption(uint64_t doc, const char *user_pw, const char *owner_pw, int32_t algorithm);
int32_t  folio_document_set_auto_bookmarks(uint64_t doc, int32_t enabled);
int32_t  folio_document_set_form(uint64_t doc, uint64_t form);

/* Callbacks — fn must NOT be NULL */
int32_t  folio_document_set_header(uint64_t doc, folio_page_decorator_fn fn, void *user_data);
int32_t  folio_document_set_footer(uint64_t doc, folio_page_decorator_fn fn, void *user_data);

/* ── Page ──────────────────────────────────────────────────────────── */

int32_t  folio_page_add_text(uint64_t page, const char *text, uint64_t font, double size, double x, double y);
int32_t  folio_page_add_text_embedded(uint64_t page, const char *text, uint64_t font, double size, double x, double y);
int32_t  folio_page_add_image(uint64_t page, uint64_t img, double x, double y, double w, double h);
int32_t  folio_page_add_link(uint64_t page, double x1, double y1, double x2, double y2, const char *uri);
int32_t  folio_page_set_opacity(uint64_t page, double alpha);
int32_t  folio_page_set_rotate(uint64_t page, int32_t degrees);

/* ── Font ──────────────────────────────────────────────────────────── */

/* Standard PDF fonts (singletons — folio_font_free is a no-op) */
uint64_t folio_font_standard(const char *name);  /* e.g. "Helvetica", "Courier-Bold" */
uint64_t folio_font_helvetica(void);
uint64_t folio_font_helvetica_bold(void);
uint64_t folio_font_times_roman(void);
uint64_t folio_font_times_bold(void);
uint64_t folio_font_courier(void);

/* Embedded fonts (TTF/OTF) — must call folio_font_free when done */
uint64_t folio_font_load_ttf(const char *path);
uint64_t folio_font_parse_ttf(const void *data, int32_t length);
void     folio_font_free(uint64_t font);  /* no-op for standard fonts */

/* ── Paragraph ─────────────────────────────────────────────────────── */

uint64_t folio_paragraph_new(const char *text, uint64_t font, double font_size);
uint64_t folio_paragraph_new_embedded(const char *text, uint64_t font, double font_size);
void     folio_paragraph_free(uint64_t para);

int32_t  folio_paragraph_set_align(uint64_t para, int32_t align);  /* FOLIO_ALIGN_* */
int32_t  folio_paragraph_set_leading(uint64_t para, double leading);
int32_t  folio_paragraph_set_space_before(uint64_t para, double pts);
int32_t  folio_paragraph_set_space_after(uint64_t para, double pts);
int32_t  folio_paragraph_set_background(uint64_t para, double r, double g, double b);
int32_t  folio_paragraph_set_first_indent(uint64_t para, double pts);

/* Add a styled text run. Font can be standard or embedded. Color is RGB 0-1. */
int32_t  folio_paragraph_add_run(uint64_t para, const char *text, uint64_t font, double font_size, double r, double g, double b);

/* ── Heading ───────────────────────────────────────────────────────── */

uint64_t folio_heading_new(const char *text, int32_t level);  /* FOLIO_H1..H6 */
uint64_t folio_heading_new_with_font(const char *text, int32_t level, uint64_t font, double font_size);
uint64_t folio_heading_new_embedded(const char *text, int32_t level, uint64_t font);
void     folio_heading_free(uint64_t heading);

int32_t  folio_heading_set_align(uint64_t heading, int32_t align);

/* ── Table ─────────────────────────────────────────────────────────── */

uint64_t folio_table_new(void);
void     folio_table_free(uint64_t table);

int32_t  folio_table_set_column_widths(uint64_t table, const double *widths, int32_t count);
int32_t  folio_table_set_border_collapse(uint64_t table, int32_t enabled);

uint64_t folio_table_add_row(uint64_t table);         /* returns row handle */
uint64_t folio_table_add_header_row(uint64_t table);   /* returns row handle */
void     folio_row_free(uint64_t row);

uint64_t folio_row_add_cell(uint64_t row, const char *text, uint64_t font, double font_size);
uint64_t folio_row_add_cell_embedded(uint64_t row, const char *text, uint64_t font, double font_size);
uint64_t folio_row_add_cell_element(uint64_t row, uint64_t element);
void     folio_cell_free(uint64_t cell);

int32_t  folio_cell_set_align(uint64_t cell, int32_t align);
int32_t  folio_cell_set_padding(uint64_t cell, double padding);
int32_t  folio_cell_set_background(uint64_t cell, double r, double g, double b);
int32_t  folio_cell_set_colspan(uint64_t cell, int32_t n);
int32_t  folio_cell_set_rowspan(uint64_t cell, int32_t n);

/* ── Image ─────────────────────────────────────────────────────────── */

uint64_t folio_image_load_jpeg(const char *path);
uint64_t folio_image_load_png(const char *path);
uint64_t folio_image_parse_jpeg(const void *data, int32_t length);
uint64_t folio_image_parse_png(const void *data, int32_t length);
int32_t  folio_image_width(uint64_t img);
int32_t  folio_image_height(uint64_t img);
void     folio_image_free(uint64_t img);

/* Image as layout element (for document flow) */
uint64_t folio_image_element_new(uint64_t img);
int32_t  folio_image_element_set_size(uint64_t elem, double w, double h);
void     folio_image_element_free(uint64_t elem);

/* ── Div (container) ───────────────────────────────────────────────── */

uint64_t folio_div_new(void);
void     folio_div_free(uint64_t div);

int32_t  folio_div_add(uint64_t div, uint64_t element);
int32_t  folio_div_set_padding(uint64_t div, double top, double right, double bottom, double left);
int32_t  folio_div_set_background(uint64_t div, double r, double g, double b);
int32_t  folio_div_set_border(uint64_t div, double width, double r, double g, double b);
int32_t  folio_div_set_width(uint64_t div, double pts);
int32_t  folio_div_set_min_height(uint64_t div, double pts);
int32_t  folio_div_set_space_before(uint64_t div, double pts);
int32_t  folio_div_set_space_after(uint64_t div, double pts);

/* ── List ──────────────────────────────────────────────────────────── */

uint64_t folio_list_new(uint64_t font, double font_size);
uint64_t folio_list_new_embedded(uint64_t font, double font_size);
void     folio_list_free(uint64_t list);

int32_t  folio_list_set_style(uint64_t list, int32_t style);  /* FOLIO_LIST_* */
int32_t  folio_list_set_indent(uint64_t list, double indent);
int32_t  folio_list_add_item(uint64_t list, const char *text);

/* ── Misc layout elements ──────────────────────────────────────────── */

uint64_t folio_line_separator_new(void);
uint64_t folio_area_break_new(void);

/* ── HTML to PDF ───────────────────────────────────────────────────── */

/* One-shot: HTML string → PDF file */
int32_t  folio_html_to_pdf(const char *html, const char *output_path);

/* HTML string → buffer handle (caller must folio_buffer_free) */
uint64_t folio_html_to_buffer(const char *html, double page_width, double page_height);

/* HTML string → document handle (for further manipulation before save) */
uint64_t folio_html_convert(const char *html, double page_width, double page_height);

/* ── Reader (PDF parsing) ──────────────────────────────────────────── */

uint64_t folio_reader_open(const char *path);
uint64_t folio_reader_parse(const void *data, int32_t length);
void     folio_reader_free(uint64_t reader);

int32_t  folio_reader_page_count(uint64_t reader);
uint64_t folio_reader_version(uint64_t reader);       /* returns buffer handle */
uint64_t folio_reader_info_title(uint64_t reader);     /* returns buffer handle */
uint64_t folio_reader_info_author(uint64_t reader);    /* returns buffer handle */
uint64_t folio_reader_extract_text(uint64_t reader, int32_t page_index);  /* returns buffer handle */
double   folio_reader_page_width(uint64_t reader, int32_t page_index);
double   folio_reader_page_height(uint64_t reader, int32_t page_index);

/* ── Forms ─────────────────────────────────────────────────────────── */

uint64_t folio_form_new(void);
void     folio_form_free(uint64_t form);

int32_t  folio_form_add_text_field(uint64_t form, const char *name, double x1, double y1, double x2, double y2, int32_t page_index);
int32_t  folio_form_add_checkbox(uint64_t form, const char *name, double x1, double y1, double x2, double y2, int32_t page_index, int32_t checked);
int32_t  folio_form_add_dropdown(uint64_t form, const char *name, double x1, double y1, double x2, double y2, int32_t page_index, const char **options, int32_t opt_count);
int32_t  folio_form_add_signature(uint64_t form, const char *name, double x1, double y1, double x2, double y2, int32_t page_index);

#ifdef __cplusplus
}
#endif

#endif /* FOLIO_H */
