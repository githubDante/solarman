package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"testing"
	"time"
)

func TestProtocol(t *testing.T) {
	req, rErr := hex.DecodeString("a5170010451400d20296490200000000000000000000000000000103003c000144061115")
	if rErr != nil {
		t.Fatalf("Cannot create request: %s\n", rErr.Error())
	}
	res, rErr := hex.DecodeString("a51d001015742bd20296490201623e6106c52900007b181f6401030a32323131303735303730c32ae715")
	if rErr != nil {
		t.Fatalf("Cannot create response: %s\n", rErr.Error())
	}
	rkl, rkErr := hex.DecodeString("a5010010477227d202964900f515")
	if rkErr != nil {
		t.Fatalf("Cannot create response: %s\n", rkErr.Error())
	}
	exc, _ := hex.DecodeString("a510001015e6f478243cd6020189460100c90f000039f87e6a05008615")

	v5Req, reqErr := RawToPacket(req)
	v5Res, resErr := RawToPacket(res)
	v5Rkl, rklErr := RawToPacket(rkl)
	v5Exc, excErr := RawToPacket(exc)
	if reqErr != nil {
		log.Fatalf("Cannot create request packet: %s\n", reqErr.Error())
	}
	if resErr != nil {
		log.Fatalf("Cannot create response packet: %s\n", resErr.Error())
	}
	if rklErr != nil {
		log.Fatalf("Cannot create keep alive packet: %s\n", rklErr.Error())
	}
	if excErr != nil {
		log.Fatalf("Cannot create exception packet: %s\n", excErr.Error())
	}

	t.Logf("Request:\n%#v\n", v5Req)
	t.Logf("Request logger serial [%d]\n", v5Req.LoggerSerial())
	t.Logf("Request valid [%t]\n", v5Req.Valid())
	t.Logf("Request rtu [%x]\n", v5Req.ModbusRTU())
	t.Logf("Request raw: %x\n", v5Req.Raw())
	t.Logf("Request SeqOut: %d\n", v5Req.(*V5Request).SeqOut)
	t.Logf("Request V5Length: %d\n", v5Req.LengthV5())
	t.Logf("Response:\n%#v\n", v5Res)
	t.Logf("Response totWork T [%d]\n", v5Res.(*V5Response).responsePayload.TotWorkTime)
	t.Logf("Response powerOn T [%d] [%s] [%s]\n", v5Res.(*V5Response).responsePayload.PowerOnTime,
		(time.Duration(v5Res.(*V5Response).responsePayload.PowerOnTime) * time.Second).String(),
		time.Unix(
			int64(v5Res.(*V5Response).responsePayload.OffsetTime+
				v5Res.(*V5Response).responsePayload.TotWorkTime),
			0).String(),
	)
	t.Logf("Response offest T [%d] [%s] [%s]\n", v5Res.(*V5Response).responsePayload.OffsetTime,
		(time.Duration(v5Res.(*V5Response).responsePayload.OffsetTime) * time.Second).String(),
		time.Unix(int64(v5Res.(*V5Response).responsePayload.OffsetTime), 0).String())
	t.Logf("KeepAlive valid [%t]\n", v5Rkl.Valid())
	t.Logf("KeepAlive serial [%d]\n", v5Rkl.LoggerSerial())
	t.Logf("KeepAlive raw [%x]\n", v5Rkl.Raw())
	t.Logf("Exception raw [%x]", v5Exc.Raw())
	t.Logf("Exception PDU [%x]", v5Exc.ModbusRTU())
	t.Logf("Exception Type [%s]", v5Exc.TypeName())

	t.Run("PacketGen", func(t *testing.T) {
		wr := NewV5Wrapper(2718270936)
		wr.seqNo = 20
		pkt := wr.WrapRequest([]uint8{0x1, 0x3, 0x0, 0x3c, 0x0, 0x1, 0x44, 0x6})
		t.Logf("Generated Request:\n%#v\n", pkt)
		t.Logf("Generated Request valid [%t]\n", pkt.Valid())
		t.Logf("Generated Request length [%d]\n", pkt.LengthV5())
	})

}

func TestSerial(t *testing.T) {
	buf := make([]byte, 0)
	buf2 := make([]byte, 0)
	buf, _ = binary.Append(buf, binary.LittleEndian, uint32(2718270936))
	buf2, _ = binary.Append(buf2, binary.LittleEndian, uint32(1234567890))
	t.Logf("[9876543210] as hex: %x\n", buf)
	t.Logf("[1234567890] as hex: %x\n", buf2)
}
