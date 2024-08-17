package updreliablity

import (
	"bytes"
	"encoding/binary"
	"errors"
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


type Testc struct {
	name     string
	input    []byte
	expected ReadRQ
	err      error
}


func TestUnmarshalBinary(t *testing.T) {
	tcs := []Testc{
		{
			name: "valid input",
			input: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpRRQ)
				buf.WriteString("test.txt")
				buf.WriteByte(0)
				buf.WriteString("octet")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
			expected: ReadRQ{
				Filename: "test.txt",
				Mode:     "octet",
			},
			err: nil,
		},
		{
			name: "invalid opcode",
			input: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpAck)
				buf.WriteString("test.txt")
				buf.WriteByte(0)
				buf.WriteString("octet")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
			expected: ReadRQ{},
			err:      errors.New("invalid read request"),
		},
		{
			name: "invalid filename",
			input: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpRRQ)
				buf.WriteByte(0)
				buf.WriteString("octet")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
			expected: ReadRQ{},
			err:      errors.New("invalid read request"),
		},
		{
			name: "invalid mode",
			input: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, OpRRQ)
				buf.WriteString("test.txt")
				buf.WriteByte(0)
				buf.WriteString("netascii")
				buf.WriteByte(0)
				return buf.Bytes()
			}(),
			expected: ReadRQ{},
			err:      errors.New("invalid read request: netascii not supported"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var r ReadRQ
			err := r.UnmarshalBinary(tc.input)

			if err != nil && tc.err == nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if err == nil && tc.err != nil {
				t.Fatalf("Expected error: %v, got nil", tc.err)
			}

			if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
				t.Fatalf("Expected error: %v, got: %v", tc.err, err)
			}

			if r != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, r)
			}
		})
	}
}