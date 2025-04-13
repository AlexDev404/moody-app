package home

import "syscall/js"

func BeginInteractivity() {
	DOM := js.Global().Get("document")
	moodMeter := DOM.Call("getElementById", "mood-input")
	// For each ux_palette_item, set the onclick function
	uxPaletteItems := DOM.Call("querySelectorAll", "ux-palette-item")
	length := uxPaletteItems.Get("length").Int()

	for i := 0; i < length; i++ {
		item := uxPaletteItems.Call("item", i)
		mood := item.Get("dataset").Get("mood").String()
		mood = "I feel " + mood

		item.Call("addEventListener", "click", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
			moodMeter.Set("value", mood)
			return nil
		}))
	}

}
