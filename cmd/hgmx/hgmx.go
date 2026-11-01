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
	rootCmd.AddCommand(genCobraCmd)
	rootCmd.AddCommand(fmtCobraCmd)
	rootCmd.AddCommand(devCobraCmd)
	rootCmd.AddCommand(checkCobraCmd)
	rootCmd.AddCommand(updateCobraCmd)
	rootCmd.AddCommand(migrateCobraCmd)
	paletteCobraCmd.Flags().StringVarP(&paletteFormat, "format", "f", "tailwind", "Output format [tailwind, css, lua]")
	paletteCobraCmd.Flags().StringVarP(&paletteOut, "out", "o", "", "Write to this file instead of stdout")
}

var genCobraCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generates templ Go code and the minified Tailwind CSS",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(genCmd())
	},
}

var fmtCobraCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Formats the site's templ and Go files",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(fmtCmd())
	},
}

var devCobraCmd = &cobra.Command{
	Use:   "dev [-- site flags...]",
	Short: "Regenerates on change and restarts the site",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(devCmd(args))
	},
}

var checkCobraCmd = &cobra.Command{
	Use:   "check",
	Short: "Runs the generation check, site rules, format check, vet, build and tests",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(checkCmd())
	},
}

var migrateCobraCmd = &cobra.Command{
	Use:                "migrate <args>...",
	Short:              "Runs the site's migrate command",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(migrateCmd(args))
	},
}

var updateCobraCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the site's hgmx library and tool to master, tidies and vendors",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(updateCmd())
	},
}

var infoCobraCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about the hgmx environment",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(infoCmd())
	},
}

var initCobraCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new hgmx project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(initCmd())
	},
}

var paletteCobraCmd = &cobra.Command{
	Use:   "palette <name>=<base>,<primary>[,<accent>]...",
	Short: "Generates semantic palette themes from a base and a primary colour",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(paletteCmd(args))
	},
}
