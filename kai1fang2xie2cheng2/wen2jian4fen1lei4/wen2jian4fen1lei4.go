package wen2jian4fen1lei4

import (
	"fmt"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/material"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

//待写入十个文件的结构体
type PersonFile struct {
	Num      int
	FileChan chan string
	File     *os.File
}

var wg sync.WaitGroup
var (
	str  string
	errs error
)
//通过原始数据分类过程创建的tag.txt，来判断文件分类是否准确的完成
func check() bool {
	fileInfo, err2 := os.Stat("./kai1fang2txt/tag.txt")
	if err2 != nil || fileInfo == nil { //标记文件尚不存在
		fmt.Println("原始数据即将开始分类")
		return false //表明需要进行文件分类操作
	} else {
		bytes, err := ioutil.ReadFile("./kai1fang2txt/tag.txt")
		material.ExitErr(err, "42行")
		tagSlice := strings.Split(strings.TrimSpace(string(bytes)), "\n")
		fmt.Println("tagSlice=", len(tagSlice))
		if len(tagSlice) == material.TxtCount+1 {
			fmt.Println("原始数据本已分类")

			return true //表明原始数据分类进11个文件的操作本就已经完成
		}
		return false //表明原始数据分表不完全成功，需要重新再来一次(向文件写入内容，无append情况下默认覆盖原内容)
	}
}

func Init1() bool {
	//先判断原始文件的分类是否已经完成
	if check() {
		return true //表明本已做好11个分类文件 无需操作
	}
	//否则开始分割原始文件
	//为十个文件和相应管道建立map，方便原始数据找到各自该前往的分类文件予以写入
	fileMap := make(map[int]*PersonFile)
	for i := 0; i <= material.TxtCount; i++ {
		file, err := os.OpenFile("./kai1fang2txt/"+strconv.Itoa(i)+".txt", os.O_CREATE|os.O_RDWR, 0644)
		material.ExitErr(err, "48行")
		//defer file.Close()  //todo:需要避免在for循环里写defer
		fileMap[i] = &PersonFile{
			File: file,
		}
		fileMap[i].Num = i
		fileMap[i].FileChan = make(chan string, 1000)
	}

	//将十一个管道里的数据写进对应的文件里
	//这一部分必须在读取原始文件代码的前面，不然没法让写入的协程在读取原始文件的同时就开始工作
	for _, v := range fileMap {
		wg.Add(1)
		go func(pf *PersonFile) {
			writeIntoTxt(pf)
			wg.Done()
		}(v)
	}

	//对40万条开房原始数据文件的数据分类写进十一个管道里
	writeIntoChan(fileMap)

	wg.Wait()
	fmt.Println("文件一分为十完毕")
	return false //表示本次操作了文件分类
}


