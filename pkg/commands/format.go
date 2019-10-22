package commands

import (
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/goreleaser/chglog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nolint: gocognit, funlen
func setupFormatCmd(config *viper.Viper) (cmd *cobra.Command) {
	var (
		input,
		output,
		pkg,
		templateName,
		templateFile string
	)
	cmd = &cobra.Command{
		Use:   "format",
		Short: "display version info",
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
	config.BindPFlag("package-name", cmd.Flags().Lookup("package-name"))

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

	cmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(c, args)
	}

	cmd.RunE = func(c *cobra.Command, args []string) (err error) {
		var (
			tpl        *template.Template
			data       []byte
			ret        string
			fmtPackage = new(chglog.PackageChangeLog)
		)
		if templateName != "" && templateFile != "" {
			return fmt.Errorf("--template and --template-file are mutually exclusive")
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
			if data, err = ioutil.ReadFile(templateFile); err != nil { // nolint: gosec
				return err
			}
			tpl, err = chglog.LoadTemplateData(string(data))
		}
		if err != nil {
			return err
		}

		fmtPackage.Name = config.GetString("package-name")

		if fmtPackage.Entries, err = chglog.Parse(input); err != nil {
			return err
		}

		if ret, err = chglog.FormatChangelog(fmtPackage, tpl); err != nil {
			return err
		}

		if output == "-" {
			fmt.Println(ret)
			return
		}

		return ioutil.WriteFile(output, []byte(ret), 0644)
	}

	return cmd
}
