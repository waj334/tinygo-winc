package utilities

import (
	"bytes"
	"testing"
)

func Test_pad(t *testing.T) {
	type args struct {
		length    int
		alignment int
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
	}{
		{
			name: "zero-length",
			args: args{
				length:    0,
				alignment: 4,
			},
			expected: []byte{0, 0, 0, 0},
		},
		{
			name: "full-length",
			args: args{
				length:    4,
				alignment: 4,
			},
			expected: []byte{}, // No padding is expected
		},
		{
			name: "half-length",
			args: args{
				length:    2,
				alignment: 4,
			},
			expected: []byte{0, 0}, // No padding is expected
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(make([]byte, 0, 4))
			Pad(tt.args.length, tt.args.alignment, buf)

			if !bytes.Equal(buf.Bytes(), tt.expected) {
				t.Errorf("\nExpected \t% #x\nGot \t\t % #x", tt.expected, buf.Bytes())
			}
		})
	}
}
