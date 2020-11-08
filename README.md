# wen2jian4xie2cheng2ru4ku4
本项目中是将一个爬虫获取到的存放于大文件的内容，使用管道、协程快速切分成小文件并入数据库，最后供客户端查询。
1、项目用到了如下协程模型
   发布者---订阅者 ： 一个协程读原始大文件，经过判断【判断过程也可以各个管道的接收协程来做】写入若干管道中当前数据该
                     去的管道，若干管道里每个管道都有若干协程争夺读取并写入小文件
   生产者---消费者 ： 一个协程【其他情境中可以多个协程写入一个管道。如何做到多协程同时竞争有序读文件？】读取小文件的
                     内容写进一个管道，若干协程争夺读取管道数据并写入数据库中对应的数据表
   一起跑---比第一：  多个读取数据库表的协程分别读取各自的数据库表，找到了目标内容的立即返回。【本例中 一起跑---比第一 
                      不是特别典型，而应该叫 分头干---合作完 :) ，毕竟是每个协程分别查各自的表，是非查不可的，而非多协
                      程干完全相同的一件事来比谁先完成】
 管道容量控制协程并发数：见jia4kao3xi4tong3的仓库
2、用键值对的值为结构体（文件字段、管道字段）的map来管理起整个工作流程，避免数量众多的管道和文件散乱难以对应起来。但要注意，结构体为map的键值对的值，最好取指针哦！！！
3、数据写入数据库时，分类后的多个小文件--每个小文件的数据入库的多个管道--每个管道都需要拿到独属于自己的数据库连接db,
   避免db不够用各管道协程争夺使用而db挂掉
4、积累一定数量后，拼接insert的sql语句批量插入数据库，这样可以减少对数据库的写操作以避免不必要的降速。
5、用map构建一个简单的缓存
