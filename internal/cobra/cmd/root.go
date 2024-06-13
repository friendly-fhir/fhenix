package cmd

import "github.com/spf13/cobra"

var Root = &cobra.Command{
	Use:   "fhenix",
	Short: "Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
}

func init() {
	flags := Root.PersistentFlags()
	flags.StringP("output", "o", "", "The output directory to write the generated code to")
}
