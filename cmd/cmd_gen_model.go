package cmd

import (
	"github.com/spf13/cobra"
	"my_github/ast2template/codegen"
	"my_github/ast2template/event/dbmodel"
)

var (
	cmdGenModelFlagTableName string
)

var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "generate gorm db model method",
	Run: func(cmd *cobra.Command, args []string) {
		for _, structName := range args {
			generator := dbmodel.DBFuncGenerator{}
			generator.StructName = structName
			generator.TableName = cmdGenModelFlagTableName

			codegen.Generate(&generator)
		}
	},
}

func init() {
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagTableName, "table-name", "t", "", "custom table name")

	cmdGen.AddCommand(cmdGenModel)
}
