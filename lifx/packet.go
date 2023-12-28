package lifx

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
)

type PaylodType uint16

//goland:noinspection ALL
const (
	Undefined PaylodType = iota

	// Device
	GetService        = 2
	StateService      = 3
	GetHostFirmware   = 14
	StateHostFirmware = 15
	GetWifiInfo       = 16
	StateWifiInfo     = 17
	GetWifiFirmware   = 18
	StateWifiFirmware = 19
	GetPower          = 20
	SetPower          = 21
	StatePower        = 22
	GetLabel          = 23
	SetLabel          = 24
	StateLabel        = 25
	GetVersion        = 32
	StateVersion      = 33
	GetInfo           = 34
	StateInfo         = 35
	SetReboot         = 38
	Acknowledgement   = 45
	GetLocation       = 48
	SetLocation       = 49
	StateLocation     = 50
	GetGroup          = 51
	SetGroup          = 52
	StateGroup        = 53
	EchoRequest       = 58
	EchoResponse      = 59

	//Light
	GetColor                   = 101
	SetColor                   = 102
	SetWaveform                = 103
	LightState                 = 107
	GetLightPower              = 116
	SetLightPower              = 117
	StateLightPower            = 118
	SetWaveformOptional        = 119
	GetInfrared                = 120
	StateInfrared              = 121
	SetInfrared                = 122
	GetHevCycle                = 142
	SetHevCycle                = 143
	StateHevCycle              = 144
	GetHevCycleConfiguration   = 145
	SetHevCycleConfiguration   = 146
	StateHevCycleConfiguration = 147
	GetLastHevCycleResult      = 148
	StateLastHevCycleResult    = 149
	StateUnhandled             = 223

	// mutlizone
	SetColorZones           = 501
	GetColorZones           = 502
	StateZone               = 503
	StateMultiZone          = 506
	GetMultiZoneEffect      = 507
	SetMultiZoneEffect      = 508
	StateMultiZoneEffect    = 509
	SetExtendedColorZones   = 510
	GetExtendedColorZones   = 511
	StateExtendedColorZones = 512

	// Relay
	GetRPower   = 816
	SetRPower   = 817
	StateRPower = 818

	// Tile
	GetDeviceChain   = 701
	StateDeviceChain = 702
	SetUserPosition  = 703
	Get64            = 707
	State64          = 711
	Set64            = 715
	GetTileEffect    = 718
	SetTileEffect    = 719
	StateTileEffect  = 720
)

var packetSource = uint32(2)

//goland:noinspection GoNameStartsWithPackageName
type LIFXPacket struct {
	// common items
	Target           net.HardwareAddr
	PayloadTypeValue PaylodType
	ResponseRequired bool
	AckRequired      bool
	// power
	PowerLevel uint16
	// color
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
	// common - used to indicate how quickly a change should happen (in milliseconds
	// TODO - change to time
	Duration uint32
}

func NewSetColorLIFXPacket(hue float64, saturation float64, brightness float64, kelvin uint16, duration uint32) *LIFXPacket {

	p := &LIFXPacket{Hue: uint16(math.Round(0xFFFF * hue)),
		Saturation: uint16(math.Round(0xFFFF * saturation)),
		Brightness: uint16(math.Round(0xFFFF * brightness)),
		Kelvin:     kelvin,
		Duration:   duration}
	p.PayloadTypeValue = SetColor
	p.Target, _ = net.ParseMAC("d0:73:d5:00:13:37")
	p.AckRequired = true
	return p
}

const hexDigit = "0123456789abcdef"

func (p LIFXPacket) Dump() {
	buffer := p.Marshal()
	s := ""
	for i, b := range buffer {
		if i%8 == 0 {
			s = s + fmt.Sprintf("%04x: ", i)
		}
		s = s + fmt.Sprintf("%c%c", hexDigit[b>>4], hexDigit[b&0x0f])
		if i%8 == 7 {
			s = s + fmt.Sprintf("\n")
		}
	}
	s = s + fmt.Sprintf("\n")
	print(s)
}

func (p LIFXPacket) Marshal() []byte {
	buffer := make([]byte, 1024)

	//
	//   ========== Frame Header
	//
	bufferIndex := uint16(0)

	// 0..1 is LIFXPacket size.  This will be computed at the end
	bufferIndex++
	bufferIndex++

	buffer[2] = 0
	bufferIndex++
	// protocol must be 1024
	buffer[3] = 0x04
	// addressable must be 1
	buffer[3] |= 0x10
	/*
		if tagged
			buffer[3] |= 0x20
		// x40 and x80 reserved and must be zero
	*/
	bufferIndex++

	// 4..7 A unique key to identify replies to this LIFXPacket
	binary.LittleEndian.PutUint32(buffer[bufferIndex:], packetSource)
	bufferIndex = bufferIndex + 4

	//
	//   ========== Frame Address
	//

	// 8..13 is a mac address when using broadcast to specify a single target
	// or all zeros for all receiving LIFX devices.
	if p.Target != nil {
		copy(buffer[bufferIndex:], p.Target[0:6])
		bufferIndex = bufferIndex + 6
	} else {
		for i := 0; i < 6; i++ {
			buffer[bufferIndex] = 0
			bufferIndex++
		}
	}
	// The specification then provides 2 extra bytes (14..15) that must
	// be zero.  Perhaps the extra bytes are for future EUI-64 support?
	buffer[bufferIndex] = 0
	bufferIndex++
	buffer[bufferIndex] = 0
	bufferIndex++

	// 16..21 are reserved for now and always set to zero
	for i := 0; i < 6; i++ {
		buffer[bufferIndex] = 0
		bufferIndex++
	}

	// 22 LIFXPacket boolean
	buffer[bufferIndex] = 0

	if p.ResponseRequired {
		buffer[bufferIndex] |= 0x01
	}
	if p.AckRequired {
		buffer[bufferIndex] |= 0x02
	}
	// remaining bits are not used
	bufferIndex++

	// 23
	// 	uint8_t  sequence;
	buffer[bufferIndex] = 1
	bufferIndex++

	// 24..31 eight reserved bytes
	for i := 0; i < 8; i++ {
		buffer[bufferIndex] = 0
		bufferIndex++
	}

	//
	//   ========== Protocol Header
	//

	// 32..33 the type of LIFXPacket.
	binary.LittleEndian.PutUint16(buffer[bufferIndex:], uint16(p.PayloadTypeValue))
	bufferIndex += 2

	// 34..35 two reserved bytes
	for i := 0; i < 2; i++ {
		buffer[bufferIndex] = 0
		bufferIndex++
	}

	switch p.PayloadTypeValue {
	case SetPower:
		binary.LittleEndian.PutUint16(buffer[bufferIndex:], p.PowerLevel)
		bufferIndex = bufferIndex + 2
		binary.LittleEndian.PutUint32(buffer[bufferIndex:], p.Duration)
		bufferIndex = bufferIndex + 4
	case SetColor:
		// reserved
		buffer[bufferIndex] = 0
		bufferIndex++
		binary.LittleEndian.PutUint16(buffer[bufferIndex:], p.Hue)
		bufferIndex += 2
		binary.LittleEndian.PutUint16(buffer[bufferIndex:], p.Saturation)
		bufferIndex += 2
		binary.LittleEndian.PutUint16(buffer[bufferIndex:], p.Brightness)
		bufferIndex += 2
		binary.LittleEndian.PutUint16(buffer[bufferIndex:], p.Kelvin)
		bufferIndex += 2
		binary.LittleEndian.PutUint32(buffer[bufferIndex:], p.Duration)
		bufferIndex += 4
	default:
		log.Fatalf("Unable to Marshall packets of PayloadTypeValue %d (not implemented)", p.PayloadTypeValue)
	}

	binary.LittleEndian.PutUint16(buffer[0:], bufferIndex)
	return buffer[0:bufferIndex]

	/* variable length payload follows */
	/*
	   typedef struct {
	     uint16_t power;
	   } payload_StateLightPower;


	   typedef struct {
	     uint8_t reserved;
	     uint16_t hue;
	     uint16_t saturation;
	     uint16_t brightness;
	     uint16_t kelvin;
	     uint32_t duration;

	   } payload_SetColor;

	   typedef struct {
	     uint16_t hue;
	     uint16_t saturation;
	     uint16_t brightness;
	     uint16_t kelvin;
	     uint16_t power;
	     char label[32];
	     char reserved[8];
	   } payload_LightState;


	*/
	return buffer
}
