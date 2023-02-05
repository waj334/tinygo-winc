package winc

import "bytes"

func pad(length, alignment int, buf *bytes.Buffer) {
	count := (alignment - (length % alignment)) % alignment
	for i := 0; i < count; i++ {
		buf.WriteByte(0)
	}
}
