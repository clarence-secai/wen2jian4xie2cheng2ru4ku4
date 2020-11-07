//todo:删掉tag.txt文件后，运行起来就可以重新进行一遍文件分类和数据入库。
// 不知为什么数据入库总是会出现个别数据入错数据库表的情况？？？？？

package main

import (
	"fmt"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/material"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/shu4ju4ru4ku4"
	"go9kai1fang2huan3cun2xie2cheng2/20200807kai1fang2xie2cheng2/wen2jian4fen1lei4"
	"strconv"
	"sync"
	"time"
)

func init() {
	if !wen2jian4fen1lei4.Init1() { //本次进行了文件的分类，则表明数据库里的数据可能乱七八糟，需再次数据入库
		shu4ju4ru4ku4.Drop()//删除数据库内所有的表（因为数据可能乱七八糟）
		shu4ju4ru4ku4.Init2()
	} else { //本次之前已进行了文件分类，现将各个文件内容入库
		shu4ju4ru4ku4.Init2()
	}
	fmt.Println("一分十一，初始化和入库完毕")
	fmt.Println("init运行结束")
}

func main() {
	var input string
	for {
		fmt.Println("请输入要查询的人名[输入exit退出查询,输入show查看缓存]")
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}
		if input == "show" {
			for k, v := range store {
				fmt.Print(k, "---")
				for _, m := range v.Res {
					fmt.Print(m)
				}
				fmt.Println()
			}
			continue
		}
		Explore(input)
	}

}

var ns int
var lk sync.Mutex
var wtgp sync.WaitGroup

func noRes() {
	lk.Lock()
	ns++
	lk.Unlock()
}

//建立缓存的map
type ResAndTag struct {
	Res []shu4ju4ru4ku4.Person
	t   int64 //最新一次查询的时间
	n   int64 //被客户查询的次数
}

var store = make(map[string]*ResAndTag)

//到缓存中查
func lookUp(input string) bool {
	if v, ok := store[input]; ok {
		fmt.Println("查询结果为：", v.Res)
		store[input].n++
		store[input].t = time.Now().Unix()
		return true
	} else {
		return false
	}
}

func Explore(input string) {
	//db := material.GetDB()//这里的db将被十一个协程争夺使用
	//defer db.Close()

	//先到缓存中查询
	if lookUp(input) {
		return
	}

	//否则进入数据库查询并加入缓存中
	resultChan := make(chan shu4ju4ru4ku4.Person, 5)
	signalChan := make(chan int, 11)

//todo:分头干--合作完  各个协程同时分别对数据库中的各个表进行查询
	for i := 0; i <= 10; i++ {
		wtgp.Add(1)
		go func(i int) {
			defer wtgp.Done() //defer会在return前一脚运行
			db := material.GetDB()
			defer db.Close() //"select * from kfinfo where name like ?"
			sql := "select name,idCard,hotel from txt_" + strconv.Itoa(i) + " where name = ?"
			var res shu4ju4ru4ku4.Person
			rows, err := db.Query(sql, input) //,"%"+input+"%"
			if err != nil {
				fmt.Println(i, "表中查询出错:", err)
				signalChan <- 1
				return
			}
			if rows == nil {
				fmt.Println("i,表中没有你要查询的内容")
				noRes()  //没找到结果的协程去 ns++ 一下
				signalChan <- 1
				return
			}
			for rows.Next() { //todo:如果sql语句里是select * from ，则下面必须还有&res.Id
				rows.Scan(&res.Name, &res.IdCard, &res.Hotel) //必须是逐个赋值,顺序也最好和数据库表一致，不能直接只 &res
				fmt.Println("res=", res)
				resultChan <- res //todo:这里如果传指针，担心后面从管道里取出来再指向原本的值时，原本的值被覆盖了
			}
			signalChan <- 1   //每个协程，找没找到结束了都要报个到，
							 // 其实也不用，直接wtgp.Done，然后close(resultChan)
		}(i)
	}
	wtgp.Add(1)
	go func() {
		for i := 0; i <= 10; i++ {
			<-signalChan
		}
		close(resultChan)
		wtgp.Done()
	}()
	//wtgp.Wait()

	//加入缓存中
	store[input] = &ResAndTag{}
	//store[input].Res = make([]*shu4ju4ru4ku4.Person,0)
	for v := range resultChan {
		//if v==nil{
		//	continue
		//}
		store[input].Res = append(store[input].Res, v)
	}
	if ns == 11 {
		fmt.Println("数据库无你要查的内容")
		delete(store, input)  //不可少，因为键值对的值是空，这种键值对也需删掉
		return
	}
	store[input].n++
	store[input].t = time.Now().Unix()

	fmt.Print("查询结果为----")
	for _, v := range store[input].Res {
		fmt.Print(v, "  ")
	}
	fmt.Println()
	//多余5条记录就对缓存进行老旧清除
	delOld()

}

//多余5条记录就对缓存进行老旧清除
func delOld() {
	if len(store) > 5 {
		mix := time.Now().Unix()
		var myK string
		for k, v := range store {
			if v.t < mix {
				myK = k
				mix = v.t
			}
		}
		delete(store, myK)
		fmt.Println("删除了一条旧缓存")
	}
}
