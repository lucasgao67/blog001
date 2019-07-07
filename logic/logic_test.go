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
