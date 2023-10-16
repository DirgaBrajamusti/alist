package teldrive

import (
	"html/template"
	"net/url"
)

func (d *Teldrive) int64min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (d *Teldrive) sanitizeHTMLURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	sanitizedURL := &url.URL{
		Scheme:   template.HTMLEscapeString(parsedURL.Scheme),
		Opaque:   template.HTMLEscapeString(parsedURL.Opaque),
		User:     parsedURL.User,
		Host:     template.HTMLEscapeString(parsedURL.Host),
		Path:     template.HTMLEscapeString(parsedURL.Path),
		RawQuery: template.HTMLEscapeString(parsedURL.RawQuery),
		Fragment: template.HTMLEscapeString(parsedURL.Fragment),
	}

	return sanitizedURL.String(), nil
}
