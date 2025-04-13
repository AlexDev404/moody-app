package pages

import (
	"syscall/js"
)

// Global variables to avoid repetitive calls
var DOM = js.Global().Get("document")
var window = js.Global().Get("window")

func BeginInteractivity(applicationPath string) {
	// For each ux-navigation, set the onclick function
	uxNavigation := DOM.Call("querySelectorAll", "ux-navigation")
	length := uxNavigation.Get("length").Int()

	for i := 0; i < length; i++ {

		item := uxNavigation.Call("item", i)
		setupNavigationItems(item)
	}

	// Set initial state based on current location
	updateActiveNavigationItem()

	// Listen for popstate events
	window.Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updateActiveNavigationItem()
		return nil
	}))
}

func setupNavigationItems(navElement js.Value) {

	navItems := navElement.Call("querySelectorAll", "ux-navigation-item")

	length := navItems.Get("length").Int()

	for i := 0; i < length; i++ {
		item := navItems.Call("item", i)
		href := item.Call("getAttribute", "href").String()

		// Setup click handler
		item.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			args[0].Call("preventDefault")

			// Push the new state
			js.Global().Get("history").Call("pushState", nil, "", href)

			// Update active item
			updateActiveNavigationItem()

			return nil
		}))
	}
}

func updateActiveNavigationItem() {
	currentPath := window.Get("location").Get("pathname").String()

	if currentPath == "" {
		currentPath = "/"
	}

	// Get all navigation items
	navItems := DOM.Call("querySelectorAll", "ux-navigation-item")
	length := navItems.Get("length").Int()

	for i := 0; i < length; i++ {
		item := navItems.Call("item", i)
		href := item.Call("getAttribute", "href").String()

		isActive := href == currentPath || (href == "/" && currentPath == "")
		updateItemStyle(item, isActive)
	}
}

func updateItemStyle(item js.Value, isActive bool) {
	// Update background class
	item.Get("classList").Call("remove", "bg-zinc-700", "bg-transparent")
	if isActive {
		item.Get("classList").Call("add", "bg-zinc-700")
	} else {
		item.Get("classList").Call("add", "bg-transparent")
	}

	// Update icon color
	icon := item.Call("querySelector", "ux-icon")
	if !icon.IsUndefined() && !icon.IsNull() {
		iconColor := "#FCFCFC"
		if isActive {
			iconColor = "#7DFCBC"
		}
		icon.Call("setAttribute", "color", iconColor)
	}
}
