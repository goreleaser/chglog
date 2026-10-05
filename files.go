package chglog

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Parse parse a changelog.yml into ChangeLogEntries.
func Parse(file string) (entries ChangeLogEntries, err error) {
	var body []byte
	body, err = os.ReadFile(file) // nolint: gosec,gocritic
	switch {
	case os.IsNotExist(err):
		return make(ChangeLogEntries, 0), nil
	case err != nil:
		return nil, fmt.Errorf("error parsing %s: %w", file, err)
	}

	if err = yaml.Unmarshal(body, &entries); err != nil {
		return entries, fmt.Errorf("error parsing %s: %w", file, err)
	}

	return entries, nil
}

// Save save ChangeLogEntries to a yml file.
func (c *ChangeLogEntries) Save(file string) (err error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error saving %s: %w", file, err)
	}
	// nolint: gosec,gocritic
	return os.WriteFile(file, data, 0o644)
}

type (
	changeLogChange ChangeLogChange
	changeLogNotes  ChangeLogNotes
)

// MarshalYAML implements yaml.Marshaler. See quotedIndentedBlocks.
func (c ChangeLogChange) MarshalYAML() (any, error) {
	return quotedIndentedBlocks(changeLogChange(c))
}

// MarshalYAML implements yaml.Marshaler. See quotedIndentedBlocks.
func (n ChangeLogNotes) MarshalYAML() (any, error) {
	return quotedIndentedBlocks(changeLogNotes(n))
}

// quotedIndentedBlocks encodes v into a node and forces double-quoted style
// on multi-line strings that start with a space or a line break.
//
// yaml.v3 writes such strings as literal blocks with an indentation indicator
// (for example "note: |4-"). Inside a "- key:" sequence item it computes the
// block indentation from a different column than the parser does, so the
// changelog it writes either fails to parse ("did not find expected key") or
// reads back with different whitespace. Double-quoted scalars round-trip
// exactly at any nesting level.
func quotedIndentedBlocks(v any) (*yaml.Node, error) {
	var node yaml.Node
	if err := node.Encode(v); err != nil {
		return nil, err
	}
	forceQuotedIndentedBlocks(&node)
	return &node, nil
}

func forceQuotedIndentedBlocks(n *yaml.Node) {
	if n.Kind == yaml.ScalarNode && n.Tag == "!!str" && strings.Contains(n.Value, "\n") &&
		(n.Value[0] == ' ' || n.Value[0] == '\n') {
		n.Style = yaml.DoubleQuotedStyle
	}
	for _, child := range n.Content {
		forceQuotedIndentedBlocks(child)
	}
}
