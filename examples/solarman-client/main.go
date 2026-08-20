package main

import (
	"context"
	"fmt"
	"time"

	"github.com/githubDante/solarman"
	"github.com/githubDante/solarman/modbus"
)

func main() {

	c, err := solarman.NewClient(&solarman.SolarmanClientOptions{
		Address:       "192.168.1.101",
		Port:          8899,
		Serial:        1234567890,
		SlaveId:       1,
		AutoReconnect: true,
	})
	if err != nil {
		fmt.Printf("Cannot create client: %s\n", err.Error())
	}

	for range 3 {
		ctx, ctxC := context.WithTimeout(context.Background(), time.Second*1)
		res, err := c.ReadHoldingRegisters(ctx, 60, 1)
		ctxC()
		if err != nil {
			fmt.Printf("Cannot read inverter data: %s\n", err.Error())
		} else {
			data := modbus.FrameToRTUData(res)
			fmt.Printf("Received [%d] register(s) response from slave [%d]. Response values: %v",
				data.DataLen,
				data.SlaveId,
				data.Data,
			)
		}
		time.Sleep(30 * time.Second)
	}
	c.Close()
}
