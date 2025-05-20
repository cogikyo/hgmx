package main

import (
	"os"

	"github.com/nosvagor/hgmx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "hgmx",
	Short:   "A component management tool for Go Templ+HTMX projects.",
	Version: hgmx.Version(),
}

var logLevel string
var linkInput string
var linkOutput string

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set log verbosity level [debug, info, warn, error]")
	rootCmd.AddCommand(infoCobraCmd)
	rootCmd.AddCommand(initCobraCmd)
	rootCmd.AddCommand(paletteCobraCmd)
	rootCmd.AddCommand(linkCobraCmd)
	linkCobraCmd.Flags().StringVarP(&linkInput, "input", "i", "../hgmx/library/*", "Source directory to link from")
	linkCobraCmd.Flags().StringVarP(&linkOutput, "output", "o", "app/*", "Transforms files in directory to symlinks")
}

var infoCobraCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about the hgmx environment",
	Run: func(cmd *cobra.Command, args []string) {
		infoCmd()
	},
}

var initCobraCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new hgmx project",
	Run: func(cmd *cobra.Command, args []string) {
		initCmd(args)
	},
}

var paletteCobraCmd = &cobra.Command{
	Use:   "palette <hex_color>",
	Short: "Generates a color palette based on the input hex color",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paletteCmd(args)
	},
}

var linkCobraCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlinks files in the output directory to the source directory",
	Run: func(cmd *cobra.Command, args []string) {
		linkCmd(linkInput, linkOutput)
	},
}
