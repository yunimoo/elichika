package db

import (
	"elichika/serverdb"

	"fmt"
	"os"
)

func Init() {
	// wipe the database if necessary, then insert the relevant table
	if len(os.Args) == 2 {
		serverdb.InitTables(false)
	} else {
		if os.Args[2] == "overwrite" {
			serverdb.InitTables(true)
		} else {
			fmt.Println("Invalid params:", os.Args)
			fmt.Println("Use", os.Args[0], "init [overwrite]")
		}
	}
}
