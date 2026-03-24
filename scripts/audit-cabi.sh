#!/usr/bin/env bash
# audit-cabi.sh — Compare Go exports, C header declarations, and built symbols.
#
# Usage:
#   ./scripts/audit-cabi.sh          # audit only (no build)
#   ./scripts/audit-cabi.sh --build  # build dylib and also verify symbols
#
# Exit codes:
#   0  everything in sync
#   1  drift detected

set -euo pipefail
cd "$(dirname "$0")/.."

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'
BOLD='\033[1m'

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

# ── 1. Extract Go //export directives ──────────────────────────────

grep -rh '//export folio_' export/*.go \
  | sed 's|.*//export ||' \
  | sort -u \
  > "$tmpdir/go_exports.txt"

go_count=$(wc -l < "$tmpdir/go_exports.txt" | tr -d ' ')
echo -e "${BOLD}Go //export directives:${NC} $go_count"

# ── 2. Extract C header declarations ───────────────────────────────

# Match function declarations: return_type folio_name(...)
grep -oE 'folio_[a-zA-Z0-9_]+' export/folio.h \
  | grep -v '^folio_page_decorator_fn$' \
  | sort -u \
  > "$tmpdir/header_funcs.txt"

header_count=$(wc -l < "$tmpdir/header_funcs.txt" | tr -d ' ')
echo -e "${BOLD}Header declarations:${NC}   $header_count"

# ── 3. Compare Go exports vs header ───────────────────────────────

in_go_not_header=$(comm -23 "$tmpdir/go_exports.txt" "$tmpdir/header_funcs.txt" || true)
in_header_not_go=$(comm -13 "$tmpdir/go_exports.txt" "$tmpdir/header_funcs.txt" || true)

drift=0

if [ -n "$in_go_not_header" ]; then
  echo ""
  echo -e "${RED}Exported in Go but MISSING from folio.h:${NC}"
  echo "$in_go_not_header" | sed 's/^/  - /'
  drift=1
fi

if [ -n "$in_header_not_go" ]; then
  echo ""
  echo -e "${RED}Declared in folio.h but MISSING Go //export:${NC}"
  echo "$in_header_not_go" | sed 's/^/  - /'
  drift=1
fi

if [ "$drift" -eq 0 ]; then
  echo -e "${GREEN}✓ Go exports and folio.h are in sync${NC}"
fi

# ── 4. Optionally verify built symbols ─────────────────────────────

if [ "${1:-}" = "--build" ]; then
  echo ""
  echo -e "${BOLD}Building shared library...${NC}"
  outlib="$tmpdir/libfolio_audit"

  case "$(uname -s)" in
    Darwin) outlib="$outlib.dylib" ;;
    Linux)  outlib="$outlib.so" ;;
    *)      outlib="$outlib.dll" ;;
  esac

  CGO_ENABLED=1 go build -o "$outlib" -buildmode=c-shared ./export/ 2>&1

  # Extract symbols
  case "$(uname -s)" in
    Darwin) nm -gU "$outlib" | grep ' T _folio_' | sed 's/.* T _//' | sort -u > "$tmpdir/symbols.txt" ;;
    *)      nm -D "$outlib" | grep ' T folio_' | sed 's/.* T //' | sort -u > "$tmpdir/symbols.txt" ;;
  esac

  sym_count=$(wc -l < "$tmpdir/symbols.txt" | tr -d ' ')
  echo -e "${BOLD}Built symbols:${NC}         $sym_count"

  in_go_not_sym=$(comm -23 "$tmpdir/go_exports.txt" "$tmpdir/symbols.txt" || true)
  in_sym_not_go=$(comm -13 "$tmpdir/go_exports.txt" "$tmpdir/symbols.txt" || true)

  if [ -n "$in_go_not_sym" ]; then
    echo ""
    echo -e "${RED}Go //export but NOT in built library:${NC}"
    echo "$in_go_not_sym" | sed 's/^/  - /'
    drift=1
  fi

  if [ -n "$in_sym_not_go" ]; then
    echo ""
    echo -e "${YELLOW}In built library but no Go //export (runtime/cgo internals):${NC}"
    echo "$in_sym_not_go" | sed 's/^/  - /'
  fi

  if [ "$drift" -eq 0 ]; then
    echo -e "${GREEN}✓ Built symbols match Go exports${NC}"
  fi
fi

# ── 5. Summary ─────────────────────────────────────────────────────

echo ""
if [ "$drift" -ne 0 ]; then
  echo -e "${RED}${BOLD}DRIFT DETECTED${NC} — header and Go exports are out of sync"
  exit 1
else
  echo -e "${GREEN}${BOLD}ALL CLEAR${NC} — $go_count functions exported, header matches"
  exit 0
fi
