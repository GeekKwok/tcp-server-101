package main

import (
	"fmt"
	"net"

	"github.com/geekkwok/tcp-server-101/frame"
	"github.com/geekkwok/tcp-server-101/packet"
)

func handleConn(c net.Conn) {
	defer func(c net.Conn) {
		_ = c.Close()
	}(c)
	frameCodec := frame.NewMyFrameCodec()

	for {
		// decode the frame to get the payload
		framePayload, err := frameCodec.Decode(c)
		if err != nil {
			fmt.Println("handleConn: frame decode error: ", err)
			return
		}

		// do something with the payload
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: handle packet error: ", err)
			return
		}

		// write ack frame to the connection
		err = frameCodec.Encode(c, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: frame encode error: ", err)
			return
		}
	}
}

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handlePacket: decode packet error: ", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload = %s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handlePacket: encode packet error: ", err)
			return nil, err
		}
		return ackFramePayload, nil
	case *packet.SubmitAck:
		submitAck := p.(*packet.SubmitAck)
		fmt.Printf("recv submit ack: id = %s, result = %d\n", submitAck.ID, submitAck.Result)
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown packet type: %T", p)
	}
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	fmt.Println("server start ok(on *.8080)")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle the new connection
		go handleConn(c)
	}
}
