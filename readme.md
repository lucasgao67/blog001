# 问题背景
最近业务上遇到这样的场景，觉得很有代表性，所以拿来说一说。我们有一个奖品发放系统，当用户申请奖品的时候，首先需要判断用户有没有申请过奖品，如果没有申请过，则去奖品总量扣除一个，然后再把用户申请记录写回数据库。

流程如下：

| 时序 | 事件               |
| ---- | ------------------ |
| t1   | 检测用户申请过奖品 |
| t2   | 奖品扣除           |
| t3   | 插入用户申请记录   |
| t4   | 返回申请成功       |

如果是正常的但线程执行完全没有问题，但是我们是并发的。所以很有可能会有好多事件都到t2导致多次申请，或者奖品多次扣除等情况。

# 解决方案
## 方案1：事务
最开始想到的就是数据库的事务了，在t1开始前启动事务，t1查询使用for update加锁,其他事务继续执行相同条件的for update的时候则会block，直到上个事务执行完成。这种方案很完美，但是只限于InnoDB，在tidb上行不通。tidb的事务是乐观锁，所以在t1的时候不会阻塞，那么可以在t3提交的时候conflict再回滚也可以啊。抱歉，tidb并不会产生插入冲突，因为tidb的锁不支持gap lock和next-key lock,所以如果我们的奖品申请表没有唯一索引冲突的话，完全可以插入（因为我某些设计，所以这里我们不能使用唯一索引）。

这个方案是行不通了，那么换种思路，既然tidb不支持for update block,那么我们是不是可以使用分布式锁来解决。

## 方案2：分布式锁
在检测用户有没有申请过奖品之前，我们可以以用户id为key申请分布式锁（可以使用redis实现），申请成功则进行下一步，其他用户阻塞等待锁释放。
具体实现如下：

```go
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
```

测试
```go
package logic

import (
	"fmt"
	"github.com/LucasGao67/blog001/util"
	"sync"
	"testing"
	"time"
)

func init() {
	util.SqlInit()
	util.RedisInit()

}

func TestAppleAward(t *testing.T) {
	// 插入1000条数据测试
	st := time.Now()
	numSucc := 0
	numFail := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			err := AppleAward(1000+int64(i), "test 001")
			if err != nil {
				numFail++
			} else {
				numSucc++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	et := time.Now()
	fmt.Printf("耗时： %ds\n", et.Unix()-st.Unix())
	fmt.Printf("成功： %d\n", numSucc)
	fmt.Printf("失败： %d\n", numFail)
}
```

这样如果userid都一样是没有问题的，但是一旦userid重复申请，都会阻塞到检测userid状态上，也就算上述的t1。所以这边加锁需要优化下

## 方案三 分布式锁优化
采用最经典的二次校验，先查询，满足条件，加锁，再校验
```go

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

```

# github地址
[https://github.com/LucasGao67/blog001](https://github.com/LucasGao67/blog001)

