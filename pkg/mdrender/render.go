package mdrender

import (
	"html/template"
	"log"
	"strings"
)

// Slide represents a single slide in the presentation
type Slide struct {
	Title   string   `json:"title"`
	Bullets []string `json:"bullets"`
	Image   string   `json:"image,omitempty"`
}

// Presentation holds all slides
type Presentation struct {
	Slides []Slide `json:"slides"`
}

const TEMPLATE = `
{{ range .Slides }}
### {{ .Title }}

{{ range .Bullets }}
- {{ . }}
{{ end }}

{{ if .Image }}
![]({{ .Image }})
{{ end }}

---
{{ end }}
`

func RenderMarkdown(presentation *Presentation) string {
	// Parse the template
	tmpl, err := template.New("presentation").Parse(TEMPLATE)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var buf strings.Builder

	// Execute the template and write to the buffer
	if err := tmpl.Execute(&buf, presentation); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// Return the resulting string
	return buf.String()
}
