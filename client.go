package solarman

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	sollog "github.com/githubDante/solarman/log"
	"github.com/githubDante/solarman/modbus"
	"github.com/githubDante/solarman/protocol"
)

var log sollog.Logger

func init() {
	log = &sollog.SilentLoger{}
}

// SetLogger logger for the solarman package
//
// The dafult logger is SilentLogger i.e. nothing will be logged/printed
func SetLogger(l sollog.Logger) {
	log = l
}

func createConnection(address string, port int) (net.Conn, error) {
	if port == 0 {
		port = defaultOptions.Port
	}
	return net.Dial("tcp4", fmt.Sprintf("%s:%d", address, port))
}

// Client is a V5 compatible a
type Client struct {
	c               net.Conn
	opt             *SolarmanClientOptions
	proto           *protocol.Wrapper
	awaitingPkt     bool
	comm            chan protocol.V5Packet
	mu              sync.Mutex
	closedByRequest bool
}

// NewClient client initialization
func NewClient(opts *SolarmanClientOptions) (*Client, error) {
	if c, cErr := createConnection(opts.Address, opts.Port); cErr != nil {
		return nil, cErr
	} else {
		cl := new(Client)
		cl.opt = opts
		cl.c = c
		cl.proto = protocol.NewV5Wrapper(cl.opt.Serial)
		cl.comm = make(chan protocol.V5Packet)
		cl.mu = sync.Mutex{}
		go cl.readLoop()
		return cl, nil
	}
}

func (c *Client) sendAndWait(ctx context.Context, req protocol.V5Packet) (protocol.V5Packet, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.c == nil {
		if !c.opt.AutoReconnect {
			return nil, fmt.Errorf("client not connected")
		}
		if c.closedByRequest {
			return nil, fmt.Errorf("client connection closed by request")
		}
		log.Warn("Will wait for client reconnect...")
		t := time.NewTicker(5 * time.Millisecond)
	reconnectWait:
		for {
			select {
			case <-t.C:
				if c.c == nil {
					continue
				} else {
					break reconnectWait
				}
			case <-ctx.Done():
				return nil, fmt.Errorf("cannot send packet: %s", ctx.Err().Error())
			}
		}
	}
	if _, wErr := c.c.Write(req.Raw()); wErr != nil {
		log.Errorf("[%d] Socket write failure: %s", wErr.Error())
		return nil, wErr
	}
	c.awaitingPkt = true
	defer func() { c.awaitingPkt = false }()
	log.Debugf("[%d] V5 packet sent: %x", c.opt.Serial, req.Raw())
	for {
		select {
		case pkt := <-c.comm:
			if c.proto.ReponseMatch(req, pkt) {
				return pkt, nil
			} else {
				continue
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("packet not received: %s", ctx.Err().Error())
		}
	}
}

func (c *Client) readLoop() {
	log.Debugf("[%d] readloop started", c.opt.Serial)
	for {
		pkt := make([]byte, 8192)
		n, e := c.c.Read(pkt)
		if e != nil {
			log.Errorf("[%d] read error: %s", c.opt.Serial, e.Error())
			log.Errorf("[%d] read loop exiting...", c.opt.Serial)
			break
		}
		log.Debugf("[%d] got packet: %x", c.opt.Serial, pkt[:n])
		v5Packet, pktEr := protocol.RawToPacket(pkt[:n])
		if pktEr != nil {
			log.Errorf("[%d] cannot decode packet: %s", c.opt.Serial, pktEr.Error())
		} else {
			if c.awaitingPkt {
				c.comm <- v5Packet
			}
		}
	}
	c.reconnect()
}

func (c *Client) reconnect() {
	if !c.opt.AutoReconnect || c.closedByRequest {
		log.Warnf("Reconnect disabled. Exiting...")
		return
	}
	log.Infof("Reconnect enabled. Trying to connect...")
	for {
		time.Sleep(1 * time.Second)
		if conn, cErr := createConnection(c.opt.Address, c.opt.Port); cErr == nil {
			c.c = conn
			go c.readLoop()
			break
		}
	}
}

// Close disconencts the client from the datalogger.
//
// This disables the auto-reconnect functionality making the client unusable after this point
func (c *Client) Close() {
	if c.c == nil || c.closedByRequest {
		log.Warnf("[%d] Connection already closed...", c.opt.Serial)
		return
	}
	log.Infof("[%d] Connection will be closed...", c.opt.Serial)
	c.closedByRequest = true
	c.c.Close()
	c.c = nil
}

// ReadInputRegisters
//
//	MODBUS function (0x04)
//
// Read data from the inverter starting from address `startAddr`
func (c *Client) ReadInputRegisters(ctx context.Context, startAddr, quantity uint16) ([]byte, error) {
	adu, mErr := modbus.CreateReadInputs(c.opt.SlaveId, startAddr, quantity)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// WriteSingleCoil
//
//	MODBUS function (0x05)
func (c *Client) WriteSingleCoil(ctx context.Context, startAddr, value uint16) ([]byte, error) {
	adu, mErr := modbus.CreateWriteSingleCoil(c.opt.SlaveId, startAddr, value)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// WriteMultipleCoils
//
//	MODBUS function (0x0f)
func (c *Client) WriteMultipleCoils(ctx context.Context, startAddr, quantity uint16, data []byte) ([]byte, error) {
	adu, mErr := modbus.CreateWriteMultipleCoils(c.opt.SlaveId, startAddr, quantity, data)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// ReadHoldingRegisters
//
//	MODBUS function (0x03)
//
// Will send request for data to the datalogger and will return
// the reposnse MODBUS RTU part after validation
func (c *Client) ReadHoldingRegisters(ctx context.Context, startAddr, quantity uint16) ([]byte, error) {
	adu, mErr := modbus.CreateReadHolding(c.opt.SlaveId, startAddr, quantity)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// WriteHoldingRegister
//
//	MODBUS function (0x06)
//
// Write of single holding register.
// The returned response is a validated MODBUS RTU
func (c *Client) WriteHoldingRegister(ctx context.Context, startAddress, value uint16) ([]byte, error) {
	adu, mErr := modbus.CreateWriteSingleRegister(c.opt.SlaveId, startAddress, value)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// WriteMultipleHoldingRegisters function (0x10)
//
// The data should be a binary encoded uint16 slice
func (c *Client) WriteMultipleHoldingRegisters(ctx context.Context, startAddress, quantity uint16, data []byte) ([]byte, error) {
	adu, mErr := modbus.CreateWriteMultipleRegisters(c.opt.SlaveId, startAddress, quantity, data)
	if mErr != nil {
		return nil, mErr
	}
	pkt := c.proto.WrapRequest(adu)
	res, err := c.sendAndWait(ctx, pkt)
	if err != nil {
		return nil, err
	}
	return validateAndReturn(res)
}

// validateAndReturn Response validator and RTU extractor
func validateAndReturn(res protocol.V5Packet) ([]byte, error) {
	if !res.Valid() {
		return nil, fmt.Errorf("invalid V5 packet received")
	}
	if mErr := modbus.ValidRTU(res.ModbusRTU()); mErr != nil {
		return nil, mErr
	}
	return res.ModbusRTU(), nil
}
