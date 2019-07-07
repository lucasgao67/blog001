package logic

import (
	"errors"
	"fmt"
	"github.com/LucasGao67/blog001/dao"
	"github.com/LucasGao67/blog001/util"
	"sync"
)

var lock sync.Locker

func AppleAward(userId int64, remark string) error {

	// 1. 去申请锁
	key := fmt.Sprintf("test:{%d}", userId)
	//lock.()
	util.LockBlock(key)
	defer util.UnLodk(key)

	// 2. 去查询

	info, err := dao.Award.FindOne(userId)
	if err != nil {
		msg := "查询奖品申请信息失败"
		fmt.Printf("%s :%s\n", msg, err.Error())
		return errors.New("查询奖品申请信息失败")
	}

	if info != nil {
		msg := "已经申请不能重复申请"
		fmt.Println(msg)
		return errors.New(msg)
	}
	// 3. 插入
	if _, err := dao.Award.InsertOne(userId, remark); err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 结束
	return nil

}
