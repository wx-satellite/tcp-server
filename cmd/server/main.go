package main

import (
	"fmt"
	"net"
	"server/frame"
	"server/packet"
)

func main() {
	srv, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	fmt.Println("server start ok, listen on :8888")

	for {
		client, e := srv.Accept()
		if e != nil {
			fmt.Println("accent client error:", e)
			continue
		}
		go handleConn(client)
	}
}

// handleConn 处理链接
func handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	stream := frame.New()
	for {
		payload, err := stream.Decode(conn)
		if err != nil {
			fmt.Println("handleConn: frame decode error:", err)
			return
		}
		ack, err := handlePayload(payload)
		if err != nil {
			fmt.Println("handlePayload: handle packet error:", err)
			return
		}
		err = stream.Encode(conn, ack)
		if err != nil {
			fmt.Println("handleConn: frame encode error:", err)
			return
		}
	}
}

// handlePayload 处理收到的 PayLoad，并返回响应的 PayLoad
func handlePayload(payload frame.PayLoad) (ack frame.PayLoad, err error) {
	p, err := packet.Decode(payload)
	if err != nil {
		fmt.Println("handleConn: payload decode error:", err)
		return
	}
	switch v := p.(type) {
	case *packet.Submit:
		fmt.Printf("recv submit: id = %s, payload=%s\n", v.ID, string(v.Payload))
		ackPacket := &packet.SubmitAck{
			ID:     v.ID,
			Result: 0, // 返回成功
		}
		ack, err = ackPacket.Encode()
		if err != nil {
			fmt.Println("handleConn: payload encode error:", err)
			return
		}
		return
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}
