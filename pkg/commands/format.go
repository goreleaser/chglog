package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/goreleaser/chglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrTemplateFlags occurs when the user tries to use both --template and --template-file flags.
var ErrTemplateFlags = errors.New("--template and --template-file are mutually exclusive")

// nolint: funlen, gocritic
func FormatCmd(config *viper.Viper) (cmd *cobra.Command) {
	var input,
		output,
		pkg,
		templateName,
		templateFile string
	cmd = &cobra.Command{
		Use:   "format",
		Short: "format entries to a specific template",
	}
	cmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"-",
		"file to save the output to (- is stdout)")
	cmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		"changelog.yml",
		"changelog.yml file to use as the basis for formatting")
	cmd.Flags().StringVarP(
		&pkg,
		"package-name",
		"p",
		"",
		"package name to use in formatting")

	cmd.Flags().StringVarP(
		&templateName,
		"template",
		"t",
		"",
		"builtin template to use ('deb', 'rpm', 'release', 'repo')")
	cmd.Flags().StringVarP(
		&templateFile,
		"template-file",
		"",
		"",
		"custom template file to use")

	cmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		return config.BindPFlag("package-name", cmd.Flags().Lookup("package-name"))
	}

	cmd.RunE = func(*cobra.Command, []string) (err error) {
		var (
			tpl        *template.Template
			data       []byte
			ret        string
			fmtPackage = new(chglog.PackageChangeLog)
		)
		if templateName != "" && templateFile != "" {
			return ErrTemplateFlags
		}

		switch strings.ToLower(templateName) {
		case "deb":
			tpl, err = chglog.DebTemplate()
		case "rpm":
			tpl, err = chglog.RPMTemplate()
		case "release":
			tpl, err = chglog.ReleaseTemplate()
		case "repo":
			tpl, err = chglog.RepoTemplate()
		default:
			// nolint: gosec, gocritic
			if data, err = os.ReadFile(templateFile); err != nil {
				return fmt.Errorf("error formatting entries: %w", err)
			}
			tpl, err = chglog.LoadTemplateData(string(data))
		}
		if err != nil {
			return fmt.Errorf("error formatting entries: %w", err)
		}

		fmtPackage.Name = config.GetString("package-name")

		if fmtPackage.Entries, err = chglog.Parse(input); err != nil {
			return fmt.Errorf("error formatting entries: %w", err)
		}

		if ret, err = chglog.FormatChangelog(fmtPackage, tpl); err != nil {
			return fmt.Errorf("error formatting entries: %w", err)
		}

		if output == "-" {
			fmt.Println(ret)

			return
		}

		// nolint: gosec, gocritic
		return os.WriteFile(output, []byte(ret), 0o644)
	}

	return cmd
}
