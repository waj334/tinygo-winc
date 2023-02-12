package utilities

import "bytes"

func Pad(length, alignment int, buf *bytes.Buffer) {
	count := alignment
	if length > 0 {
		count = (alignment - (length % alignment)) % alignment
	}
	for i := 0; i < count; i++ {
		buf.WriteByte(0)
	}
}
