package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	stream := New()
	if stream == nil {
		t.Errorf("want non-nill, got nil")
	}

	fmt.Println(bytes.Join([][]byte{[]byte{'a'},[]byte{'b'}},[]byte{'z'}))
}


func TestFrame_Encode(t *testing.T) {
	stream := New()
	buf := make([]byte,128)
	rw := bytes.NewBuffer(buf)

	err := stream.Encode(rw, []byte("hello"))
	if err != nil {
		t.Errorf("want nil, got %s", err.Error())
	}

	var total int32
	err = binary.Read(rw, binary.BigEndian, &total)
	if err != nil {
		t.Errorf("want nil, got %s",err.Error())
	}

	if total != 9 {
		t.Errorf("want 9, got %d",total)
	}

	// binary 从 rw 读取了4个字节，得到 total，rw.Bytes() 返回剩下的字节，也就是 hello
	if string(rw.Bytes()) != "hello" {
		t.Errorf("want hello, got %s",string(rw.Bytes()))
	}
}


func TestFrame_Decode(t *testing.T) {
	stream := New()
	payload,err := stream.Decode(bytes.NewReader([]byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}))
	if err != nil {
		t.Errorf("want nil, got %s",err.Error())
	}

	if string(payload) != "hello" {
		t.Errorf("want hello, got %s",string(payload))
	}
}