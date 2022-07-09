package myCmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/stoffand/go-validation/generator"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vgen",
		Short: "cli tool to generate validation logic",
		Long:  "generate validation logic from exsiting go struct",
	}
	genVerbose bool
	genCmd     = &cobra.Command{
		Use:  "generate",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			generator.Generate(args, genVerbose)
		},
	}
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().BoolVarP(&genVerbose, "verbose", "v", false, "display extra information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
