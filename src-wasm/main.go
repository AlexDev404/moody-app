package main

import (
	"fmt"
	"syscall/js"
)

type WasmApplication struct {
	Path string
}

func (application *WasmApplication) go_setpath(this js.Value, p []js.Value) interface{} {
	application.Path = p[0].String()
	application.updateDOMContent()
	return nil
}

func (application *WasmApplication) updateDOMContent() {
	document := js.Global().Get("document")
	element := document.Call("getElementById", "button1")
	element.Set("innerText", "Updated content from Go!")
}

func (application *WasmApplication) init() {
	js.Global().Set("go_setpath", js.FuncOf(application.go_setpath))
}
func main() {
	application := &WasmApplication{}
	application.init()
	ch := make(chan string)
	fmt.Println("[WASM]: Channel created")
	<-ch // Prevent the program from exiting immediately
}
