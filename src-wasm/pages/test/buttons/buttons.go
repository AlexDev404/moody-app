package buttons

import "syscall/js"

func Begin_Interactivity() {
	document := js.Global().Get("document")
	element := document.Call("getElementById", "button1")
	element.Set("innerText", "Updated content from Go!")
}
