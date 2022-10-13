package main

import (
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

const (
	EnvPWD       = "PWD"
	EnvGoPackage = "GOPACKAGE"
	EnvGoFile    = "GOFILE"
	EnvGoLine    = "GOLINE"

	ModeSql = "sql"
	ModeApi = "api"

	FileMode = 0b_0110_0100_0100 // 0644
)

var (
	PackageName = os.Getenv(EnvGoPackage)
	CurrentDir  = os.Getenv(EnvPWD)
	CurrentFile = os.Getenv(EnvGoFile)
	LineNum, _  = strconv.Atoi(os.Getenv(EnvGoLine))
)

var (
	mode     string
	features []string
	output   string
	pointer  bool
)

var loadc = &cobra.Command{
	Use: "loadc",
	Args: func(cmd *cobra.Command, args []string) error {
		switch mode {
		case ModeSql:
			return cobra.ExactArgs(1)(cmd, args)
		case ModeApi:
			return cobra.NoArgs(cmd, args)
		default:
			return nil
		}
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch mode {
		case ModeSql:
			return genSql(cmd, args)
		case ModeApi:
			return genApi(cmd, args)
		default:
			return nil
		}
	},
}

func init() {
	loadc.Flags().StringVarP(&mode, "mode", "m", "", "mode=[sql, api, ...]")
	loadc.Flags().StringArrayVarP(&features, "features", "f", nil, "features")
	loadc.Flags().StringVarP(&output, "output", "o", "", "output file name")
	loadc.Flags().BoolVar(&pointer, "pointer", false, "mode=sql: make 'SqlLoader' pointer type (*ident)")
}

func main() {
	cobra.CheckErr(loadc.Execute())
}
