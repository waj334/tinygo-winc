package winc

import "C"
import (
	"bytes"
	"encoding/binary"
)

type ECPoint struct {
	XY           [64]byte
	Size         uint16
	PrivateKeyID uint16

	// 68 bytes
}

func (e *ECPoint) read(buf []byte) {
	reader := bytes.NewReader(buf)
	reader.Read(e.XY[:])
	binary.Read(reader, binary.LittleEndian, &e.Size)
	binary.Read(reader, binary.LittleEndian, &e.PrivateKeyID)
}

func (e *ECPoint) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 68))

	buf.Write(e.XY[:])
	binary.Write(buf, binary.LittleEndian, e.Size)
	binary.Write(buf, binary.LittleEndian, e.PrivateKeyID)

	return buf.Bytes()
}

type EcdhReqInfo struct {
	REQ       uint16
	Status    uint16
	UserData  uint32
	SeqNo     uint32
	PublicKey ECPoint
	Key       [32]byte

	// 112 bytes
}

func (e *EcdhReqInfo) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &e.REQ)
	binary.Read(reader, binary.LittleEndian, &e.Status)
	binary.Read(reader, binary.LittleEndian, &e.UserData)
	binary.Read(reader, binary.LittleEndian, &e.SeqNo)

	e.PublicKey.read(buf[12:])
	reader.ReadAt(e.Key[:], 80)

}

func (e *EcdhReqInfo) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 112))

	binary.Write(buf, binary.LittleEndian, e.REQ)
	binary.Write(buf, binary.LittleEndian, e.Status)
	binary.Write(buf, binary.LittleEndian, e.UserData)
	binary.Write(buf, binary.LittleEndian, e.SeqNo)
	buf.Write(e.PublicKey.bytes())
	buf.Write(e.Key[:])

	return buf.Bytes()
}

// EcdsaVerifyReqInfo as declared in include/ecc_types.h:221
type EcdsaVerifyReqInfo struct {
	REQ       uint16
	Status    uint16
	UserData  uint32
	SeqNo     uint32
	Signature uint32

	// 16 bytes
}

func (e *EcdsaVerifyReqInfo) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &e.REQ)
	binary.Read(reader, binary.LittleEndian, &e.Status)
	binary.Read(reader, binary.LittleEndian, &e.UserData)
	binary.Read(reader, binary.LittleEndian, &e.SeqNo)
	binary.Read(reader, binary.LittleEndian, &e.Signature)

}

func (e *EcdsaVerifyReqInfo) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 118))

	binary.Write(buf, binary.LittleEndian, e.REQ)
	binary.Write(buf, binary.LittleEndian, e.Status)
	binary.Write(buf, binary.LittleEndian, e.UserData)
	binary.Write(buf, binary.LittleEndian, e.SeqNo)
	binary.Write(buf, binary.LittleEndian, e.Signature)

	return buf.Bytes()
}

// EcdsaSignReqInfo as declared in include/ecc_types.h:236
type EcdsaSignReqInfo struct {
	REQ       uint16
	Status    uint16
	UserData  uint32
	SeqNo     uint32
	CurveType EcNamedCurve
	HashSize  uint16

	// 16 bytes
}

func (e *EcdsaSignReqInfo) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &e.REQ)
	binary.Read(reader, binary.LittleEndian, &e.Status)
	binary.Read(reader, binary.LittleEndian, &e.UserData)
	binary.Read(reader, binary.LittleEndian, &e.SeqNo)
	binary.Read(reader, binary.LittleEndian, &e.CurveType)
	binary.Read(reader, binary.LittleEndian, &e.HashSize)

}

func (e *EcdsaSignReqInfo) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 118))

	binary.Write(buf, binary.LittleEndian, e.REQ)
	binary.Write(buf, binary.LittleEndian, e.Status)
	binary.Write(buf, binary.LittleEndian, e.UserData)
	binary.Write(buf, binary.LittleEndian, e.SeqNo)
	binary.Write(buf, binary.LittleEndian, e.CurveType)
	binary.Write(buf, binary.LittleEndian, e.HashSize)

	return buf.Bytes()
}

type sslSetActiveCsList struct {
	CipherSuiteBitmap uint32
}

func (s *sslSetActiveCsList) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &s.CipherSuiteBitmap)
}

func (s *sslSetActiveCsList) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	binary.Write(buf, binary.LittleEndian, s.CipherSuiteBitmap)

	return buf.Bytes()
}
