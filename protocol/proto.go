package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type PacketType uint16

const (
	V5ModbusResponsePacket  PacketType = 0x1510
	V5ModbusRequestPacket   PacketType = 0x4510
	V5KeepAlivePacket       PacketType = 0x4710
	V5LoggerCloudInit       PacketType = 0x4110
	V5LoggerCloudReport     PacketType = 0x4210
	V5LoggerCloudWifiReport PacketType = 0x4310
	V5LoggerCloudUnknReport PacketType = 0x4810

	V5CloudServerIAck PacketType = 0x1110 // Bitwise AND NOT (&^) of the V5LoggerCloudInit
	V5CloudServerRAck PacketType = 0x1210 // Bitwise AND NOT (&^) of the V5LoggerCloudReport
	V5CloudServerWAck PacketType = 0x1310
	V5CloudServerUAck PacketType = 0x1810
)

var packetName = map[PacketType]string{
	V5ModbusRequestPacket:   "RequestPacket",
	V5ModbusResponsePacket:  "ResponsePacket",
	V5KeepAlivePacket:       "KeepAlive/TimeSync",
	V5LoggerCloudReport:     "LoggerCloudReport",
	V5LoggerCloudWifiReport: "LoggerCloudWiFiReport",
	V5LoggerCloudUnknReport: "LoggerCloudUnknownReport",
	V5LoggerCloudInit:       "LoggerCloudInit",
	V5CloudServerIAck:       "CloudSrvInitAck",
	V5CloudServerRAck:       "CloudSrvReportAck",
	V5CloudServerWAck:       "CloudSrvWiFiAck",
	V5CloudServerUAck:       "CloudSrvUnknownReportAck",
}

const (
	V5Start     uint8 = 0xa5
	V5End       uint8 = 0x15
	V5FrameType uint8 = 0x02
)

// v5Common is the SolarmanV5 header of the packet
type v5Common struct {
	Header   uint8
	Length   uint16 // V5 Packet length - length wihtout header & trailer
	CCode    uint16
	SeqOut   uint8
	SeqIn    uint8
	LoggerSn uint32
}

type requestPayload struct {
	FrameType   byte
	SensorType  uint16
	TotWorkTime uint32
	PowerOnTime uint32
	OffsetTime  uint32
}

type responsePayload struct {
	FrameType   uint8
	Status      uint8
	TotWorkTime uint32
	PowerOnTime uint32
	OffsetTime  uint32
}

type v5Trailer struct {
	Checksum uint8
	End      uint8
}

type V5Request struct {
	v5Common
	requestPayload
	RTU []byte
	v5Trailer
}

type V5Response struct {
	v5Common
	responsePayload
	RTU []byte
	v5Trailer
}

type V5KeepAlive struct {
	v5Common
	End uint8
}

type V5CloudAck struct {
	v5Common
	Ack []byte
	v5Trailer
}

type V5Packet interface {
	// LoggerSerial returns the serial number of the datalogger embeded in the V5 packet
	LoggerSerial() uint32
	// MODBUS RTU part of the packet
	ModbusRTU() []byte
	// Valid indicates whether this is a valid SolarmanV5 packet
	//
	// The validation incldues header, trailer & V5 checksum
	Valid() bool
	// Length
	LengthV5() uint16
	// Checksum returens the V5 checksum of the packet (if available)
	//
	// The value is missing from the KeppAlive type
	ChecksumV5() uint8
	// Raw returns encoded V5 packet
	Raw() []byte
	// Type
	Type() PacketType
	// TypeName packet type as human readable string
	TypeName() string
	// SequenceNum returns the sequence number used by solraman
	SequenceNum() uint8
}

func (rq *V5Request) LoggerSerial() uint32 {
	return rq.LoggerSn
}

func (rq *V5Request) ModbusRTU() []byte {
	return rq.RTU
}

func (rq *V5Request) LengthV5() uint16 {
	return rq.Length
}

func (rq *V5Request) ChecksumV5() uint8 {
	return rq.Checksum
}

func (rq *V5Request) Raw() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rq.v5Common)
	binary.Write(buf, binary.LittleEndian, rq.requestPayload)
	binary.Write(buf, binary.BigEndian, rq.RTU)
	binary.Write(buf, binary.LittleEndian, rq.v5Trailer)
	return buf.Bytes()
}

func (rq *V5Request) Valid() bool {
	return rq.Header == V5Start && rq.v5Trailer.End == V5End && v5Checksum(rq) == rq.v5Trailer.Checksum
}

func (rq *V5Request) Type() PacketType {
	return PacketType(rq.CCode)
}

func (rq *V5Request) TypeName() string {
	return packetName[rq.Type()]
}

func (rq *V5Request) SequenceNum() uint8 {
	return rq.SeqOut
}

func (rs *V5Response) LoggerSerial() uint32 {
	return rs.LoggerSn
}

func (rs *V5Response) ModbusRTU() []byte {
	return rs.RTU
}

func (rs *V5Response) LengthV5() uint16 {
	return rs.Length
}

func (rs *V5Response) ChecksumV5() uint8 {
	return rs.Checksum
}

func (rs *V5Response) Raw() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rs.v5Common)
	binary.Write(buf, binary.LittleEndian, rs.responsePayload)
	binary.Write(buf, binary.BigEndian, rs.RTU)
	binary.Write(buf, binary.LittleEndian, rs.v5Trailer)
	return buf.Bytes()
}

func (rs *V5Response) Valid() bool {
	return rs.Header == V5Start && rs.v5Trailer.End == V5End && v5Checksum(rs) == rs.v5Trailer.Checksum
}

func (rs *V5Response) Type() PacketType {
	return PacketType(rs.CCode)
}

func (rs *V5Response) TypeName() string {
	return packetName[rs.Type()]
}

func (rs *V5Response) SequenceNum() uint8 {
	return rs.SeqOut
}

func (kl *V5KeepAlive) LoggerSerial() uint32 {
	return kl.LoggerSn
}

func (kl *V5KeepAlive) ModbusRTU() []byte {
	return []byte(nil)
}

func (kl *V5KeepAlive) LengthV5() uint16 {
	return kl.Length
}

func (kl *V5KeepAlive) ChecksumV5() uint8 {
	return 0
}

func (kl *V5KeepAlive) Raw() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, kl.v5Common)
	binary.Write(buf, binary.LittleEndian, kl.End)
	return buf.Bytes()
}

func (kl *V5KeepAlive) Valid() bool {
	return kl.Header == V5Start && kl.End == V5End
}

func (kl *V5KeepAlive) Type() PacketType {
	return PacketType(kl.CCode)
}

func (kl *V5KeepAlive) TypeName() string {
	return packetName[kl.Type()]
}

func (kl *V5KeepAlive) SequenceNum() uint8 {
	return kl.SeqOut
}

func (ca *V5CloudAck) LoggerSerial() uint32 {
	return ca.LoggerSn
}

func (ca *V5CloudAck) ModbusRTU() []byte {
	return ca.Ack
}

func (ca *V5CloudAck) LengthV5() uint16 {
	return ca.Length
}

func (ca *V5CloudAck) ChecksumV5() uint8 {
	return ca.Checksum
}

func (ca *V5CloudAck) Raw() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ca.v5Common)
	binary.Write(buf, binary.BigEndian, ca.Ack)
	binary.Write(buf, binary.LittleEndian, ca.v5Trailer)
	return buf.Bytes()
}

func (ca *V5CloudAck) Valid() bool {
	return ca.Header == V5Start && ca.End == V5End && v5Checksum(ca) == ca.v5Trailer.Checksum
}

func (ca *V5CloudAck) Type() PacketType {
	return PacketType(ca.CCode)
}

func (ca *V5CloudAck) TypeName() string {
	return packetName[ca.Type()]
}

func (ca *V5CloudAck) SequenceNum() uint8 {
	return ca.SeqOut
}

// RawToPacket
//
// This is a parser for raw bytes -> SolarmanV5 packet
func RawToPacket(payload []byte) (V5Packet, error) {

	if len(payload) < 5 {
		return nil, errors.New("packet too short")
	}
	var CCode uint16
	var err error
	CCode = binary.LittleEndian.Uint16(payload[3:5])

	switch PacketType(CCode) {
	case V5ModbusRequestPacket:
		pkt := new(V5Request)
		buf := bytes.NewReader(payload[:len(payload)-2])
		err = binary.Read(buf, binary.LittleEndian, &pkt.v5Common)
		if err != nil {
			return nil, err
		}
		err = binary.Read(buf, binary.LittleEndian, &pkt.requestPayload)
		if err != nil {
			return nil, err
		}
		pkt.RTU = make([]byte, buf.Len())
		err = binary.Read(buf, binary.BigEndian, &pkt.RTU)
		if err != nil {
			return nil, err
		}
		pkt.v5Trailer.Checksum = payload[len(payload)-2]
		pkt.v5Trailer.End = payload[len(payload)-1]
		return pkt, nil
	case V5ModbusResponsePacket, // Logger responses V5 & Cloud
		V5LoggerCloudReport,
		V5LoggerCloudInit,
		V5LoggerCloudWifiReport,
		V5LoggerCloudUnknReport:

		pkt := new(V5Response)
		buf := bytes.NewReader(payload[:len(payload)-2])
		err = binary.Read(buf, binary.LittleEndian, &pkt.v5Common)
		if err != nil {
			return nil, err
		}
		err = binary.Read(buf, binary.LittleEndian, &pkt.responsePayload)
		if err != nil {
			return nil, err
		}
		pkt.RTU = make([]byte, buf.Len())
		err = binary.Read(buf, binary.BigEndian, &pkt.RTU)
		if err != nil {
			return nil, err
		}
		pkt.v5Trailer.Checksum = payload[len(payload)-2]
		pkt.v5Trailer.End = payload[len(payload)-1]
		return pkt, nil
	case V5KeepAlivePacket:
		pkt := new(V5KeepAlive)
		buf := bytes.NewReader(payload[:len(payload)-1])
		err = binary.Read(buf, binary.LittleEndian, &pkt.v5Common)
		if err != nil {
			return nil, err
		}
		pkt.End = payload[len(payload)-1]
		return pkt, nil
	case V5CloudServerIAck, // Cloud Server
		V5CloudServerRAck,
		V5CloudServerWAck,
		V5CloudServerUAck:

		pkt := new(V5CloudAck)
		buf := bytes.NewReader(payload)
		err = binary.Read(buf, binary.LittleEndian, &pkt.v5Common)
		if err != nil {
			return nil, err
		}
		pkt.Ack = make([]byte, buf.Len()-2)
		err = binary.Read(buf, binary.BigEndian, &pkt.Ack)
		if err != nil {
			return nil, err
		}
		err = binary.Read(buf, binary.BigEndian, &pkt.v5Trailer)
		if err != nil {
			return nil, err
		}
		return pkt, nil

	default:
		return nil, errors.New("unsupported packet type")

	}
}

func v5Checksum(packet V5Packet) uint8 {
	var checksum uint8
	pkt := packet.Raw()
	for i := 1; i < len(pkt)-2; i++ {
		checksum += pkt[i] & 0xff
	}
	return checksum
}

func v5LengthCalc(pkt V5Packet) uint16 {
	b := pkt.Raw()
	var hdr v5Common
	var trl v5Trailer
	return uint16(len(b) - binary.Size(&trl) - binary.Size(&hdr))
}

// MinPacketLength returns the length of the V5 packet Header & Trailer
func MinPacketLength() int {
	var hdr v5Common
	var trl v5Trailer
	return binary.Size(&trl) + binary.Size(&hdr)
}
