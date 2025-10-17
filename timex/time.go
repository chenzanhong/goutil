package timex

import (
	"time"
)

// Every 以固定时间间隔执行函数 f。
// 推荐使用闭包来绑定参数。
//
// 示例：
//
//	count := 0
//	go Every(1*time.Second, func() {
//	    count++
//	    fmt.Println("第", count, "个定时任务")
//	})
func Every(duration time.Duration, f func()) {
	if f == nil || duration <= 0 {
		return
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		f()
	}
}
