package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

var funcMap = template.FuncMap{
	"CapitalizeFirst": CapitalizeFirst,
}

func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToLower(s[1:])
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// RedirectToSafeURL redirects to a safe URL after authentication
// If the next URL is provided, it will redirect to that URL if it's safe
// Otherwise, it will redirect to the default URL
func RedirectToSafeURL(w http.ResponseWriter, r *http.Request, defaultURL string, logger *slog.Logger) {
	next := r.URL.Query().Get("next")
	if next != "" {
		// Parse the URL to ensure it's not an open redirect vulnerability
		parsedURL, err := url.Parse(next)
		if err != nil || parsedURL.Host != "" {
			// If the URL is invalid or has a host (external URL), redirect to the default URL
			logger.Warn("Potentially unsafe redirect detected", "url", next)
			http.Redirect(w, r, defaultURL, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}
	// If no next parameter is provided, redirect to the default URL
	http.Redirect(w, r, defaultURL, http.StatusSeeOther)
}
