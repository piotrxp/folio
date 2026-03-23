package html

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	folioimage "github.com/carlos7ags/folio/image"
	"github.com/carlos7ags/folio/layout"
)

func TestImageInterceptor(t *testing.T) {
	elems, err := Convert(`<img src="photo.jpg"/>`, &Options{
		ImageInterceptor: func(src string) (*folioimage.Image, error) {
			log.Printf("ImageInterceptor called with src: %s, returning error", src)
			return nil, fmt.Errorf("Loading external images is not allowed")
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(elems) != 1 {
		t.Fatalf("expected 1 element, got %d", len(elems))
	}

	// element should be a layout.Paragraph element, not an image element, since the interceptor returned an error to prevent loading
	_, ok := elems[0].(*layout.Paragraph)
	if !ok {
		t.Fatalf("expected a Paragraph element, got %T", elems[0])
	}

	// now test it with an interceptor that returns an image
	elems, err = Convert(`<img src="photo.jpg"/>`, &Options{
		ImageInterceptor: func(src string) (*folioimage.Image, error) {
			pngBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==" // black 1px x 1px PNG
			log.Printf("ImageInterceptor called with src: %s, returning dummy PNG image", src)

			pngBytes, err := base64.StdEncoding.DecodeString(pngBase64)
			if err != nil {
				return nil, err
			}
			return folioimage.NewPNG(pngBytes)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(elems) != 1 {
		t.Fatalf("expected 1 element, got %d", len(elems))
	}

	// element should be a layout.ImageElement element, not a Paragraph element, since the interceptor returned an image
	_, ok = elems[0].(*layout.ImageElement)
	if !ok {
		t.Fatalf("expected an ImageElement element, got %T", elems[0])
	}
}
