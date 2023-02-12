/*
MIT License

Copyright (c) 2022 waj334

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package winc

import (
	"encoding/binary"
	"time"

	"github.com/waj334/tinygo-winc/debug"
	"github.com/waj334/tinygo-winc/protocol"
)

// Wifi opcodes

const (
	OpcodeWifiReqGetConnInfo = 5
	OpcodeWifiRespConnInfo   = 6
	OpcodeWifiRespGetSysTime = 27
)

const (
	OpcodeWifiReqConnect protocol.OpcodeId = iota + 40
	OpcodeWifiReqDefaultConnect
	OpcodeWifiRespDefaultConnect
	OpcodeWifiReqDisconnect
	OpcodeWifiRespConStateChanged
	OpcodeWifiReqSleep
	OpcodeWifiReqWpsScan
	OpcodeWifiReqWps
	OpcodeWifiReqStartWps
	OpcodeWifiReqDisableWps
	OpcodeWifiReqDhcpConf
	OpcodeWifiRespIpConfigured
	OpcodeWifiRespIpConflict
	OpcodeWifiReqEnableMonitoring
	OpcodeWifiReqDisableMonitoring
	OpcodeWifiRespWifiRxPacket
	OpcodeWifiReqSendWifiPacket
	OpcodeWifiReqLsnInt
	OpcodeWifiReqDoze
	OpcodeWifiReqConn
	OpcodeWifiIndConnParam
	OpcodeWifiReqDhcpFailure
)

const (
	OpcodeWifiReqEnableAP protocol.OpcodeId = iota + 70
	OpcodeWifiReqDisableAP
)

type WifiCredOption uint8

const (
	WifiCredDontSave        WifiCredOption = 0
	WifiCredSaveUnencrypted WifiCredOption = 0x01
	WifiCredSaveEncrypted   WifiCredOption = 0x02
)

type WifiChannel uint8

const (
	WifiChannel1 WifiChannel = iota + 1
	WifiChannel2
	WifiChannel3
	WifiChannel4
	WifiChannel5
	WifiChannel6
	WifiChannel7
	WifiChannel8
	WifiChannel9
	WifiChannel10
	WifiChannel11
	WifiChannel12
	WifiChannel13
	WifiChannel14
	WifiChannelAll WifiChannel = 255
)

type WifiSecurityType uint8

const (
	WifiSecurityOpen WifiSecurityType = iota + 1
	WifiSecurityWpaPsk
	WifiSecurityWep
	WifiSecurity8021X
)

type WifiState uint8

const (
	WifiStateDisconnected WifiState = iota
	WifiStateConnected
	WifiStateRoamed
	WifiStateUnknown
)

type WifiConnectionSettings struct {
	Bssid      [6]byte
	Ssid       string
	Channel    WifiChannel
	Passphrase string
	Storage    WifiCredOption
	Security   WifiSecurityType
}

func (w *WINC) WifiConnectPsk(settings WifiConnectionSettings) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	debug.DEBUG("WIFI: WifiConnectPsk - BEGIN")
	defer debug.DEBUG("WIFI: WifiConnectPsk - END")

	opcode := OpcodeWifiReqConn
	var data []byte

	if settings.Security == 0 {
		// Default to WPA
		settings.Security = WifiSecurityWpaPsk
	}

	if len(settings.Ssid) == 0 {
		return ErrInvalidParameter
	}

	if settings.Security != WifiSecurityOpen {
		if len(settings.Passphrase) == 0 || len(settings.Passphrase) >= 64 {
			return ErrInvalidParameter
		}

		opcode |= protocol.OpcodeReqDataPkt

		// TODO: Handle PSK - 2

		credentials := wifiPsk{
			passphrase: []byte(settings.Passphrase),
			psk:        nil,
			pskValue:   0,
		}

		data = credentials.bytes()
	}

	control := wifiConnectionHeader{
		hdr: wifiConnectionCredentialsHeader{
			credentialsSize: 152,
			storageFlags:    byte(settings.Storage),
			channel:         byte(settings.Channel),
		},
		cmn: wifiConnectionCredentialsCommon{
			ssid:     []byte(settings.Ssid),
			options:  0, //TODO
			bssid:    settings.Bssid[:],
			authType: byte(settings.Security),
		},
	}

	// Send the HIF command
	if err = w.hif.Send(GroupWIFI, opcode, control.bytes(), data, 48); err != nil {
		return
	}

	return nil
}

func (w *WINC) WifiDisconnect() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.hif.Send(GroupWIFI, OpcodeWifiReqDisconnect, nil, nil, 0)
}

func (w *WINC) GetWifiState() WifiState {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.wifiState
}

func (w *WINC) GetConnectionInfo() (strConnInfo *WifiConnectionInfo, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err = w.hif.Send(GroupWIFI, OpcodeWifiReqGetConnInfo, nil, nil, 0); err != nil {
		return
	}

	// Wait for reply
	select {
	case reply := <-w.callbackChan:
		strConnInfo = reply.(*WifiConnectionInfo)
	}

	return
}

func (w *WINC) EnableAP(config APModeConfig) (err error) {
	// Validate configuration
	if config.SecurityType == WifiSecurityWep {
		return ErrInvalidParameter
	} else if config.SecurityType == WifiSecurityWpaPsk && len(config.WPAKey) == 0 {
		return ErrInvalidParameter
	} else if config.Channel < WifiChannel1 || config.Channel > WifiChannel14 {
		return ErrInvalidParameter
	} // TODO check IP

	if err = w.hif.Send(GroupWIFI, OpcodeWifiReqEnableAP|protocol.OpcodeReqDataPkt,
		nil, config.bytes(), 0); err != nil {

		return
	}

	return
}

func (w *WINC) DisableAP() (err error) {
	if err = w.hif.Send(GroupWIFI, OpcodeWifiReqDisableAP, nil, nil, 0); err != nil {
		return
	}

	return
}

func (w *WINC) wifiCallback(id protocol.OpcodeId, sz uint16, address uint32) (obj any, err error) {
	switch id {
	case OpcodeWifiRespGetSysTime:
		data := make([]byte, 8)
		if err = w.hif.Receive(address, data, false); err != nil {
			return
		}

		strSysTime := &SystemTime{}
		strSysTime.read(data)

		obj = strSysTime
	case OpcodeWifiReqDhcpConf:
		data := make([]byte, 24)
		if err = w.hif.Receive(address, data, false); err != nil {
			return
		}

		strIpConfig := &IpConfig{}
		strIpConfig.read(data)

		w.ipAddr.IP = make([]byte, 4)
		w.ipAddr.Mask = make([]byte, 4)
		binary.BigEndian.PutUint32(w.ipAddr.IP, strIpConfig.StaticIP)
		binary.BigEndian.PutUint32(w.ipAddr.Mask, strIpConfig.SubnetMask)

		obj = strIpConfig
	case OpcodeWifiRespConStateChanged:
		data := make([]byte, 4)
		if err = w.hif.Receive(address, data, false); err != nil {
			return
		}

		strState := &WifiStateChanged{}
		strState.read(data)

		w.wifiState = WifiState(strState.CurrentState)
		obj = strState
	case OpcodeWifiRespConnInfo:
		data := make([]byte, 48)
		if err = w.hif.Receive(address, data, false); err != nil {
			return
		}

		strConnInfo := &WifiConnectionInfo{}
		strConnInfo.read(data)

		w.callbackChan <- strConnInfo
		obj = strConnInfo
	}

	return
}

func SysTimeToDate(strSysTime *SystemTime) time.Time {
	return time.Date(
		int(strSysTime.Year),
		time.Month(strSysTime.Month),
		int(strSysTime.Day),
		int(strSysTime.Hour),
		int(strSysTime.Minute),
		int(strSysTime.Second),
		0,
		time.UTC,
	)
}
