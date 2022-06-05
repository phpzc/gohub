package cmd

import (
	"gohub/database/migrations"
	"gohub/pkg/migrate"

	"github.com/spf13/cobra"
)

var CmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	//所有 migrate 下的子命令都会执行 以下代码
}

//执行迁移文件的命令
var CmdMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Run unmigrated migrations",
	Run:   runUp,
}

func init() {
	CmdMigrate.AddCommand(
		CmdMigrateUp,
	)
}

func migrator() *migrate.Migrator {
	//注册 database/migrations下的所有迁移文件
	migrations.Initialize()
	//初始化 migrator
	return migrate.NewMigrator()
}

func runUp(cmd *cobra.Command, args []string) {
	migrator().Up()
}

func runMigrate(cmd *cobra.Command, args []string) {

}
