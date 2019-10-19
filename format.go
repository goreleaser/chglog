package chglog

import (
	"bytes"
	"text/template"
)

// FormatChangelog format pkgLogs from a text/template
func FormatChangelog(pkgLogs *PackageChangeLog, tpl *template.Template) (string, error) {
	var data bytes.Buffer
	err := tpl.Execute(&data, pkgLogs)

	return data.String(), err
}
