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
	paletteLua bool
	skip       []string
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
	rootCmd.AddCommand(bufCobraCmd)
	rootCmd.AddCommand(updateCobraCmd)
	rootCmd.AddCommand(skillsCobraCmd)
	rootCmd.AddCommand(migrateCobraCmd)
	rootCmd.AddCommand(credentialsCobraCmd)
	credentialsCobraCmd.AddCommand(credentialsCheckCobraCmd)
	for _, c := range []*cobra.Command{genCobraCmd, checkCobraCmd} {
		c.Flags().StringSliceVar(&skip, "skip", nil, "Native platforms to leave unchecked and ungenerated: android, ios")
	}
	paletteCobraCmd.Flags().BoolVar(&paletteLua, "lua", false, "Write a Lua table of hex shades to stdout instead of CSS")
}

var genCobraCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generates templ Go code, the minified Tailwind CSS and the protobuf contract code",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(genCmd(skip))
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
	Short: "Runs the contract, generation, pin and credentials checks, site rules, format check, vet, build and tests",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(checkCmd(skip))
	},
}

var bufCobraCmd = &cobra.Command{
	Use:                "buf [buf args...]",
	Short:              "Runs hgmx's pinned buf in the current directory",
	Hidden:             true,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(bufPassCmd(args))
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

var credentialsCobraCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Writes credentials.json from the site's credential declarations",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(credentialsCmd())
	},
}

var credentialsCheckCobraCmd = &cobra.Command{
	Use:   "check [name...] [-- site flags...]",
	Short: "Tests the site's SMTP and Google credentials without printing them",
	Run: func(cmd *cobra.Command, args []string) {
		names, flags := args, []string(nil)
		if i := cmd.ArgsLenAtDash(); i >= 0 {
			names, flags = args[:i], args[i:]
		}
		os.Exit(credentialsCheckCmd(names, flags))
	},
}

var updateCobraCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the site's hgmx library, tool and agent skills to master",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(updateCmd())
	},
}

var skillsCobraCmd = &cobra.Command{
	Use:   "skills",
	Short: "Installs hgmx's agent skills into .agents/skills and links them for Claude Code",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(skillsCmd())
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
	Use:   "palette <hex_color>",
	Short: "Generates a color palette based on the input hex color",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if code := paletteCmd(args); code != 0 {
			os.Exit(code)
		}
	},
}
