package main

import (
	"fmt"
	"github.com/ayumu83s/go-mysqldiff/mysqldiff"
)

func main() {

	localConfig, srcConfig, dstConfig := mysqldiff.InitializeConfig()
	//fmt.Printf("[local]host: %s, port: %s, user: %s, password: %s\n",
	//	localDB.DBHost, localDB.DBPort, localDB.DBUser, localDB.DBPassword,
	//)
	//fmt.Printf("[src]host: %s, port: %s, user: %s, password: %s, dbname: %s, filePath: %s\n",
	//	srcConfig.DBHost, srcConfig.DBPort, srcConfig.DBUser, srcConfig.DBPassword, srcConfig.DBName, srcConfig.FilePath,
	//)
	//fmt.Printf("[dst]host: %s, port: %s, user: %s, password: %s, dbname: %s, filePath: %s\n",
	//	dstConfig.DBHost, dstConfig.DBPort, dstConfig.DBUser, dstConfig.DBPassword, dstConfig.DBName, dstConfig.FilePath,
	//)

	src, err := mysqldiff.GetTables(localConfig, srcConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	dst, err := mysqldiff.GetTables(localConfig, dstConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	mysqldiff.Diff(src, dst)
}
