package main

// contain codes for manipulating the db outside without using the server or the game
import (
	"elichika/exec/db"
	// "elichika/exec/db/account"
	// "elichika/exec/db/gacha"

	"fmt"

	"os"
)

func main() {
	switch os.Args[1] {
	case "init":
		db.Init()
	case "account":
		db.Account()
	case "gacha":
		db.Gacha()
	default:
		fmt.Println("Invalid params:", os.Args)
		return
	}
}
