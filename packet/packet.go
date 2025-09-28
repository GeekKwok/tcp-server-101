package packet

import (
	"bytes"
	"fmt"
)

// Packet 协议定义
/**
+--------+--------+--------+--------+--------+--------+--------+--------+
packet header:
---
1 byte : commandID
+--------+--------+--------+--------+--------+--------+--------+--------+
submit packet body:
---
8 byte : ID string
任意字节 : payload
+--------+--------+--------+--------+--------+--------+--------+--------+
reply packet body:
---
8 byte : ID string
1 byte : result
+--------+--------+--------+--------+--------+--------+--------+--------+
*/

const (
	CommandConn   = iota + 0x01 // 0x01
	CommandSubmit               // 0x02
)

const (
	CommandConnAck   = iota + 0x81 // 0x81
	CommandSubmitAck               // 0x82
)

type Packet interface {
	Decode([]byte) error     // []byte -> Packet
	Encode() ([]byte, error) // Packet -> []byte
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(packetBody []byte) error {
	s.ID = string(packetBody[:8])
	s.Payload = packetBody[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(packetBody []byte) error {
	s.ID = string(packetBody[:8])
	s.Result = packetBody[8]
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	packetBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		s := Submit{}
		err := s.Decode(packetBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(packetBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(packet Packet) ([]byte, error) {
	var commandID uint8
	var packetBody []byte
	var err error

	switch packet.(type) {
	case *Submit:
		commandID = CommandSubmit
		packetBody, err = packet.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		packetBody, err = packet.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown packet type [%T]", packet)
	}

	return bytes.Join([][]byte{[]byte{commandID}, packetBody}, nil), nil
}
