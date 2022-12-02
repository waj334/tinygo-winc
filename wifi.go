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

	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/types"
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

type WifiCredOption uint8

const (
	WifiCredDontSave        WifiCredOption = 0
	WifiCredSaveUnencrypted WifiCredOption = 0x01
	WifiCredSaveEncrypted   WifiCredOption = 0x02
)

type WifiChannel uint8

const (
	WifiChannel1 WifiChannel = iota
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
	security   WifiSecurityType
}

func (w *WINC) WifiConnectPsk(settings WifiConnectionSettings) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	settings.security = WifiSecurityWpaPsk
	if len(settings.Ssid) == 0 {
		return ErrInvalidParameter
	}

	if len(settings.Passphrase) == 0 || len(settings.Passphrase) > 64 {
		return ErrInvalidParameter
	}

	// TODO: Handle PSK - 2

	var ctrl types.M2mWifiConnHdr
	ctrl.StrConnCredHdr.U8CredStoreFlags = byte(settings.Storage)
	ctrl.StrConnCredHdr.U8Channel = byte(settings.Channel)
	ctrl.StrConnCredHdr.U16CredSize = 44

	copy(ctrl.StrConnCredCmn.Au8Ssid[:], settings.Ssid)
	ctrl.StrConnCredCmn.U8SsidLen = uint8(len(settings.Ssid))
	ctrl.StrConnCredCmn.U8AuthType = uint8(settings.security)
	ctrl.StrConnCredCmn.U8Options = uint8(0) // TODO
	ctrl.StrConnCredCmn.Au8Bssid = settings.Bssid

	psk := types.M2mWifiPsk{
		U8PassphraseLen: uint8(len(settings.Passphrase)),
	}

	copy(psk.Au8Passphrase[:], settings.Passphrase)

	// Send the HIF command
	if err = w.hif.Send(GroupWIFI,
		OpcodeWifiReqConn|protocol.OpcodeReqDataPkt,
		ctrl.Bytes(),
		psk.Bytes(),
		uint16(len(ctrl.Bytes())),
	); err != nil {
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

func (w *WINC) GetConnectionInfo() (strConnInfo types.M2MConnInfo, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.pendingCallback = true
	if err = w.hif.Send(GroupWIFI, OpcodeWifiReqGetConnInfo, nil, nil, 0); err != nil {
		return
	}

	// Wait for reply
	select {
	case reply := <-w.callbackChan:
		strConnInfo = reply.(types.M2MConnInfo)
		w.pendingCallback = false
	}

	return
}

func (w *WINC) wifiCallback(id protocol.OpcodeId, sz uint16, address uint32) (data any, err error) {
	switch id {
	case OpcodeWifiRespGetSysTime:
		var strSysTime types.SystemTime
		if err = w.hif.Receive(address, strSysTime.Bytes(), false); err != nil {
			return
		}

		strSysTime.Deref()
		strSysTime.Free()

		data = strSysTime
	case OpcodeWifiReqDhcpConf:
		var strIpConfig types.M2MIPConfig
		if err = w.hif.Receive(address, strIpConfig.Bytes(), false); err != nil {
			return
		}

		strIpConfig.Deref()
		strIpConfig.Free()

		w.ipAddr.IP = make([]byte, 4)
		w.ipAddr.Mask = make([]byte, 4)
		binary.BigEndian.PutUint32(w.ipAddr.IP, strIpConfig.U32StaticIP)
		binary.BigEndian.PutUint32(w.ipAddr.Mask, strIpConfig.U32SubnetMask)

		data = strIpConfig
	case OpcodeWifiRespConStateChanged:
		var strState types.M2mWifiStateChanged
		if err = w.hif.Receive(address, strState.Bytes(), false); err != nil {
			return
		}

		strState.Deref()
		strState.Free()

		w.wifiState = WifiState(strState.U8CurrState)
		data = strState
	case OpcodeWifiRespConnInfo:
		var strConnInfo types.M2MConnInfo
		if err = w.hif.Receive(address, strConnInfo.Bytes(), false); err != nil {
			return
		}

		strConnInfo.Deref()
		strConnInfo.Free()

		w.callbackChan <- strConnInfo
		data = strConnInfo
	}

	return
}

func SysTimeToDate(strSysTime types.SystemTime) time.Time {
	return time.Date(
		int(strSysTime.U16Year),
		time.Month(strSysTime.U8Month),
		int(strSysTime.U8Day),
		int(strSysTime.U8Hour),
		int(strSysTime.U8Minute),
		int(strSysTime.U8Second),
		0,
		time.UTC,
	)
}
