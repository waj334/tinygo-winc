package winc

import (
	"bytes"
	"encoding/binary"
)

type wifiConnectionCredentialsHeader struct {
	credentialsSize int16
	storageFlags    byte
	channel         byte
	// 4 Bytes
}

func (w *wifiConnectionCredentialsHeader) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	binary.Write(buf, binary.LittleEndian, w.credentialsSize)
	buf.WriteByte(w.storageFlags)
	buf.WriteByte(w.channel)

	return buf.Bytes()
}

type wifiConnectionCredentialsCommon struct {
	/* ssidLen byte */
	ssid     []byte // 32 bytes
	options  byte
	bssid    []byte // 6 bytes
	authType byte
	/* reserved [3]byte */
	// 44 Bytes
}

func (w *wifiConnectionCredentialsCommon) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 44))

	if len(w.ssid) > 32 {
		w.ssid = w.ssid[:32]
	}

	if len(w.bssid) > 6 {
		w.bssid = w.bssid[:6]
	}

	buf.WriteByte(byte(len(w.ssid)))

	buf.Write(w.ssid)
	pad(len(w.ssid), 32, buf)

	buf.WriteByte(w.options)

	buf.Write(w.bssid)
	pad(len(w.bssid), 6, buf)

	buf.WriteByte(w.authType)

	buf.WriteByte(0) // Reserved padding
	buf.WriteByte(0) // Reserved padding
	buf.WriteByte(0) // Reserved padding

	return buf.Bytes()
}

type wifiConnectionHeader struct {
	hdr wifiConnectionCredentialsHeader
	cmn wifiConnectionCredentialsCommon
	// 48 bytes
}

func (w *wifiConnectionHeader) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 48))
	buf.Write(w.hdr.bytes())
	buf.Write(w.cmn.bytes())

	return buf.Bytes()
}

type wifiPsk struct {
	/* passphraseLen uint8 */
	passphrase []byte // 64 bytes
	psk        []byte // 40 bytes
	pskValue   byte
	/* reserved [2]byte */
	// 108 bytes
}

func (w *wifiPsk) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 108))
	if len(w.passphrase) > 64 {
		w.passphrase = w.passphrase[:64]
	}

	if len(w.psk) > 40 {
		w.psk = w.psk[:40]
	}

	buf.WriteByte(byte(len(w.passphrase)))

	buf.Write(w.passphrase)
	pad(len(w.passphrase), 64, buf)

	buf.Write(w.psk)
	pad(len(w.psk), 40, buf)

	buf.WriteByte(w.pskValue)
	buf.WriteByte(0) // Reserved padding
	buf.WriteByte(0) // Reserved padding

	return buf.Bytes()
}

type SystemTime struct {
	Year   uint16
	Month  byte
	Day    byte
	Hour   byte
	Minute byte
	Second byte
	/* padding byte */
	// 8 bytes
}

func (s *SystemTime) read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, &s.Year)
	s.Month, _ = reader.ReadByte()
	s.Day, _ = reader.ReadByte()
	s.Hour, _ = reader.ReadByte()
	s.Minute, _ = reader.ReadByte()
	s.Second, _ = reader.ReadByte()
}

type IpConfig struct {
	StaticIP      uint32
	Gateway       uint32
	DNS           uint32
	AlternateDNS  uint32
	SubnetMask    uint32
	DhcpLeaseTime uint32
}

func (i *IpConfig) read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, &i.StaticIP)
	binary.Read(reader, binary.LittleEndian, &i.Gateway)
	binary.Read(reader, binary.LittleEndian, &i.DNS)
	binary.Read(reader, binary.LittleEndian, &i.AlternateDNS)
	binary.Read(reader, binary.LittleEndian, &i.SubnetMask)
	binary.Read(reader, binary.LittleEndian, &i.DhcpLeaseTime)
}

type WifiStateChanged struct {
	CurrentState WifiState
	ErrorCode    byte
	/* padding [2]byte */
}

func (w *WifiStateChanged) read(data []byte) {
	reader := bytes.NewReader(data)

	state, _ := reader.ReadByte()
	w.CurrentState = WifiState(state)

	w.ErrorCode, _ = reader.ReadByte()
}

type WifiConnectionInfo struct {
	SSID           string // 33 bytes
	SecurityType   WifiSecurityType
	IPAddress      [4]byte
	MACAddress     [6]byte
	RSSI           byte
	CurrentChannel WifiChannel
	/* padding [2]byte */
	// 48 bytes
}

func (w *WifiConnectionInfo) read(data []byte) {
	reader := bytes.NewBuffer(data)
	w.SSID, _ = reader.ReadString(0)
	reader.Next(33 - len(w.SSID))

	secType, _ := reader.ReadByte()
	w.SecurityType = WifiSecurityType(secType)

	reader.Read(w.IPAddress[:])
	reader.Read(w.MACAddress[:])
	w.RSSI, _ = reader.ReadByte()

	channel, _ := reader.ReadByte()
	w.CurrentChannel = WifiChannel(channel)
}
