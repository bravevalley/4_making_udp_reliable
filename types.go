package updreliablity

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// TFTP -- RFC 1350

// Packet size
const (
	DatagramSize     = 516              // Safest datagram size tranferreable over network without fragmentation
	MaxTFTPBlockSize = DatagramSize - 4 // Datagram size - 4bytes datagram header
)

// TFTP Headers code
type Opcode uint16

const (
	OpRRQ  Opcode = iota + 1 // Read request Code
	_                 // Write request not supported
	OpData            // Data (DATA) code 
	OpAck             // Acknowledgment (ACK) code
	OpErr             // Error (ERROR) code
)

// TFTP Error codes
type ErrCode uint16

const (
	ErrNotDefined      = iota // 0 - Not defined, see error message (if any)
	ErrFileNotFound           // 1 - File not found
	ErrAccessViolation        // 2 - Access violation
	ErrDiskFull               // 3 - Disk full or allocation exceeded
	ErrIllegalOp              // 4 -  Illegal TFTP operation
	ErrUnknownID              // 5 - Unknown transfer ID
	ErrFileExists             // 6 - File already exist
	ErrNoUser                 // 7 - No such user
)




// 2 bytes     string    1 byte     string   1 byte
// ------------------------------------------------
// | Opcode |  Filename  |   0  |    Mode    |   0  |
// ------------------------------------------------

// Read request

type ReadRQ struct {
	Filename string
	Mode string
}

// MarshalBinary writes the protocol header which contains the filename 
// and mode to a buffer. It returns a slice of bytes and a non-nil error if
// conversion is successful 
func (r ReadRQ) MarshalBinary() ([]byte, error) {
	mode := "octet"

	if r.Mode != "" {
		mode = r.Mode
	}

	buf := new(bytes.Buffer)
	capacity := 2+len(r.Filename)+1+len(r.Mode)+1
	buf.Grow(capacity)

	err := binary.Write(buf, binary.LittleEndian, OpRRQ)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(r.Filename)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}	

	_, err = buf.WriteString(mode)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *ReadRQ) UnmarshalBinary(p []byte) error {
	buf := bytes.NewBuffer(p)

	var opcode Opcode

	// Read Operation code off
	err := binary.Read(buf, binary.LittleEndian, &opcode)
	if err != nil {
		return errors.New("invalid read request")
	}

	if opcode != OpRRQ {
		return errors.New("invalid read request")
	}

	// Read filename with delimiter
	plxhlder, err := buf.ReadString(0)
	if err != nil {
		return errors.New("invalid read request")
	}

	r.Filename = strings.TrimRight(plxhlder, "\x00")
	if len(r.Filename) <= 0 {
		return errors.New("invalid read request")
	}


	// Read mode
	plxhlder, err = buf.ReadString(0)
	if err != nil {
		return errors.New("invalid read request")
	}

	mode := strings.TrimRight(plxhlder, "\x00")
	mode = strings.ToLower(mode)
	if mode != "octet" {
		return fmt.Errorf("invalid read request: %s not supported", mode)
	} else {
		r.Mode = mode
	}

	return nil 
}