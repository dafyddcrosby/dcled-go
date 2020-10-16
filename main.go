// Copyright David Crosby, 2020
// MIT License

package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/google/gousb"
)

const (
	VENDORID  = 0x1d34
	PRODUCTID = 0x0013
)

// Technically, only 7 rows of LEDs, but the USB packet needs an imaginary
// 8th row, which can be blank bytes
type Board struct {
	device             *gousb.Device
	default_brightness byte
	leds               [8][]byte
}

func (b Board) write_screen() {
	for l := 0; l < 8; l += 2 {
		var d []byte
		d = append(d, uint8(0), uint8(l)) // brightness, row
		d = append(d, b.leds[l]...)       // line 1
		d = append(d, b.leds[l+1]...)     // line 2
		b.write_packet(d)
	}

}

func (b Board) write_packet(data []byte) {
	b.device.Control(0x21, 0x09, 0x000, 0x0000, data)
}

func (b Board) test_random() {
	for {
		for i := 0; i < 7; i += 1 {
			b.leds[i] = []byte{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))}
		}
		b.write_screen()
	}
}
func (b Board) test_diamond() {
	b.leds[0] = []byte{0xff, 0xfe, 0xff}
	b.leds[1] = []byte{0xff, 0xfd, 0x7f}
	b.leds[2] = []byte{0xff, 0xfb, 0xbf}
	b.leds[3] = []byte{0xff, 0xf7, 0xdf}
	b.leds[4] = []byte{0xff, 0xfb, 0xbf}
	b.leds[5] = []byte{0xff, 0xfd, 0x7f}
	b.leds[6] = []byte{0xff, 0xfe, 0xff}
	b.write_screen()
}

func is_dc_board(vendor gousb.ID, product gousb.ID) bool {
	return vendor == VENDORID && product == PRODUCTID
}

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devs, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if is_dc_board(desc.Vendor, desc.Product) {
			log.Println("Found a DC board")
			return true
		}

		return false
	})

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	if len(devs) == 0 {
		log.Fatal("No DC board found")
	}
	dev := devs[0]

	board := Board{device: dev}
	board.test_diamond()
	board.test_random()
}
