package modbus

import (
	"bytes"
	"encoding/binary"
)

type RTUData struct {
	SlaveId uint8
	FnCode  uint8
	DataLen uint8
	Crc     uint16
	Data    []uint16
}

// FrameToRTUData wiil deconstruct the raw RTU packet and will convert the data values in it to uint16
func FrameToRTUData(frame []byte) *RTUData {
	d := new(RTUData)
	if len(frame) < 7 { // At least 1 uint16 in the response is required
		return d
	}
	d.SlaveId = frame[0]
	d.FnCode = frame[1]
	d.DataLen = frame[2]
	d.Data = make([]uint16, d.DataLen/2)
	rdr := bytes.NewReader(frame[3 : len(frame)-2])
	binary.Read(rdr, binary.BigEndian, &d.Data)
	d.Crc = binary.LittleEndian.Uint16(frame[len(frame)-2:])
	return d
}
