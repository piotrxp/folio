.PHONY: build test vet fmt fuzz check clean wasm cshared test-cabi

build:
	go build ./...

test:
	go test ./... -race -count=1 -timeout 300s

vet:
	go vet ./...

fmt:
	gofmt -s -w .

fmt-check:
	@test -z "$$(gofmt -s -l .)" || (echo "Run 'make fmt' to fix formatting:" && gofmt -s -l . && exit 1)

fuzz:
	go test ./reader/... -fuzz=FuzzTokenizer -fuzztime=30s || true
	go test ./reader/... -fuzz=FuzzParse -fuzztime=30s || true

check: fmt-check vet test

coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html

wasm:
	GOOS=js GOARCH=wasm go build -o folio.wasm ./cmd/wasm/
	@echo "Built folio.wasm ($$(du -h folio.wasm | cut -f1))"

cshared:
	CGO_ENABLED=1 go build -buildmode=c-shared -o libfolio.dylib ./export/
	@echo "Built libfolio.dylib ($$(du -h libfolio.dylib | cut -f1)) — header: libfolio.h"

test-cabi: cshared
	cc -o export/testdata/test_cabi export/testdata/test_cabi.c -L. -lfolio -Wl,-rpath,.
	./export/testdata/test_cabi

clean:
	rm -f coverage.out coverage.html
	rm -f folio samples showcase
	rm -f folio.wasm
	rm -f libfolio.dylib libfolio.h libfolio.so folio.dll
	rm -f export/testdata/test_cabi
	rm -f *.pdf
