package chglog

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

const (
	rpmTpl = `
{{- range .Entries }}{{$version := semver .Semver}}
* {{ .Date | date "Mon Jan 2 2006" }} {{ .Packager }} - {{ $version.Major }}.{{ $version.Minor }}.{{ $version.Patch }}{{if $version.Prerelease}}-{{ $version.Prerelease }}{{end}}
{{- range .Changes }}
- {{ .Note }} 
{{- end }}
{{- end }}
`
	debTpl = `{{$name := .Name}}
{{- range .Entries }}
{{ $name }} ({{ .Semver }}) {{if .Deb}}{{default "" (.Deb.Distributions | join " ")}}; urgency={{default "low" .Deb.Urgency}}{{end}}
{{range .Changes }}{{$note := splitList "\n" .Note}}
  * {{ first $note }}
  {{ range $i,$n := (rest $note) }}- {{$n}}
  {{end}}
{{end}}

-- {{ .Packager }} {{ .Date | date "Mon, 2 Jan 2006 03:04:05 -0700" }}
{{- end }}
`
	repoTpl = `
repo: {{ .Name }}
`
)

// LoadTemplateData load a template from string with all of the sprig.TxtFuncMap loaded
func LoadTemplateData(data string) (*template.Template, error) {
	return template.New("base").Funcs(sprig.TxtFuncMap()).Parse(data)
}
