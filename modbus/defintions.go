package modbus

type ModbusFn uint8

func (m ModbusFn) Uint8() uint8 {
	return uint8(m)
}

const (
	ReadCoilsFn              ModbusFn = 0x01
	ReadDicreteInputsFn      ModbusFn = 0x02
	ReadHoldingRegistersFn   ModbusFn = 0x03
	ReadInputRegistersFn     ModbusFn = 0x04
	WriteSingleCoilFn        ModbusFn = 0x05
	WriteSingleRegisterFn    ModbusFn = 0x06
	WriteMultipleCoilsFn     ModbusFn = 0x0f
	WriteMultipleRegistersFn ModbusFn = 0x10
)

// ModbusExceptionBase - exceptions are returned in format
//
//	exception function code + 0x80
//	exception code
const ModbusExceptionBase uint8 = 0x80

const (
	IllegalFunctionException           uint8 = 1
	IllegalDataAddress                 uint8 = 2
	IllegalDataValue                   uint8 = 3
	ServerServiceFailure               uint8 = 4
	Acknowledge                        uint8 = 5
	ServerDeviceBusy                   uint8 = 6
	MemoryParityError                  uint8 = 8
	GatewayParityError                 uint8 = 10
	GatewayTargetDeviceFailedToRespond uint8 = 11
)

var ExceptionsMap = map[uint8]string{
	IllegalFunctionException:           "IlleIllegalFunctionException",
	IllegalDataAddress:                 "IllegalDataAddress",
	IllegalDataValue:                   "IllegalDataValue",
	ServerServiceFailure:               "ServerServiceFailure",
	Acknowledge:                        "Acknowledge",
	ServerDeviceBusy:                   "ServerDeviceBusy",
	MemoryParityError:                  "MemoryParityError",
	GatewayParityError:                 "GatewayParityError",
	GatewayTargetDeviceFailedToRespond: "GatewayTargetDeviceFailedToRespond",
}
