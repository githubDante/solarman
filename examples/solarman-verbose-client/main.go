package main

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/githubDante/solarman"
	"github.com/githubDante/solarman/log"
	"github.com/githubDante/solarman/modbus"
)

func main() {
	l := log.NewSolarmanLogger()
	solarman.SetLogger(l)
	c, err := solarman.NewClient(&solarman.SolarmanClientOptions{
		Address:       "192.168.1.101",
		Port:          8899,
		Serial:        1234567890,
		SlaveId:       1,
		AutoReconnect: true,
	})
	if err != nil {
		l.Errorf("Cannot create client: %s\n", err.Error())
	}

	for range 3 {
		ctx, ctxC := context.WithTimeout(context.Background(), time.Second*3)
		buf := make([]byte, 0)
		buf = binary.BigEndian.AppendUint16(buf, 2500)
		wRes, wErr := c.WriteMultipleHoldingRegisters(ctx, 155, 1, buf)
		if wErr != nil {
			l.Errorf("Cannot write to logger: %s", wErr, err)
		} else {
			f := modbus.FrameToRTUData(wRes)
			l.Infof("Write successfull. Received response: %x", f.Data)
		}
		res, err := c.ReadHoldingRegisters(ctx, 151, 1)
		ctxC()
		if err != nil {
			l.Errorf("Cannot read inverter data: %s\n", err.Error())
		} else {
			data := modbus.FrameToRTUData(res)
			l.Infof("Received [%d] register(s) response from slave [%d]. Response values: %v",
				data.DataLen,
				data.SlaveId,
				data.Data,
			)
		}
		time.Sleep(30 * time.Second)
	}
	c.Close()
}
