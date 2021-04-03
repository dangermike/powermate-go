package message

import (
	"encoding/binary"
	"errors"
	"strconv"
)

type Message struct {
	IsDown   bool
	Rotation int8
}

var invalid = Message{false, -128}

const (
	sig = uint32(0x004F1000)
)

func Parse(buf []byte) (Message, error) {
	marker := binary.BigEndian.Uint32(buf[2:6])
	if marker != sig {
		return invalid, errors.New("Invalid marker")
	}
	return Message{
		IsDown:   buf[0] == 1,
		Rotation: int8(buf[1]),
	}, nil
}

func (m Message) String() string {
	if m.Rotation == 0 {
		if m.IsDown {
			return "press"
		}
		return "release"
	}

	retval := ""

	if m.IsDown {
		retval += "down_"
	}
	retval += "rotate("

	return retval + strconv.Itoa(int(m.Rotation)) + ")"
}
