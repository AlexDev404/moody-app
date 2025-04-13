package main

import (
	"bytes"
	"net/http"

	"baby-blog/types"
	"html/template"
)

// @title Render
// @description Very basic function. Render is a function that renders a template and writes it to the response writer. Built for humans.
// @description It uses the template library to parse the template and write it to the response writer.
// @description It also uses the buffer pool to get a buffer and write it to the response writer.
// @param w http.ResponseWriter
// @param r *http.Request
// @param t *template.Template
// @param pageData map[string]interface{}
// @return void
// @example Render(w, r, t, pageData)
func (app *Application) Render(w http.ResponseWriter, r *http.Request, t *template.Template, pageData map[string]interface{}) {
	path := app.getPath(r)
	if app.isDisallowedRoute(path) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if pageData == nil {
		pageData = map[string]interface{}{
			"Path":    path,
			"Errors":  map[string]string{},
			"Message": nil,
			"Failure": nil,
		}
	}

	pageData = app.runHooks(pageData, r, w)

	contentData := &types.TemplateData{
		Data: map[string]interface{}{
			"Path":     path,
			"PageData": pageData,
		},
	}

	var templateContent string

	// Check if the template exists in the template cache
	// If it does, execute it and write the result to the response
	var tmpl *template.Template = t.Lookup(path)
	var contentErr error
	if path != "/" {
		if tmpl != nil {
			var buf bytes.Buffer
			// Write it to a buffer to execute afterwards
			contentErr = tmpl.Execute(&buf, contentData)
			templateContent = buf.String()
		} else {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "./static/errors/404.html")
			return
		}
	}

	data := &types.TemplateData{
		Data: map[string]interface{}{
			"Path":     path,
			"HTML":     template.HTML(templateContent),
			"PageData": pageData,
		},
	}

	// Execute the layout template with the data
	// If the layout template is not found, use the default layout
	// If the layout template is found, execute it with the data
	layout := t.Lookup("layout/" + path)
	pageBuf := app.bufferPool.Get().(*bytes.Buffer)
	pageBuf.Reset()
	defer app.bufferPool.Put(pageBuf)

	// ---------------- THIS IS THE ERROR CATCHING PART ----------------
	// If the contentErr is nil, execute the layout template
	// If the layout template is nil, execute the default layout template
	// If the layout template is not found, use the default layout
	// If the layout template is found, execute it with the data
	// If the contentErr is not nil, write the error to the response
	var err error
	if contentErr == nil {
		if layout == nil {
			err = t.ExecuteTemplate(pageBuf, "layout/app", data)
		} else {
			err = layout.Execute(pageBuf, data)
		}
		if err != nil {
			app.Logger.Error("Error rendering page", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = pageBuf.WriteTo(w)
	}
	if contentErr != nil || err != nil {
		app.Logger.Error("Failed to write template to response: ", "error", err, "content_err", contentErr)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "./static/errors/500.html")
	}
}
