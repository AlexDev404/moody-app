package main

import (
	"fmt"
	"syscall/js"
)

func init() {
	js.Global().Set("go_setpath", js.FuncOf(go_setpath))
}

func go_setpath(this js.Value, p []js.Value) interface{} {
	fmt.Println(this, p)
	return nil
}

func updateDOMContent() {
	document := js.Global().Get("document")
	element := document.Call("getElementById", "myParagraph")
	element.Set("innerText", "Updated content from Go!")
}

func main() {
	ch := make(chan string)
	fmt.Println("[WASM]: Channel created")
	<-ch // Prevent the program from exiting immediately
}
