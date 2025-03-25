package buttons

import "syscall/js"

func BeginInteractivity() {
	setButtonProperties("button1", "Updated content from Go!", func(this js.Value, p []js.Value) interface{} {
		js.Global().Call("alert", "Button clicked!")
		return nil
	})

	setButtonProperties("button2", "Updated content from Go!", func(this js.Value, p []js.Value) interface{} {
		js.Global().Get("console").Call("log", "Button clicked!")
		this.Set("innerText", "Updated Console!")
		return nil
	})

	setButtonProperties("button3", "Updated content from Go!", func(this js.Value, p []js.Value) interface{} {
		js.Global().Get("console").Call("warn", "Button clicked!")
		this.Set("innerText", "Updated Console!")
		return nil
	})

	setButtonProperties("button4", "Updated content from Go!", func(this js.Value, p []js.Value) interface{} {
		js.Global().Get("localStorage").Call("setItem", "button4", "Button clicked!")
		this.Set("innerText", "Updated LocalStorage!")
		return nil
	})

	setButtonProperties("button5", "Updated content from Go!", func(this js.Value, p []js.Value) interface{} {
		js.Global().Call("fetch", "https://jsonplaceholder.typicode.com/posts/1").
			Call("then", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
				js.Global().Get("console").Call("log", p[0])
				return nil
			}))
		this.Set("innerText", "Updated Console!")
		return nil
	})
}

func setButtonProperties(buttonID, innerText string, onClickFunc func(this js.Value, p []js.Value) interface{}) {
	element := js.Global().Get("document").Call("getElementById", buttonID)
	element.Set("innerText", innerText)
	element.Set("onclick", js.FuncOf(onClickFunc))
}
