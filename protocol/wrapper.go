package protocol

import (
	"math/rand"
)

type Wrapper struct {
	serial uint32
	seqNo  uint8
}

// NewV5Wrapper SolarmanV5 wrpaper initializtion
func NewV5Wrapper(serial uint32) *Wrapper {
	w := new(Wrapper)
	w.serial = serial
	w.seqNo = uint8(rand.Int31n(255))
	return w
}

// WrapRequest will create SolarmanV5 request packet from a raw MODBUS ADU
func (w *Wrapper) WrapRequest(adu []byte) V5Packet {
	defer func() { w.seqNo++ }()
	r := new(V5Request)
	r.v5Common.Header = V5Start
	r.v5Common.LoggerSn = w.serial
	r.v5Common.SeqOut = w.seqNo
	r.v5Common.CCode = uint16(V5ModbusRequestPacket)
	r.requestPayload.FrameType = V5FrameType

	r.RTU = adu
	r.v5Trailer.End = V5End

	r.v5Common.Length = v5LengthCalc(r)
	r.v5Trailer.Checksum = v5Checksum(r)
	return r
}

// WrapResponse will create SolarmanV5 response packet from a raw MODBUS ADU
func (w *Wrapper) WrapResponse(adu []byte, req V5Packet) V5Packet {
	r := new(V5Response)
	rq := req.(*V5Request)
	if rq == nil {
		return r
	}
	r.v5Common.Header = V5Start
	r.v5Common.LoggerSn = w.serial
	r.v5Common.SeqIn = w.seqNo
	r.v5Common.SeqOut = rq.SeqOut
	r.v5Common.CCode = uint16(V5ModbusResponsePacket)
	r.responsePayload.FrameType = V5FrameType

	r.RTU = adu
	r.v5Trailer.End = V5End

	r.v5Common.Length = v5LengthCalc(r)
	r.v5Trailer.Checksum = v5Checksum(r)

	return r
}

// ReponseMatch checks a SolarmanV5 response packet against request wrapped previously by this wrapper
//
// It should be noted that the Sequence number is a single byte, so false positivies are possible
func (w *Wrapper) ReponseMatch(req, res V5Packet) bool {
	rq := req.(*V5Request)
	rs := res.(*V5Response)
	if rq == nil || rs == nil {
		return false
	}
	return rq.SeqOut == rs.SeqOut
}
