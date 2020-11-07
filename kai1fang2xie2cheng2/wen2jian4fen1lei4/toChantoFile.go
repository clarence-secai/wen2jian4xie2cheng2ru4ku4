package wen2jian4fen1lei4

import (
	"bufio"
	"fmt"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/material"
	"io"
	"os"
	"strconv"
	"strings"
)
//todo 分发者
func writeIntoChan(fileMap map[int]*PersonFile) {
	file, err := os.Open("./kai1fang2txt/kf.txt")
	material.ExitErr(err, "20行")
	defer file.Close()
	reader := bufio.NewReader(file)
	var n int //记录写入十个文件的开房记录条数
	for {
		str, errs = reader.ReadString('\n')//todo:改进：如何开多个协程接力地读同一个文件？
		if errs != nil {
			fmt.Println("有一行内容读取失败")
			continue
		}
		if errs == io.EOF {
			//要关闭十个文件的管道
			for i := 0; i <= material.TxtCount; i++ {
				close(fileMap[i].FileChan)
			}
			break
		}
		strSlice := strings.Split(str, ",") //这里的这个逗号最好从原本文中复制过来，因为不确定是英文的还是中文的
		if len(strSlice) != 3 {             //原始文件末尾可能有格式，即虽然实际读完了但因为格式之类看不见的东西的存在导致不会io.EOF
			fmt.Println("合法数据已经读完")
			for i := 0; i <= material.TxtCount; i++ {
				close(fileMap[i].FileChan)
			}
			break
		}
		name := strings.TrimSpace(strSlice[0])
		personNum, err := strconv.Atoi(name[4:]) //借此来确定这个人的信息该放进第几号文件
		material.GoonErr(err, "81行")
		//fmt.Println(personNum)
		//todo:方式一：
		//switch personNum/40000 {
		//case 0 :
		//	fileMap[0].FileChan <- str
		//case 1 :
		//	fileMap[1].FileChan <- str
		//case 2 :
		//	fileMap[2].FileChan <- str
		//case 3 :
		//	fileMap[3].FileChan <- str
		//case 4 :
		//	fileMap[4].FileChan <- str
		//case 5 :
		//	fileMap[5].FileChan <- str
		//case 6 :
		//	fileMap[6].FileChan <- str
		//case 7 :
		//	fileMap[7].FileChan <- str
		//case 8 :
		//	fileMap[8].FileChan <- str
		//case 9 :
		//	fileMap[9].FileChan <- str
		//case 10:
		//	fileMap[10].FileChan <- str
		//}
		//todo:方式二：
		if v, ok := fileMap[personNum/40000]; ok { //todo:通过map以键找值
			v.FileChan <- str
		}
		n++
	}
	fmt.Println("读取原始数据条数：", n)
}

//todo 订阅者
func writeIntoTxt(pf *PersonFile) {
	var m int
	for v := range pf.FileChan {
		//_, err := pf.File.WriteString(v+"\n")//todo:坑在这里了！！！！！！72行在读取时回车符也包含在str里
		_, err := pf.File.WriteString(v) //todo:坑在这里了！！！！！！
		material.GoonErr(err, "117行")
		m++
	}
	fmt.Printf("写入文件%v开房记录条数%d：\n", pf.Num, m)
	pf.File.Close() //记得打开的文件要关闭

	//向tag.txt文件记录一下分割进入十一个文件中的某一个文件的完成情况
	tagFile, err := os.OpenFile("./kai1fang2txt/tag.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	material.ExitErr(err, "135行")
	defer tagFile.Close()
	_, err = tagFile.WriteString("文件" + strconv.Itoa(pf.Num) + "完成写入\n")
	material.ExitErr(err, "137行")
}