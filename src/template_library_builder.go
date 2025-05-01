package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"html/template"
)

func getTemplates() (*template.Template, error) {
	log.Println("Parsing 'base' templates...")
	templates, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		return nil, err
	}
	log.Println("Parsing 'partial' templates...")
	templates, err = templates.Funcs(funcMap).ParseGlob("templates/partials/**/*.tmpl")
	if err != nil {
		return nil, err
	}
	log.Println("Parsing 'route' templates...")
	err = filepath.Walk("templates/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			var tmpl *template.Template
			tmpl, err = templates.Funcs(funcMap).ParseFiles(path)
			if err != nil {
				return err
			}
			templates = tmpl
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
