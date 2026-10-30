package main

import (
	"os"

	"github.com/cogikyo/hgmx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "hgmx",
	Short:   "A component management tool for Go Templ+HTMX projects.",
	Version: hgmx.Version(),
}

var logLevel string
var (
	paletteFormat string
	paletteOut    string
)

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
	paletteCobraCmd.Flags().StringVarP(&paletteFormat, "format", "f", "tailwind", "Output format [tailwind, css, lua]")
	paletteCobraCmd.Flags().StringVarP(&paletteOut, "out", "o", "", "Write to this file instead of stdout")
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
		if code := initCmd(args); code != 0 {
			os.Exit(code)
		}
	},
}

var paletteCobraCmd = &cobra.Command{
	Use:   "palette <name>=<base>,<primary>[,<accent>]...",
	Short: "Generates semantic palette themes from a base and a primary colour",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if code := paletteCmd(args); code != 0 {
			os.Exit(code)
		}
	},
}
