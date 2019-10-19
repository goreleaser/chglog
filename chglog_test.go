package chglog

func createNote(header, footer string) *ChangeLogNotes {
	if header == "" && footer == "" {
		return nil
	}
	ret := &ChangeLogNotes{}
	if header != "" {
		ret.Header = &header
	}
	if footer != "" {
		ret.Footer = &footer
	}

	return ret
}
