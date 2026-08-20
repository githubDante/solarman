package solarman

type SolarmanClientOptions struct {
	// Address is the IP address of the datalogger
	Address string
	// Port used as listener by the datalogger/proxy
	//
	// If not set the default one (8899) will be used
	Port int
	// Serial serial number of the datalogger
	Serial uint32
	// SlaveId of the inverter usually 1
	SlaveId uint8
	// AutoReconnect will force the client will to attempt reconenct on connection lost
	AutoReconnect bool
}

var defaultOptions = &SolarmanClientOptions{
	Address:       "127.0.0.1",
	Port:          8899,
	Serial:        1234567890,
	AutoReconnect: false,
	SlaveId:       1,
}
