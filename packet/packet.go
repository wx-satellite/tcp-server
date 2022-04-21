package packet

import (
	"bytes"
	"fmt"
)

/*
packet 定义： header + body

header：
	1 字节：类型


body：
	8字节的：ID
	任意字节：payload
ack body：
	8字节：ID
	1字节：result
 */


const (
	CommandConn = iota + 0x01
	CommandSubmit
)

const (
	CommandConnAck = iota + 0x81
	CommandSubmitAck
)

type Packet interface {
	Decode([]byte) error
	Encode()([]byte,error)
}


type Submit struct {
	ID string
	Payload []byte
}

func (s*Submit)Decode(body []byte)(err error) {
	s.ID = string(body[:8])
	s.Payload = body[9:]
	return
}

func (s *Submit) Encode()(body []byte, err error) {
	body = bytes.Join([][]byte{[]byte(s.ID)[:8],s.Payload},nil)
	return
}


type SubmitAck struct {
	ID string
	Result uint8 // byte 就是 uint8 的别名
}

func (s *SubmitAck) Decode(body []byte)(err error) {
	s.ID = string(body[:8])
	s.Result = body[8]
	return
}

func (s *SubmitAck) Encode()(body []byte, err error) {
	body = bytes.Join([][]byte{[]byte(s.ID),{s.Result}},nil)
	return
}


func Decode(packet []byte)(p Packet, err error) {
	header := packet[0]
	body := packet[1:]
	switch header {
	case CommandConn:
		return
	case CommandConnAck:
		return
	case CommandSubmit:
		submit := &Submit{}
		err = submit.Decode(body)
		if err != nil {
			return
		}
		p = submit
		return
	case CommandSubmitAck:
		ack := &SubmitAck{}
		err = ack.Decode(body)
		if err != nil {
			return
		}
		p = ack
		return
	default:
		err = fmt.Errorf("unknown header")
		return
	}
}

func Encode(p Packet)(packet []byte, err error) {
	var (
		header uint8
		body []byte
	)
	switch v := p.(type) {
	case *SubmitAck:
		header = CommandSubmit
		body, err= v.Encode()
		if err != nil {
			return
		}
	case *Submit:
		header = CommandSubmitAck
		body, err = v.Encode()
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unkown type")
		return
	}
	packet = bytes.Join([][]byte{{header},body},nil)
	return
}