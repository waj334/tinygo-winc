/*
MIT License

Copyright (c) 2022 waj334

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package protocol

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

const (
	_NBIT31 uint32 = 0x80000000
	_NBIT30 uint32 = 0x40000000
	_NBIT29 uint32 = 0x20000000
	_NBIT28 uint32 = 0x10000000
	_NBIT27 uint32 = 0x08000000
	_NBIT26 uint32 = 0x04000000
	_NBIT25 uint32 = 0x02000000
	_NBIT24 uint32 = 0x01000000
	_NBIT23 uint32 = 0x00800000
	_NBIT22 uint32 = 0x00400000
	_NBIT21 uint32 = 0x00200000
	_NBIT20 uint32 = 0x00100000
	_NBIT19 uint32 = 0x00080000
	_NBIT18 uint32 = 0x00040000
	_NBIT17 uint32 = 0x00020000
	_NBIT16 uint32 = 0x00010000
	_NBIT15 uint32 = 0x00008000
	_NBIT14 uint32 = 0x00004000
	_NBIT13 uint32 = 0x00002000
	_NBIT12 uint32 = 0x00001000
	_NBIT11 uint32 = 0x00000800
	_NBIT10 uint32 = 0x00000400
	_NBIT9  uint32 = 0x00000200
	_NBIT8  uint32 = 0x00000100
	_NBIT7  uint32 = 0x00000080
	_NBIT6  uint32 = 0x00000040
	_NBIT5  uint32 = 0x00000020
	_NBIT4  uint32 = 0x00000010
	_NBIT3  uint32 = 0x00000008
	_NBIT2  uint32 = 0x00000004
	_NBIT1  uint32 = 0x00000002
	_NBIT0  uint32 = 0x00000001

	_HAVE_SDIO_IRQ_GPIO_BIT     = _NBIT0
	_HAVE_USE_PMU_BIT           = _NBIT1
	_HAVE_SLEEP_CLK_SRC_RTC_BIT = _NBIT2
	_HAVE_SLEEP_CLK_SRC_XO_BIT  = _NBIT3
	_HAVE_EXT_PA_INV_TX_RX      = _NBIT4
	_HAVE_LEGACY_RF_SETTINGS    = _NBIT5
	_HAVE_LOGS_DISABLED_BIT     = _NBIT6
	_HAVE_ETHERNET_MODE_BIT     = _NBIT7
	_HAVE_RESERVED1_BIT         = _NBIT8
	_HAVE_RESERVED2_BIT         = _NBIT9
	_HAVE_XO_XTALGM2_DIS_BIT    = _NBIT10

	_rNMI_GP_REG_0       uint32 = 0x149c
	_rNMI_GP_REG_1       uint32 = 0x14A0
	_rNMI_GP_REG_2       uint32 = 0xc0008
	_rNMI_GLB_RESET      uint32 = 0x1400
	_rNMI_BOOT_RESET_MUX uint32 = 0x1118

	_NMI_STATE_REG     uint32 = 0x108c
	_BOOTROM_REG       uint32 = 0xc000c
	_NMI_REV_REG       uint32 = 0x207ac
	_NMI_REV_REG_ATE   uint32 = 0x1048
	_WAIT_FOR_HOST_REG uint32 = 0x207bc
	_FINISH_INIT_STATE uint64 = 0x02532636
	_FINISH_BOOT_ROM   uint64 = 0x10add09e
	_START_FIRMWARE    uint64 = 0xef522f61
	_START_PS_FIRMWARE uint64 = 0x94992610

	_NMI_PERIPH_REG_BASE uint32 = 0x1000
	_NMI_INTR_REG_BASE   uint32 = _NMI_PERIPH_REG_BASE + 0xa00
	_NMI_CHIPID          uint32 = _NMI_PERIPH_REG_BASE
	_NMI_PIN_MUX_0       uint32 = _NMI_PERIPH_REG_BASE + 0x408
	_NMI_INTR_ENABLE     uint32 = _NMI_INTR_REG_BASE

	_NMI_GP_REG_0 uint32 = 0x149c
	_NMI_GP_REG_1 uint32 = 0x14A0
	_NMI_GP_REG_2 uint32 = 0xc0008

	_NMI_SPI_REG_BASE         uint32 = 0xe800
	_NMI_SPI_CTL              uint32 = _NMI_SPI_REG_BASE
	_NMI_SPI_MASTER_DMA_ADDR  uint32 = _NMI_SPI_REG_BASE + 0x4
	_NMI_SPI_MASTER_DMA_COUNT uint32 = _NMI_SPI_REG_BASE + 0x8
	_NMI_SPI_SLAVE_DMA_ADDR   uint32 = _NMI_SPI_REG_BASE + 0xc
	_NMI_SPI_SLAVE_DMA_COUNT  uint32 = _NMI_SPI_REG_BASE + 0x10
	_NMI_SPI_TX_MODE          uint32 = _NMI_SPI_REG_BASE + 0x20
	_NMI_SPI_PROTOCOL_CONFIG  uint32 = _NMI_SPI_REG_BASE + 0x24
	_NMI_SPI_INTR_CTL         uint32 = _NMI_SPI_REG_BASE + 0x2c
	_NMI_SPI_MISC_CTRL        uint32 = _NMI_SPI_REG_BASE + 0x48

	_NMI_SPI_PROTOCOL_OFFSET uint32 = _NMI_SPI_PROTOCOL_CONFIG - _NMI_SPI_REG_BASE

	_SPI_BASE uint32 = _NMI_SPI_REG_BASE

	_CORT_HOST_COMM uint32 = 0x10
	_HOST_CORT_COMM uint32 = 0x0B
	_WAKE_CLOCK_REG uint32 = 0x1
	_CLOCKS_EN_REG  uint32 = 0xF

	_WIFI_NORMAL_MODE            = 1
	_WIFI_HOST_RCV_CTRL_0 uint32 = 0x1070
	_WIFI_HOST_RCV_CTRL_1 uint32 = 0x1084
	_WIFI_HOST_RCV_CTRL_2 uint32 = 0x1078
	_WIFI_HOST_RCV_CTRL_3 uint32 = 0x106c
	_WIFI_HOST_RCV_CTRL_4 uint32 = 0x150400
	_WIFI_HOST_RCV_CTRL_5 uint32 = 0x1088

	_FIRMWARE_VERSION uint32 = 0x1377 // 19.7.7
	_DRIVER_VERSION   uint32 = 0x1330 // 19.3.0
	_VERSION                 = _FIRMWARE_VERSION | (_DRIVER_VERSION << 16)
	_REV_2B0          uint32 = 0x2B0
	_REV_B0           uint32 = 0x2B0
	_REV_3A0          uint32 = 0x3A0

	// The maximum transmission unit for sending data blocks over the SPI bus
	_SPI_BUS_MTU = 2048 - 8
)

type sequence byte

const (
	first    sequence = 0x01
	sendRecv sequence = 0x02
	last     sequence = 0x03
	reserved sequence = 0xFF
	mask     sequence = 0x0F
)

func getSequence(header byte) sequence {
	return sequence(header) & mask
}

type numeric interface {
	constraints.Float | constraints.Signed | constraints.Unsigned
}

func min[T numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func PrintBytes(s []byte) {
	printBytes(s)
}

func printBytes(s []byte) {
	for _, b := range s {
		PrintByte(b)
		print(" ")
	}

	print("\n\r")
}

func printBytesL(s []byte) {
	for _, b := range s {
		PrintByte(b)
		print(" ")
	}
}

func PrintByte(b byte) {
	upper := (b >> 4) & 0x0F
	lower := b & 0x0F

	print(hex(upper))
	print(hex(lower))
}

func hex(b byte) string {
	switch b {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "A"
	case 11:
		return "B"
	case 12:
		return "C"
	case 13:
		return "D"
	case 14:
		return "E"
	case 15:
		return "F"
	default:
		return "NOP"
	}
}

func ToUint32(buffer []byte) (val uint32) {
	val |= uint32(buffer[3]) << 24
	val |= uint32(buffer[2]) << 16
	val |= uint32(buffer[1]) << 8
	val |= uint32(buffer[0])
	return
}

// Itoa converts val to a decimal string.
func Itoa(val int) string {
	if val < 0 {
		return "-" + Uitoa(uint(-val))
	}
	return Uitoa(uint(val))
}

// Uitoa converts val to a decimal string.
func Uitoa(val uint) string {
	if val == 0 { // avoid string allocation
		return "0"
	}
	var buf [20]byte // big enough for 64bit value base 10
	i := len(buf) - 1
	for val >= 10 {
		q := val / 10
		buf[i] = byte('0' + val - q*10)
		i--
		val = q
	}
	// val < 10
	buf[i] = byte('0' + val)
	return string(buf[i:])
}

func Bytes[T any](in *T) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(in)), unsafe.Sizeof(*in))
}
