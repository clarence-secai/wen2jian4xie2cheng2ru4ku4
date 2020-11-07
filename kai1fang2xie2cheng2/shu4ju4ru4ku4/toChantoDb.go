package shu4ju4ru4ku4

import (
	"bufio"
	"database/sql"
	"fmt"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/material"
	"io"
	"os"
	"strconv"
	"strings"
)

//todo 生产者  一个生产协程从分类小文件读取数据写入一个管道 //todo:改进：如何多协程同时有序读同一个文件？
func readTxt2chan(i int) {
	//fmt.Println("开始读文件",i)
	file, err := os.Open("./kai1fang2txt/" + strconv.Itoa(i) + ".txt")
	material.ExitErr(err, "16行")
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("文件%v读取完毕\n", i)
			close(tenChanMap[i].channel)
			fmt.Printf("文件%v管道关闭了", i)
			//tenChanMap[i]=Achan{ChannelCondition:true}
			//todo:上面一行不能要，否则相当于不要管道了；因为大括号赋值是一次性赋值，不赋值的字段是给默认值 go0ji1chu3/18struct/01struct.go
			break
		}
		if len(strings.Split(str, ",")) != 3 {
			fmt.Printf("文件%v末尾有格式，实际已读完\n", i)
			close(tenChanMap[i].channel)
			fmt.Println("管道关闭了")
			//tenChanMap[i]=Achan{ChannelCondition:true}
			//todo:上面一行不能要，否则相当于不要管道了；因为大括号赋值是一次性赋值，不赋值的字段是给默认值 go0ji1chu3/18struct/01struct.go
			break
		}
		if err != nil {
			fmt.Printf("读取文件%v一行失败\n", i)
			continue
		}
		strSlice = strings.Split(str, ",")
		name := strings.TrimSpace(strSlice[0])
		idCard := strings.TrimSpace(strSlice[1])
		hotel := strings.TrimSpace(strSlice[2])
		person = &Person{Name: name, IdCard: idCard, Hotel: hotel}
		if !strings.Contains(person.Name, "jack") {
			fmt.Println(i, "号文件有一条无jack")
			continue
		}
		//todo:向对应的管道中放进数据
		tenChanMap[i].channel <- person
		//fmt.Println("分类文件",i,"向对应管道放进数据~")
	}

}

//todo 消费者,有4个消费协程从一个管道读取输入后写入数据库
func fromChan2Db(i int) {
	//fmt.Println("创建相应数据表和积累100条")
	//创建相应的数据库中的数据表
	//todo:每一个写数据进数据库的协程都有一个自己的mydb,本例中共44个,避免db被多个协程轮着用而挂掉
	mydb, err := sql.Open("mysql", "root:413188ok@tcp(localhost:3306)/kai1fang2")
	material.ExitErr(err, "64行")
	defer mydb.Close()
	//todo:一定要加 if not exists
	_, err = mydb.Exec(`create table if not exists txt_` + strconv.Itoa(i) + `(
											id int auto_increment primary key,
											name varchar(20) not null,
											idCard int not null,
											hotel varchar(20) not null);`)
	//fmt.Println("result=",result)
	material.ExitErr(err, "66行")

	//var hundredSlice = make([]*Person,99)
	var hundredSlice = make([]*Person, 0)
	//var r int
	//for v := range tenChanMap[i].channel{
	for {
		v, ok := <-tenChanMap[i].channel
		// fmt.Printf("从第%v个管道读出%v入库~\n",i,v)
		if !ok { //把最后一次不足99条的记录写进数据库
			//if hundredSlice[0]==nil {
			//	break
			//}
			//hundredSlice2 := hundredSlice[:r]
			if len(hundredSlice) > 0 { //todo:避免刚好99个的都批量录入了，没有零头而把空切片也录入一回

				Insert2Db(mydb, hundredSlice, i)//将不够100的也累计数据后批量插入
			}
			// cleanHundredSlice(hundredSlice)
			break
		}
		//hundredSlice[r]=v
		//r++
		hundredSlice = append(hundredSlice, v)
		if len(hundredSlice) == 100 { //tenChanMap[i].ChannelCondition用这个不太好
			//todo:每一百条记录插入一次数据库表，从而节省时间提高效率，因为操作一次数据库较耗时，论次不论量
			Insert2Db(mydb, hundredSlice, i)
			//cleanHundredSlice(hundredSlice)
			hundredSlice = make([]*Person, 0) //todo:清空切片，重新来过
			//r = 0
		}
	}
	//}
}

var accumulate int //计算有多个个100条或不足100条的批次记录入库，也即对数据库实行写操作的次数
func Insert2Db(mydb *sql.DB, hundredSliceNeedle []*Person, i int) {
	//fmt.Println("开始插入数据库")
	sqls := "insert into txt_" + strconv.Itoa(i) + "(name,idCard,hotel)values"
	var sqlStr string
	for _, p := range hundredSliceNeedle {
		sqlStr += `,("` + p.Name + `","` + p.IdCard + `","` + p.Hotel + `")`
	}
	sqls += sqlStr[1:] //去掉第一个逗号
	//for _,p := range *hundredSliceNeedle {
	//	sqlStr += "("+ p.Name +","+p.IdCard+","+p.Hotel+"),"
	//}
	//fmt.Println("sqls=",sqls)
	//sql += sqlStr[:len(sqlStr)-1]
	//fmt.Println("执行插入数据表")
	mydb.Exec(sqls)
	accumulate++
	fmt.Println("一批100条插入数据表完成成功了！", accumulate)
}