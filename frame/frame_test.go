package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestNewMyFrameCodec(t *testing.T) {
	codec := NewMyFrameCodec()
	if codec == nil {
		t.Errorf("NewMyFrameCodec() should not return nil")
	}
}

func TestEncode(t *testing.T) {
	codec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)

	err := codec.Encode(rw, []byte("hello world"))
	if err != nil {
		t.Errorf("Encode() failed: %v", err)
	}

	// verify the encoded frame
	var totalLen int32
	err = binary.Read(rw, binary.BigEndian, &totalLen)
	if err != nil {
		t.Errorf("binary.Read() failed: %v", err)
	}

	if totalLen != 15 {
		t.Errorf("totalLen should be 15, but got %d", totalLen)
	}

	left := rw.Bytes()
	if string(left) != "hello world" {
		t.Errorf("encoded frame should be 'hello world', but got %s", left)
	}
}

func TestDecode(t *testing.T) {
	codec := NewMyFrameCodec()
	data := []byte{0, 0, 0, 15, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}

	payload, err := codec.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("Decode() failed: %v", err)
	}

	if string(payload) != "hello world" {
		t.Errorf("Decoded payload should be 'hello world', but got %s", payload)
	}
}

type ReturnErrorWriter struct {
	W  io.Writer
	Wn int // 第几次调用 Write 返回错误
	wc int // 写操作次数计数
}

func (w *ReturnErrorWriter) Write(p []byte) (n int, err error) {
	w.wc++
	if w.wc >= w.Wn {
		return 0, errors.New("write error")
	}
	return w.W.Write(p)
}

type ReturnErrorReader struct {
	R  io.Reader
	Rn int // 第几次调用 Read 返回错误
	rc int // 读操作次数计数
}

func (r *ReturnErrorReader) Read(p []byte) (n int, err error) {
	r.rc++
	if r.rc >= r.Rn {
		return 0, errors.New("read error")
	}
	return r.R.Read(p)
}

func TestEncodeWithWriteFail(t *testing.T) {
	codec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	w := bytes.NewBuffer(buf)

	// 模拟 binary.Write 写入返回错误
	err := codec.Encode(&ReturnErrorWriter{W: w, Wn: 1}, []byte("hello world"))
	if err == nil {
		t.Errorf("Encode() should return error")
	}

	// 模拟 w.Write 写入返回错误
	err = codec.Encode(&ReturnErrorWriter{W: w, Wn: 2}, []byte("hello world"))
	if err == nil {
		t.Errorf("Encode() should return error")
	}
}

func TestDecodeWithReadFail(t *testing.T) {
	codec := NewMyFrameCodec()
	data := []byte{0, 0, 0, 15, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}

	// 模拟 binary.Read 读取返回错误
	_, err := codec.Decode(&ReturnErrorReader{R: bytes.NewReader(data), Rn: 1})
	if err == nil {
		t.Errorf("Decode() should return error")
	}
}
