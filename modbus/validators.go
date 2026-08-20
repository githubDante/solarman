package modbus

import (
	"fmt"
)

var solarmanErrorCodesMap = map[uint8]string{
	5: "IllegalAddress",
	6: "InvalidSerialNumber",
}

type RTUTooShortError struct{}

func (r RTUTooShortError) Error() string {
	return "RTU packet is too short"
}

type SolamranExceptionError struct {
	code uint8
	val  string
}

func (e SolamranExceptionError) Error() string {
	if e.val != "" {
		return fmt.Sprintf("Solarman exception [%d]: %s", e.code, e.val)
	}
	return fmt.Sprintf("Solarman exception [%d]", e.code)
}

type CRCError struct{}

func (e CRCError) Error() string {
	return "CRC error - validation failed"
}

type ModbusError struct {
	code   uint8
	exType string
}

func (e ModbusError) Error() string {
	return fmt.Sprintf("MODBUS error for function [%d] - %s", e.code, e.exType)
}

// ValidRTU checks the packet returned by the datalogger for erorrs
func ValidRTU(rtu []byte) error {
	if len(rtu) < 3 {
		if len(rtu) == 2 {
			if _, ok := solarmanErrorCodesMap[rtu[0]]; ok {
				return SolamranExceptionError{code: rtu[0], val: solarmanErrorCodesMap[rtu[0]]}
			} else {
				return SolamranExceptionError{code: rtu[0], val: "unknown exception code"}
			}
		}
		return RTUTooShortError{}
	}
	if rtu[1]&ModbusExceptionBase != 0 {
		if ex, ok := ExceptionsMap[rtu[2]]; ok {
			return ModbusError{
				code:   rtu[1] - ModbusExceptionBase,
				exType: ex,
			}
		}
		return ModbusError{
			code:   rtu[1] - ModbusExceptionBase,
			exType: fmt.Sprintf("Unknown code [%d]", rtu[2]),
		}
	}
	if !ValidCRC(rtu) {
		return CRCError{}
	}

	return nil
}

func requestValidation(fn ModbusFn, qty uint16, n uint8) bool {
	switch fn {
	case ReadCoilsFn, ReadDicreteInputsFn:
		return qty >= 1 && qty <= 0x07d0
	case ReadHoldingRegistersFn, ReadInputRegistersFn:
		return qty >= 1 && qty <= 0x7d
	case WriteMultipleCoilsFn:
		return qty >= 1 && qty <= 0x07b0
	case WriteSingleCoilFn:
		return qty == 0x0000 || qty == 0xff00 // This is the value to be written
	case WriteMultipleRegistersFn:
		return qty >= 1 && qty <= 0x07b0 && n == uint8(2*qty)
	case WriteSingleRegisterFn:
		return true
	}
	return false
}
