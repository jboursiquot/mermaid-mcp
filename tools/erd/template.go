package erd

import "text/template"

var mermaidERDTemplate = template.Must(template.New("erd").Parse(`
erDiagram
	direction {{ .Direction }}
	{{- range .Tables }}
	{{ .Name }} {
	{{- range .Columns }}
		{{ .Type }} {{ .Name }}
	{{- end }}
	}
	{{- end }}
	
	{{- range .ForeignKeys }}
	{{ .Table }} }o--|| {{ .ForeignTable }} : "{{ .Column }} ➝ {{ .ForeignColumn }}"
	{{- end }}
	`))
