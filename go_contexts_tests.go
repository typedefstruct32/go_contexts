package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 2
	ctx := context.Background()                                 //在外层，我们需要新建一个控制上下文的Context ，然后为这个Context包装上一层又一层的功能（使用了装饰器设计模式）。
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10) //这里用了WithTimeout，让它在指定时间内停止目标函数。
	defer cancelFunc()                                          //在这里加上defer cancelFunc()保证外层函数必定终止目标函数。

	// 3	开个go协程运行目标函数
	go targetFunc(ctx, cancelFunc)

	//手动执行终止函数
	/* 		go func() {
	   		time.Sleep(time.Second * 1)
	   		cancelFunc()
	   	}() */

	// 4
	for {
		select {
		case <-ctx.Done():
			// 5
			switch ctx.Err() {
			case context.DeadlineExceeded: //这里列举了两个结束信号，DeadlineExceeded ：意味着到目标函数的结束时间了（不管完成与否）
				fmt.Println("context timeout exceeded")
				return
			case context.Canceled: //Canceled：手动终止了目标函数。
				fmt.Println("context cancelled by force")
				return
			}
		default:
			time.Sleep(time.Second * 1)
			fmt.Println("sleep 1s")
		}
	}
}

// 1 	这个是需要进行控制的目标函数，第一个参数必定为一个contetx.Context接口（以下简称Context）（这是context的标准使用方式）
// 然后在目标函数里的逻辑自由发挥（所以不需要在这个函数里面考虑context。控制窗口在外层）。
func targetFunc(ctx context.Context, cancelFunc context.CancelFunc) {
	defer cancelFunc() //targetFunc完成任务后，自动手动终止，否则外层会一直等待超时/手动终止信号
	time.Sleep(time.Second * 3)
	fmt.Println("u r here.")
}
