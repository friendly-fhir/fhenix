package cmd

import (
	"fmt"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/template"
	"github.com/friendly-fhir/fhenix/model"
	"github.com/friendly-fhir/fhenix/model/raw"
	"github.com/spf13/cobra"
)

var Evaluate = &cobra.Command{
	Use:     "evaluate <template>",
	Aliases: []string{"eval"},
	Short:   "Evaluate a template string",
	Long: strings.Join([]string{
		"Evaluate a template string and print the result to stdout using test data.",
		"This can be used to evaluate whether a template works as expected.",
	}, "\n"),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}
		content := args[0]
		tmpl, err := template.Parse("evaluate", content)
		if err != nil {
			return err
		}
		t := &model.Type{
			Source: &model.TypeSource{
				Package:             fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
				File:                "StructureDefinition-string.json",
				StructureDefinition: &raw.StructureDefinition{},
			},
			Kind:       "primitive-type",
			Name:       "String",
			Base:       nil,
			IsAbstract: true,
			URL:        "http://hl7.org/fhir/StructureDefinition/string",
		}
		err = tmpl.Execute(cmd.OutOrStdout(), t)
		fmt.Fprintln(cmd.OutOrStdout())
		return err
	},
}

func init() {
	Root.AddCommand(Evaluate)
}
