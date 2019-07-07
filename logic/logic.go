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

func lockUtil(key string, exec func() error) (needUnlock bool, err error) {
	if err := exec(); err != nil {
		return false, err
	}
	if err := util.LockBlock(key); err != nil {
		// 申请锁失败
		return false, err
	}
	// 需要二次校验
	if err := exec(); err != nil {
		return true, err
	}
}

func AppleAwardV2(userId int64, remark string) error {
	// 1. 去申请锁
	key := fmt.Sprintf("test:{%d}", userId)
	needUnlock, err := lockUtil(key, func() error {
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
		return nil
	})
	if needUnlock {
		util.UnLodk(key)
	}
	if err != nil {
		return err
	}
	// 3. 插入
	if _, err := dao.Award.InsertOne(userId, remark); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
