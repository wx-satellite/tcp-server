package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

/*
frame 定义：frameHeader + framePayload

frameHeader 定义：
	4 bytes：值表示整个frame的长度，其中头的长度是固定的--4字节

framePayload 定义：
	Packet
 */

type PayLoad []byte

type Stream interface {
	Encode(writer io.Writer, load PayLoad)error // 打包
	Decode(reader io.Reader)(PayLoad, error) // 解包
}


var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")


type frame struct {

}


func New() Stream {
	return &frame{}
}


func (f *frame) Encode(w io.Writer, payload PayLoad) (err error) {
	var totalLength  = 4 + int32(len(payload))
	err = binary.Write(w, binary.BigEndian, &totalLength)
	if err != nil {
		return
	}
	n, err :=w.Write(payload)
	if err != nil {
		return
	}
	if n != len(payload) {
		err = ErrShortWrite
		return
	}
	return
}


func (f *frame) Decode(r io.Reader)(payLoad PayLoad, err error) {
	var totalLength int32
	err = binary.Read(r, binary.BigEndian, &totalLength)
	if err != nil {
		return
	}
	buf := make([]byte, totalLength-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return
	}
	if n != int(totalLength-4) {
		err = ErrShortRead
		return
	}
	payLoad = buf
	return
}