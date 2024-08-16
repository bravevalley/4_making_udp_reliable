package updreliablity

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type testCase struct {
	name     string
	input    ReadRQ
	expected []byte
}

func TestMarshalBinary(t *testing.T) {
	testCases := []testCase{
		{
			name: "default mode",
			input: ReadRQ{
				Filename: "test.txt",
			},
			expected: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpRRQ)
				buf.WriteString("test.txt")
				buf.WriteByte(0)
				buf.WriteString("octet")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
		},
		{
			name: "custom mode",
			input: ReadRQ{
				Filename: "test.txt",
				Mode:     "netascii",
			},
			expected: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpRRQ)
				buf.WriteString("test.txt")
				buf.WriteByte(0)
				buf.WriteString("netascii")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.input.MarshalBinary()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}