package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

// HighlightQuery highlights query words in text
func HighlightQuery(query, text string) template.HTML {
	words := strings.Fields(query)
	for _, word := range words {
		word = regexp.QuoteMeta(word) // Escape special characters
		re := regexp.MustCompile(`(?i)\b` + word + `\b`) // Case-insensitive match
		text = re.ReplaceAllString(text, `<span class="highlight">$0</span>`)
	}
	return template.HTML(text)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results := []struct {
		Title   string
		Excerpt string
	}{
		{"Sample Title", "This is a sample excerpt"},
		{"Another Title", "Another example of a result"},
	}

	tmpl := template.Must(template.New("results").Funcs(template.FuncMap{
		"highlight": HighlightQuery,
	}).ParseFiles("results.html"))

	tmpl.Execute(w, map[string]interface{}{
		"Query":   query,
		"Results": results,
	})
}
