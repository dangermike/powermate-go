package main

import (
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/dangermike/powermate-go/message"
	"github.com/karalabe/hid"
	"github.com/micmonay/keybd_event"
	"github.com/sqp/pulseaudio"
)

const (
	// see http://www.linux-usb.org/usb.ids
	vendorID  = uint16(0x077D)
	productID = uint16(0x0410)
	maxInt8   = int8(127)
)

func main() {
	client, err := pulseaudio.New()
	if err != nil {
		panic(err)
	}
	defer client.Close()
}

func mainxx() {
	devs := hid.Enumerate(vendorID, productID)
	if len(devs) == 0 {
		fmt.Printf("Failed to find device")
	}
	dev, err := devs[0].Open()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, need to wait 2 seconds
	if runtime.GOOS == "linux" {
		fmt.Println("initializing keyboard")
		time.Sleep(2 * time.Second)
	}

	buf := make([]byte, 8)
	for {
		cnt, err := dev.Read(buf)
		if err == io.EOF {
			fmt.Println("device closed")
			return
		}
		if err != nil {
			panic(err)
		}

		msg, err := message.Parse(buf[:cnt])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(msg)
			if msg.Rotation == 0 && msg.IsDown {
			} else {
				// evt := keybd_event.VK_VOLUMEUP
				// if msg.Rotation < 0 {
				// 	evt = keybd_event.VK_VOLUMEDOWN
				// }
				evt := 0x48
				magnitude := int(abs8(msg.Rotation))
				for i := 0; i < magnitude; i++ {
					kb.AddKey(evt)
				}
				err = kb.Launching()
				if err != nil {
					panic(err)
				}
				kb.Clear()
			}
		}
	}
}

func abs8(val int8) int8 {
	if val >= 0 {
		return val
	}
	return (maxInt8 - (val & maxInt8)) + 1
}
