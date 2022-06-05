package seeders

import "gohub/pkg/seed"

func Initialize() {

	//触发本目录下其他文件中的init方法

	//指定优先于同目录下的其他文件运行
	seed.SetRunOrder([]string{
		"SeedUsersTable",
	})
}
