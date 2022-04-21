package main

import (
	"fmt"
	"net"
	"server/frame"
	"server/packet"
	"sync"
	"time"
)

func main() {
	var waitGroup sync.WaitGroup
	num := 5
	waitGroup.Add(num)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer waitGroup.Done()
			startClient(i)
		}(i + 1)
	}
	waitGroup.Wait()
}
func startClient(i int) {
	quit := make(chan struct{})
	done := make(chan struct{})
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	fmt.Printf("[client %d]: dial ok", i)

	var (
		counter int
		stream  = frame.New()
	)

	// 接受服务端返回的数据
	go func() {
		for {
			select {
			case <-quit:
				done <- struct{}{}
				return
			default:
			}

			// 使用 SetReadDeadline 方法设置了读超时，这主要是考虑该 Goroutine 可以在收到退出通知时，能及时从 Read 阻塞中跳出来
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 1))
			payload, err := stream.Decode(conn)
			if err != nil {
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						continue
					}
				}
				panic(err)
			}
			body, _ := packet.Decode(payload)
			submitAck, ok := body.(*packet.SubmitAck)
			if !ok {
				panic("not submitAck")
			}
			fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", i, submitAck.ID, submitAck.Result)
		}
	}()

	// 发送数据
	for {
		counter++
		id := fmt.Sprintf("%8d", counter)

		// 内容
		submit := &packet.Submit{
			ID:      id,
			Payload: []byte(fmt.Sprintf("ancd%d", counter)),
		}

		// 内容拼上头信息
		payload, err := packet.Encode(submit)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame length = %d\n",
			i, submit.ID, submit.Payload, len(payload)+4)

		// 发送数据
		err = stream.Encode(conn, payload)
		if err != nil {
			panic(err)
		}

		// 阻塞 1 秒
		time.Sleep(1 * time.Second)

		// 发送10次之后推出
		if counter >= 10 {
			quit <- struct{}{}
			<-done
			fmt.Printf("[client %d]: exit ok\n", i)
			return
		}
	}
}
