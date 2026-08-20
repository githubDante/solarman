package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/go-jose/go-jose/v4/testutils/assert"
)

func TestModbus(t *testing.T) {
	// 0103003c00014406
	pdu, _ := hex.DecodeString("0103003c0001")
	pduFull, _ := hex.DecodeString("0103003c00014406")
	exc, _ := hex.DecodeString("0500")
	exc2, _ := hex.DecodeString("01900500")

	t.Run("CRCValidation", func(t *testing.T) {
		vCRC := binary.LittleEndian.Uint16([]byte{0x44, 0x06})

		assert.Equal(t, vCRC, crcNew().calc(pdu))
		assert.EqualSlice(t, pduFull, CalcCRC(pdu))
	})
	t.Run("CRCVerification", func(t *testing.T) {
		assert.Equal(t, ValidCRC(pduFull), true)
	})
	t.Run("ExceptionHandling", func(t *testing.T) {
		t.Logf("Actual Error: %#v\n", ValidRTU(exc))
		t.Logf("Actual Error: %#v\n", ValidRTU(exc2))
		assert.Equal(t, errors.Is(ValidRTU(exc), SolamranExceptionError{code: 0x05, val: solarmanErrorCodesMap[0x05]}), true)
		assert.Equal(t, errors.Is(ValidRTU(exc2), ModbusError{code: 0x10, exType: ExceptionsMap[Acknowledge]}), true)
	})
	t.Run("Validation-WriteMultipleCoils", func(t *testing.T) {
		target, _ := hex.DecodeString("110F0013000A02CD01BF0B")
		adu, err := CreateWriteMultipleCoils(17, 19, 10, []byte{0xcd, 0x01})
		if err != nil {
			t.Fatalf("Create request failed: %s\n", err.Error())
		}
		for ix, b := range target {
			if adu[ix] != b {
				t.Fatalf("Validation failed at pos [%d]: %v - %x -> target: %x\n", ix, adu, adu, target)
			}
		}
	})
	t.Run("Validation-WriteSingleCoil", func(t *testing.T) {
		target, _ := hex.DecodeString("110500ACFF004E8B")
		adu, err := CreateWriteSingleCoil(17, 172, 0xff00)
		if err != nil {
			t.Fatalf("Create request failed: %s\n", err.Error())
		}
		for ix, b := range target {
			if adu[ix] != b {
				t.Fatalf("Validation failed at pos [%d]: %v - %x -> target: %x\n", ix, adu, adu, target)
			}
		}
	})
	t.Run("Validation-ReadInputRegisters", func(t *testing.T) {
		target, _ := hex.DecodeString("110400080001B298")
		adu, err := CreateReadInputs(17, 8, 1)
		if err != nil {
			t.Fatalf("Create request failed: %s\n", err.Error())
		}
		for ix, b := range target {
			if adu[ix] != b {
				t.Fatalf("Validation failed at pos [%d]: %v - %x -> target: %x\n", ix, adu, adu, target)
			}
		}
	})
	t.Run("Validation-WriteMultipleHoldingRegisters", func(t *testing.T) {
		target, _ := hex.DecodeString("11100001000204000A0102C6F0")
		adu, err := CreateWriteMultipleRegisters(17, 1, 2, []byte{0x0, 0xa, 0x1, 0x2})
		if err != nil {
			t.Fatalf("Create request failed: %s\n", err.Error())
		}
		for ix, b := range target {
			if adu[ix] != b {
				t.Fatalf("Validation failed at pos [%d]: %v - %x -> target: %x\n", ix, adu, adu, target)
			}
		}
	})

}

func TestBigEndianEncoding(t *testing.T) {
	val := int16(-12000)
	val2 := float32(25.5)
	fmt.Printf("Neg as uint: %d -> int16 %d\n", uint16(val), int16(uint16(val)))
	fmt.Printf("Float as uint: %d -> int16 %.2f\n", uint16(val2), float32(uint16(val2)))
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint16(buf, uint16(val))
	fmt.Printf("Encoded value: %x -> %v\n", buf, buf)
}

func BenchmarkModbus(b *testing.B) {
	pdu, _ := hex.DecodeString("0103003c0001")
	for b.Loop() {
		CalcCRC(pdu)
	}
}
