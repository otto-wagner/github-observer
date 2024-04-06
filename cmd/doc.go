package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
)

const defaultOutputDirectory = "./"

func init() {
	docCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", defaultOutputDirectory, "Output directory")
}

var (
	outputDir string
	docCmd    = &cobra.Command{
		Use:   "generate-docs",
		Short: "Generate Markdown Documentation for CLI",
		Long:  `Creates a markdown file for the server command.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := doc.GenMarkdownTree(rootCmd, outputDir)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)
