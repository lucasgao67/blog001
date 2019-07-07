package main

import (
	"github.com/LucasGao67/blog001/dao"
	"github.com/LucasGao67/blog001/util"
)

func main() {
	util.SqlInit()
	util.RedisInit()
	//fmt.Printf("12")
	dao.Award.InsertOne(12, "a test")
}
