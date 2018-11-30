package cmd

import (
	"github.com/spf13/cobra"
	"my_github/ast2template/codegen"
	"my_github/ast2template/event/dbmodel"
)

var (
	cmdGenModelFlagDatabase  string
	cmdGenModelFlagTableName string
)

var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "generate gorm db model method",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdGenModelFlagDatabase == "" {
			panic("database must be defined")
		}

		for _, structName := range args {
			generator := dbmodel.DBFuncGenerator{}
			generator.StructName = structName
			generator.Database = cmdGenModelFlagDatabase
			generator.TableName = cmdGenModelFlagTableName

			codegen.Generate(&generator)
		}
	},
}

func init() {
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagDatabase, "database", "", "", "(required) register model to database var")
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagTableName, "table-name", "t", "", "custom table name")

	cmdGen.AddCommand(cmdGenModel)
}
