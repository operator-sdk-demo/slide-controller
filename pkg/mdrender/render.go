package mdrender

import (
	"html/template"
	"log"
	"strings"

	presentationsv1alpha1 "github.com/operator-sdk-demo/slide-controller/api/v1alpha1"
)

const mdTemplate = `
{{ range .Slides }}
### {{ .Title }}

{{ range .Bullets }}
- {{ . }}
{{ end }}

{{ range .Images }}
![]({{ . }})
{{ end }}

---
{{ end }}
`

func RenderMarkdown(presentation *presentationsv1alpha1.PresentationSpec) string {
	tmpl, err := template.New("presentation").Parse(mdTemplate)
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
