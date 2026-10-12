//go:build js && wasm

// Attr: https://www.sprkweb.dev/en/posts/go-wasm-with-typescript/
// Attr: https://permify.co/post/wasm-go/
// Attr: https://pkg.go.dev/syscall/js@go1.24.1

package main

import (
	"fmt"

	"baby-blog/wasm/pages"
	"baby-blog/wasm/pages/home"
	"baby-blog/wasm/pages/test/buttons"

	"github.com/AlexDev404/folio/wasm"
)

func main() {
	// The navigation lives outside the swapped #content, so it's set up once.
	pages.BeginInteractivity(wasm.Route())

	wasm.Page("index", func() func() {
		home.BeginInteractivity()
		return nil
	})
	wasm.Page("test/buttons", func() func() {
		buttons.BeginInteractivity()
		return nil
	})

	fmt.Println("[WASM]: ready")
	wasm.Run()
}
