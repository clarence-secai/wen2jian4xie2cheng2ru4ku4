package shu4ju4ru4ku4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/material"
	"runtime"
	"strconv"
	"sync"
)

var (
	wg         sync.WaitGroup
	tenChanMap map[int]Achan
	strSlice   []string
	person     *Person
)

type Achan struct {
	channel          chan *Person
	ChannelCondition bool
}
type Person struct {
	Name   string
	IdCard string
	Hotel  string
}

func check() bool {
	//db, err := sql.Open("mysql", "root:413188ok!@tcp(localhost:3306)/kai1fang2")
	//material.ExitErr(err,"34行")
	db := material.GetDB()
	defer db.Close()
	var num int
	for i := 0; i <= 10; i++ {
		sql := "select count(*) from txt_" + strconv.Itoa(i)
		var row int
		db.QueryRow(sql).Scan(&row)
		num += row
	}
	fmt.Println("num=", num)
	if num == 405000 {
		fmt.Println("原始数据本已入库,无需再进行入库操作")
		return true //说明原始数据本已入库,无需再进行入库操作
	} else {
		return false //表明之前的入库有问题，需重新进行数据入库操作
	}
}

func Drop() {
	db := material.GetDB()
	for i := 0; i <= 10; i++ {
		db.Exec("drop table if exists txt_" + strconv.Itoa(i))
	}
}

func Init2() bool {
	if check() {
		return true //说明原始数据本已入库,无需再进行入库操作
	}

	strSlice = make([]string, 3)

	//建立十一个管道的map，方便分别从对应的十一个分类文件中读取数据并放进相应的管道
	tenChanMap = make(map[int]Achan)
	for i := 0; i <= material.TxtCount; i++ {
		ch := make(chan *Person, 1000) //todo:...
		//tenChanMap[i].channel = ch //todo:为什么不能这样赋值？//因为Achan没加& 参见go0ji1chu3/17map/06map0struct.go:9
		//tenChanMap[i].ChannelCondition=false
		tenChanMap[i] = Achan{channel: ch, ChannelCondition: false}
	}
	fmt.Printf("创建管道map完成\n")

	//对每个分类文件对应的读取管道开启四条写入数据库的协程
	for i := 0; i <= material.TxtCount; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func(i int) {
				fromChan2Db(i) //todo:数据库不怕多协程同时写入数据表，数据库自带互斥锁
				fmt.Println("done1", i, "号管道读完入库完")
				wg.Done()
			}(i)
		}
	}
	//将每个分类文件的数据写入对应的管道
	for i := 0; i <= material.TxtCount; i++ {
		wg.Add(1)
		go func(i int) {
			readTxt2chan(i)
			fmt.Println("done2", i, "号文件已全部写入管道")
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("十一个文件读取入库完毕\n")
	return false //表明进行过数据入库操作
}


func cleanHundredSlice(slice []*Person) {  //todo:手动回收内存垃圾，一般不用，GC会自动回收垃圾
	for i := 0; i < len(slice); i++ {
		slice[i] = nil
	}
	fmt.Println("GC回收")
	runtime.GC()
}


