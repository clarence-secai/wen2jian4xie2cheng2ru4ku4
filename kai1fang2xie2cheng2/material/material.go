package material

import (
	"database/sql"
	"fmt"
	"os"
)

func ExitErr(err error, where string) {
	if err != nil {
		fmt.Println(where, "--报错--", err)
		os.Exit(1)
	}
}
func GoonErr(err error, where string) {
	if err != nil {
		fmt.Println(where, "--报错--", err)
	}
}

const TxtCount = 10 //设定0-10共十一个文件

func GetDB() *sql.DB {
	mydb, err := sql.Open("mysql", "root:413188ok@tcp(localhost:3306)/kai1fang2")
	ExitErr(err, "64行")
	return mydb //不能在这个函数里defer close(mydb)不然调用该函数取不到db了,但调用者要defer close(mydb)
}
