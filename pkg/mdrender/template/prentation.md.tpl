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
