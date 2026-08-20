package modbus

import (
	"encoding/binary"
	"errors"
)

// CreateReadHolding will construct ADU for FunctionCode 0x03
func CreateReadHolding(slaveid uint8, address, quantity uint16) (adu []byte, err error) {
	if ok := requestValidation(ReadHoldingRegistersFn, quantity, 0); !ok {
		return adu, errors.New("MODBUS vaildation failed")
	}
	adu = []byte{slaveid, ReadHoldingRegistersFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, address)
	adu, _ = binary.Append(adu, binary.BigEndian, quantity)

	return CalcCRC(adu), nil
}

// CreateReadInputs will construct ADU for FunctionCode 0x04
func CreateReadInputs(slaveId uint8, address, quantity uint16) (adu []byte, err error) {
	if ok := requestValidation(ReadInputRegistersFn, quantity, 0); !ok {
		return adu, errors.New("MODBUS vaildation failed")
	}
	adu = []byte{slaveId, ReadInputRegistersFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, address)
	adu, _ = binary.Append(adu, binary.BigEndian, quantity)

	return CalcCRC(adu), nil
}

// CreateWriteMultipleCoils - Function 0x0f ADU
func CreateWriteMultipleCoils(slaveId uint8, startAddress, quantity uint16, data []byte) (adu []byte, err error) {
	n := uint8(quantity / 8)
	if uint8(quantity%8) != 0 {
		n++
	}
	if ok := requestValidation(WriteMultipleCoilsFn, quantity, 0); !ok {
		return adu, errors.New("MODBUS vaildation failed")
	}
	adu = []byte{slaveId, WriteMultipleCoilsFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, startAddress)
	adu, _ = binary.Append(adu, binary.BigEndian, quantity)
	adu, _ = binary.Append(adu, binary.BigEndian, n)
	adu, _ = binary.Append(adu, binary.BigEndian, data)
	return CalcCRC(adu), nil
}

// CreateWriteSingleCoil - Function 0x05 ADU
func CreateWriteSingleCoil(slaveId uint8, startAddress, value uint16) (adu []byte, err error) {
	if ok := requestValidation(WriteSingleCoilFn, value, 0); !ok {
		return adu, errors.New("MODBUS vaildation failed")
	}
	adu = []byte{slaveId, WriteSingleCoilFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, startAddress)
	adu, _ = binary.Append(adu, binary.BigEndian, value)

	return CalcCRC(adu), nil
}

// CreateWriteMultipleRegisters - Function 0x10 ADU
func CreateWriteMultipleRegisters(slaveId uint8, startAddress, quantity uint16, data []byte) (adu []byte, err error) {
	n := uint8(quantity * 2)
	if ok := requestValidation(WriteMultipleRegistersFn, quantity, n); !ok {
		return adu, errors.New("MODBUS vaildation failed")
	}
	adu = []byte{slaveId, WriteMultipleRegistersFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, startAddress)
	adu, _ = binary.Append(adu, binary.BigEndian, quantity)
	adu, _ = binary.Append(adu, binary.BigEndian, n)
	adu, _ = binary.Append(adu, binary.BigEndian, data)
	return CalcCRC(adu), nil
}

// CreateWriteSingleRegister - Function 0x06 ADU
func CreateWriteSingleRegister(slaveId uint8, startAddress, value uint16) (adu []byte, err error) {
	// Does not need validation - it's restricted by the uint16 (addr & value)
	adu = []byte{slaveId, WriteSingleRegisterFn.Uint8()}
	adu, _ = binary.Append(adu, binary.BigEndian, startAddress)
	adu, _ = binary.Append(adu, binary.BigEndian, value)

	return CalcCRC(adu), nil
}
