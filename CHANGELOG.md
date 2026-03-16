# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-03-16

### Added
- CLI `extract` command with pluggable strategies (simple, location)
- CLI `sign` command for PAdES digital signatures
- Table `SetBorderCollapse(true)` for CSS-style collapsed borders
- CSS `calc()` support in HTML-to-PDF (e.g., `width: calc(100% - 40px)`)
- CSS `@page` rule parsing (size, margins) in HTML-to-PDF
- CSS `orphans`/`widows` properties in HTML-to-PDF
- CSS `break-before`/`break-after`/`break-inside` modern syntax
- Remote image loading (`<img src="https://...">`) in HTML-to-PDF
- Data URI image support (`<img src="data:image/png;base64,...">`)
- PDF metadata extraction from HTML `<title>` and `<meta>` tags
- Content stream processor with full graphics state (CTM, color, font)
- Pluggable text extraction strategies (Simple, Location, Region)
- Path and image extraction from content streams
- Per-glyph span extraction (opt-in)
- Text rendering mode awareness (invisible text filtering)
- Marked content tag tracking (BMC/BDC/EMC)
- Form XObject recursion in content processing
- Actual glyph widths from font metrics (replaces estimation)
- Auto-bookmarks from layout headings
- Viewer preferences (page layout, mode, UI options)
- Page labels (decimal, Roman, alpha)
- Page geometry boxes (CropBox, BleedBox, TrimBox, ArtBox)
- SVG package in README

### Changed
- CLI version bumped to 0.1.1
- README updated with extract and sign commands, border-collapse, SVG package

### Fixed
- Table border-collapse: adjacent cells no longer draw double borders
- Tables section in README had undefined variable

## [0.1.0] - 2026-03-15

### Added
- Initial release
- PDF generation with layout engine (Paragraph, Heading, Table, List, Div, Image, Float, Flex, Columns)
- PDF reader with tokenizer, parser, xref streams, object streams
- PDF merge and modify
- HTML-to-PDF conversion with CSS support
- Digital signatures (PAdES B-B, B-T, B-LT)
- Interactive forms (AcroForms)
- Barcodes (Code128, QR, EAN-13)
- Tagged PDF and PDF/A compliance
- SVG rendering
- CLI tool (merge, info, pages, text, create, blank)
- Font embedding and subsetting (TrueType)
- JPEG, PNG, TIFF image support
- Encryption (AES-256, AES-128, RC4)
