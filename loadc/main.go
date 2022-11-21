package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

const (
	EnvPWD       = "PWD"
	EnvGoPackage = "GOPACKAGE"
	EnvGoFile    = "GOFILE"
	EnvGoLine    = "GOLINE"

	ModeSql  = "sql"
	ModeApi  = "api"
	ModeSqlx = "sqlx"

	FileMode = 0b_0110_0100_0100 // 0644
)

var (
	PackageName = os.Getenv(EnvGoPackage)
	CurrentDir  = os.Getenv(EnvPWD)
	CurrentFile = os.Getenv(EnvGoFile)
	LineNum, _  = strconv.Atoi(os.Getenv(EnvGoLine))

	FileContent []byte
)

var (
	mode     string
	features []string
	output   string
	pointer  bool
)

var loadc = &cobra.Command{
	Use:     "loadc",
	Version: "v0.4.4",
	Args: func(cmd *cobra.Command, args []string) error {
		switch mode {
		case ModeSql:
			return cobra.ExactArgs(1)(cmd, args)
		case ModeApi, ModeSqlx:
			return cobra.NoArgs(cmd, args)
		default:
			return nil
		}
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkFeatures(features); err != nil {
			return err
		}

		switch mode {
		case ModeSql:
			return genSql(cmd, args)
		case ModeApi:
			return genApi(cmd, args)
		case ModeSqlx:
			return genSqlx(cmd, args)
		default:
			return nil
		}
	},
}

var validFeatures = []string{
	FeatureApiCache,
	FeatureApiLog,
	FeatureApiClient,
	FeatureSqlxLog,
}

func checkFeatures(features []string) error {
	if len(features) == 0 {
		return nil
	}

Check:
	for _, feature := range features {
		for _, valid := range validFeatures {
			if feature == valid {
				continue Check
			}
		}

		return fmt.Errorf("checkFeatures: invalid feature %s, available features are: \n\n%s\n\n",
			quote(feature),
			printStrings(features))
	}

	return nil
}

func init() {
	loadc.Flags().StringVarP(&mode, "mode", "m", "", "mode=[sql, api, sqlx, ...]")
	loadc.Flags().StringSliceVarP(&features, "features", "f", nil, "features")
	loadc.Flags().StringVarP(&output, "output", "o", "", "output file name")
	loadc.Flags().BoolVar(&pointer, "pointer", false, "mode=sql: make 'SqlLoader' pointer type (*ident)")
}

func init() {
	var err error
	FileContent, err = read(join(CurrentDir, CurrentFile))
	cobra.CheckErr(err)
}

func main() {
	cobra.CheckErr(loadc.Execute())
}
