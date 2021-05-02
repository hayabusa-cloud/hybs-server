package main

import (
	"encoding/binary"
	"math"
	"sync"
)

type packet struct {
	header  uint8
	length  uint16
	offset  uint16
	code    uint16
	payload []uint8

	networkByteOrder binary.ByteOrder
}

const payloadMaxLength = 0x1000

var packetPool = &sync.Pool{
	New: func() interface{} {
		return &packet{
			header:  0,
			length:  0,
			offset:  2,
			code:    0,
			payload: make([]byte, payloadMaxLength),

			networkByteOrder: binary.BigEndian,
		}
	},
}

func newPacket() *packet {
	var ret = packetPool.Get().(*packet)
	return ret.Reset()
}
func releasePacket(pkt *packet) {
	packetPool.Put(pkt)
}

func (pkt *packet) ToBytes() []byte {
	var totalLength = pkt.offset + 3
	var headerLength = []byte{pkt.header, uint8(totalLength >> 8), uint8(totalLength & 0xff)}
	var ret = append(headerLength, pkt.payload[:pkt.offset]...)
	return ret
}

func (pkt *packet) ReadBool(val *bool) *packet {
	if pkt.offset >= pkt.length {
		return pkt
	}
	*val = pkt.payload[pkt.offset] == 1
	pkt.offset++
	return pkt
}
func (pkt *packet) ReadInt8(val *int8) *packet {
	if pkt.offset >= pkt.length {
		return pkt
	}
	*val = int8(pkt.payload[pkt.offset])
	pkt.offset++
	return pkt
}
func (pkt *packet) ReadUint8(val *uint8) *packet {
	if pkt.offset >= pkt.length {
		return pkt
	}
	*val = pkt.payload[pkt.offset]
	pkt.offset++
	return pkt
}
func (pkt *packet) ReadInt16(val *int16) *packet {
	if pkt.offset+2 > pkt.length {
		return pkt
	}
	*val = int16(pkt.networkByteOrder.Uint16(pkt.payload[pkt.offset:]))
	pkt.offset += 2
	return pkt
}
func (pkt *packet) ReadUint16(val *uint16) *packet {
	if pkt.offset+2 > pkt.length {
		return pkt
	}
	*val = pkt.networkByteOrder.Uint16(pkt.payload[pkt.offset:])
	pkt.offset += 2
	return pkt
}
func (pkt *packet) ReadInt32(val *int32) *packet {
	if pkt.offset+4 > pkt.length {
		return pkt
	}
	*val = int32(pkt.networkByteOrder.Uint32(pkt.payload[pkt.offset:]))
	pkt.offset += 4
	return pkt
}
func (pkt *packet) ReadUint32(val *uint32) *packet {
	if pkt.offset+4 > pkt.length {
		return pkt
	}
	*val = pkt.networkByteOrder.Uint32(pkt.payload[pkt.offset:])
	pkt.offset += 4
	return pkt
}
func (pkt *packet) ReadInt64(val *int64) *packet {
	if pkt.offset+8 > pkt.length {
		return pkt
	}
	*val = int64(pkt.networkByteOrder.Uint64(pkt.payload[pkt.offset:]))
	pkt.offset += 8
	return pkt
}
func (pkt *packet) ReadUint64(val *uint64) *packet {
	if pkt.offset+8 > pkt.length {
		return pkt
	}
	*val = pkt.networkByteOrder.Uint64(pkt.payload[pkt.offset:])
	pkt.offset += 8
	return pkt
}
func (pkt *packet) ReadFloat32(val *float32) *packet {
	if pkt.offset+4 > pkt.length {
		return pkt
	}
	*val = math.Float32frombits(pkt.networkByteOrder.Uint32(pkt.payload[pkt.offset:]))
	pkt.offset += 4
	return pkt
}
func (pkt *packet) ReadFloat64(val *float64) *packet {
	if pkt.offset+8 > pkt.length {
		return pkt
	}
	*val = math.Float64frombits(pkt.networkByteOrder.Uint64(pkt.payload[pkt.offset:]))
	pkt.offset += 8
	return pkt
}
func (pkt *packet) ReadBytes() []byte {
	if pkt.offset+2 > pkt.length {
		return nil
	}
	var byteLen = pkt.networkByteOrder.Uint16(pkt.payload[pkt.offset:])
	pkt.offset += 2
	if pkt.offset+byteLen > pkt.length {
		return nil
	}
	var ret = pkt.payload[pkt.offset : pkt.offset+byteLen]
	pkt.offset += byteLen
	return ret
}
func (pkt *packet) ReadString() string {
	return string(pkt.ReadBytes())
}
func (pkt *packet) InPayload() []byte {
	return pkt.payload
}
func (pkt *packet) WriteBool(val bool) *packet {
	if (pkt.offset+1)&0xfff != 0 {
		return pkt
	}
	if val {
		pkt.payload[pkt.offset] = 1
	} else {
		pkt.payload[pkt.offset] = 0
	}
	pkt.offset++
	return pkt
}
func (pkt *packet) WriteInt8(val int8) *packet {
	if (pkt.offset+1)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset] = uint8(val)
	pkt.offset++
	return pkt
}
func (pkt *packet) WriteUint8(val uint8) *packet {
	if (pkt.offset+1)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset] = val
	pkt.offset++
	return pkt
}
func (pkt *packet) WriteInt16(val int16) *packet {
	if (pkt.offset+2)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	pkt.payload[pkt.offset] = uint8(val >> 8)
	pkt.offset += 2
	return pkt
}
func (pkt *packet) WriteUint16(val uint16) *packet {
	if (pkt.offset+2)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	pkt.payload[pkt.offset] = uint8(val >> 8)
	pkt.offset += 2
	return pkt
}
func (pkt *packet) WriteInt32(val int32) *packet {
	if (pkt.offset+4)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+3] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+2] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset] = uint8(val & 0xff)
	pkt.offset += 4
	return pkt
}
func (pkt *packet) WriteUint32(val uint32) *packet {
	if (pkt.offset+4)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+3] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+2] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset] = uint8(val & 0xff)
	pkt.offset += 4
	return pkt
}
func (pkt *packet) WriteInt64(val int64) *packet {
	if (pkt.offset+8)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+7] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+6] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+5] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+4] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+3] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+2] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset] = uint8(val & 0xff)
	pkt.offset += 8
	return pkt
}
func (pkt *packet) WriteUint64(val uint64) *packet {
	if (pkt.offset+8)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+7] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+6] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+5] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+4] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+3] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+2] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset+1] = uint8(val & 0xff)
	val >>= 8
	pkt.payload[pkt.offset] = uint8(val & 0xff)
	pkt.offset += 8
	return pkt
}
func (pkt *packet) WriteFloat32(val float32) *packet {
	return pkt.WriteUint32(math.Float32bits(val))
}
func (pkt *packet) WriteFloat64(val float64) *packet {
	return pkt.WriteUint64(math.Float64bits(val))
}
func (pkt *packet) WriteBytes(val []byte) *packet {
	var l = uint16(len(val))
	if (pkt.offset+1)&0xc000 != 0 {
		return pkt
	}
	pkt.payload[pkt.offset+1] = uint8(l & 0xff)
	pkt.payload[pkt.offset] = uint8(l >> 8)
	pkt.offset += 2
	copy(pkt.payload[pkt.offset:], val)
	pkt.offset += l
	return pkt
}
func (pkt *packet) WriteString(val string) *packet {
	return pkt.WriteBytes([]byte(val))
}
func (pkt *packet) WriteBytesNoLen(val []byte) *packet {
	if (pkt.offset+1)&0xc000 != 0 {
		return pkt
	}
	copy(pkt.payload[pkt.offset:], val)
	pkt.offset += uint16(len(val))
	return pkt
}
func (pkt *packet) SetHeader(header uint8) *packet {
	pkt.header = header
	return pkt
}
func (pkt *packet) SetHeader3(pv, pd, cmd uint8) *packet {
	pkt.header = pv | pd | cmd
	return pkt
}
func (pkt *packet) SetEventCode(code uint16) *packet {
	pkt.code = code
	pkt.payload[0], pkt.payload[1] = uint8(code>>8), uint8(code&0xff)
	return pkt
}
func (pkt *packet) Reset() *packet {
	pkt.offset = 2
	return pkt
}
