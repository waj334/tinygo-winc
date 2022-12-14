// WARNING: This file has automatically been generated
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package types

/*
#include "../include/m2m_types.h"
#include "../include/socket.h"
#include "../include/m2m_socket_host_if.h"
#include "../include/m2m_hif.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

// HifHdr as declared in include/m2m_hif.h:63
type HifHdr struct {
	U8Gid          byte
	U8Opcode       byte
	U16Length      uint16
	ref2e9cbbc5    *C.tstrHifHdr
	allocs2e9cbbc5 interface{}
}

// SockAddr as declared in include/m2m_socket_host_if.h:226
type SockAddr struct {
	U16Family     uint16
	U16Port       uint16
	U32IPAddr     uint32
	reff979731    *C.tstrSockAddr
	allocsf979731 interface{}
}

// UIPSockAddr as declared in include/m2m_socket_host_if.h:230
type UIPSockAddr struct {
	U16Family      uint16
	U16Port        uint16
	U32IPAddr      uint32
	ref73ec2e23    *C.tstrUIPSockAddr
	allocs73ec2e23 interface{}
}

// DnsReply as declared in include/m2m_socket_host_if.h:244
type DnsReply struct {
	AcHostName     [64]byte
	U32HostIP      uint32
	ref1b0d1c9e    *C.tstrDnsReply
	allocs1b0d1c9e interface{}
}

// BindCmd as declared in include/m2m_socket_host_if.h:255
type BindCmd struct {
	StrAddr        SockAddr
	Sock           int8
	U8Void         byte
	U16SessionID   uint16
	reffd41baba    *C.tstrBindCmd
	allocsfd41baba interface{}
}

// BindReply as declared in include/m2m_socket_host_if.h:265
type BindReply struct {
	Sock           int8
	S8Status       int8
	U16SessionID   uint16
	ref17cd835c    *C.tstrBindReply
	allocs17cd835c interface{}
}

// ListenCmd as declared in include/m2m_socket_host_if.h:275
type ListenCmd struct {
	Sock           int8
	U8BackLog      byte
	U16SessionID   uint16
	ref42e0429d    *C.tstrListenCmd
	allocs42e0429d interface{}
}

// ListenReply as declared in include/m2m_socket_host_if.h:293
type ListenReply struct {
	Sock           int8
	S8Status       int8
	U16SessionID   uint16
	ref7e6ef4b6    *C.tstrListenReply
	allocs7e6ef4b6 interface{}
}

// AcceptReply as declared in include/m2m_socket_host_if.h:308
type AcceptReply struct {
	StrAddr          SockAddr
	SListenSock      int8
	SConnectedSock   int8
	U16AppDataOffset uint16
	ref29a2c828      *C.tstrAcceptReply
	allocs29a2c828   interface{}
}

// ConnectCmd as declared in include/m2m_socket_host_if.h:319
type ConnectCmd struct {
	StrAddr        SockAddr
	Sock           int8
	U8SslFlags     byte
	U16SessionID   uint16
	refc65d5400    *C.tstrConnectCmd
	allocsc65d5400 interface{}
}

// ConnectReply as declared in include/m2m_socket_host_if.h:358
type ConnectReply struct {
	Sock           int8
	S8Error        int8
	U16ExtraData   uint16
	ref4e55310d    *C.tstrConnectReply
	allocs4e55310d interface{}
}

// CloseCmd as declared in include/m2m_socket_host_if.h:364
type CloseCmd struct {
	Sock           int8
	U8Dummy        byte
	U16SessionID   uint16
	ref52d2c6ab    *C.tstrCloseCmd
	allocs52d2c6ab interface{}
}

// ConnectAlpnReply as declared in include/m2m_socket_host_if.h:380
type ConnectAlpnReply struct {
	StrConnReply     ConnectReply
	U8AppProtocolIdx byte
	__PAD24__        [3]byte
	ref696c1533      *C.tstrConnectAlpnReply
	allocs696c1533   interface{}
}

// SendCmd as declared in include/m2m_socket_host_if.h:393
type SendCmd struct {
	Sock           int8
	U8Void         byte
	U16DataSize    uint16
	StrAddr        SockAddr
	U16SessionID   uint16
	U16Void        uint16
	ref4b2aedbb    *C.tstrSendCmd
	allocs4b2aedbb interface{}
}

// SendReply as declared in include/m2m_socket_host_if.h:409
type SendReply struct {
	Sock           int8
	U8Void         byte
	S16SentBytes   int16
	U16SessionID   uint16
	U16Void        uint16
	reffbd9c021    *C.tstrSendReply
	allocsfbd9c021 interface{}
}

// RecvCmd as declared in include/m2m_socket_host_if.h:421
type RecvCmd struct {
	U32Timeoutmsec uint32
	Sock           int8
	U8Void         byte
	U16SessionID   uint16
	U16BufLen      uint16
	refefddbdaa    *C.tstrRecvCmd
	allocsefddbdaa interface{}
}

// RecvReply as declared in include/m2m_socket_host_if.h:436
type RecvReply struct {
	StrRemoteAddr  SockAddr
	S16RecvStatus  int16
	U16DataOffset  uint16
	Sock           int8
	U8Void         byte
	U16SessionID   uint16
	refc36b1632    *C.tstrRecvReply
	allocsc36b1632 interface{}
}

// SetSocketOptCmd as declared in include/m2m_socket_host_if.h:447
type SetSocketOptCmd struct {
	U32OptionValue uint32
	Sock           int8
	U8Option       byte
	U16SessionID   uint16
	ref2d3fd62a    *C.tstrSetSocketOptCmd
	allocs2d3fd62a interface{}
}

// SSLSocketCreateCmd as declared in include/m2m_socket_host_if.h:453
type SSLSocketCreateCmd struct {
	SslSock        int8
	__PAD24__      [3]byte
	refa5924504    *C.tstrSSLSocketCreateCmd
	allocsa5924504 interface{}
}

// SSLSetSockOptCmd as declared in include/m2m_socket_host_if.h:465
type SSLSetSockOptCmd struct {
	Sock           int8
	U8Option       byte
	U16SessionID   uint16
	U32OptLen      uint32
	Au8OptVal      [64]byte
	ref75a4c4a8    *C.tstrSSLSetSockOptCmd
	allocs75a4c4a8 interface{}
}

// PingCmd as declared in include/m2m_socket_host_if.h:476
type PingCmd struct {
	U32DestIPAddr  uint32
	U32CmdPrivate  uint32
	U16PingCount   uint16
	U8TTL          byte
	__PAD8__       byte
	ref1fb598b3    *C.tstrPingCmd
	allocs1fb598b3 interface{}
}

// PingReply as declared in include/m2m_socket_host_if.h:487
type PingReply struct {
	U32IPAddr      uint32
	U32CmdPrivate  uint32
	U32RTT         uint32
	U16Success     uint16
	U16Fail        uint16
	U8ErrorCode    byte
	__PAD24__      [3]byte
	ref136f9b05    *C.tstrPingReply
	allocs136f9b05 interface{}
}

// SslCertExpSettings as declared in include/m2m_socket_host_if.h:502
type SslCertExpSettings struct {
	U32CertExpValidationOpt uint32
	reff3e2216d             *C.tstrSslCertExpSettings
	allocsf3e2216d          interface{}
}

// M2mPwrMode as declared in include/m2m_types.h:614
type M2mPwrMode struct {
	U8PwrMode      byte
	__PAD24__      [3]byte
	ref247585ac    *C.tstrM2mPwrMode
	allocs247585ac interface{}
}

// M2mTxPwrLevel as declared in include/m2m_types.h:645
type M2mTxPwrLevel struct {
	U8TxPwrLevel   byte
	__PAD24__      [3]byte
	ref6aa1b7e0    *C.tstrM2mTxPwrLevel
	allocs6aa1b7e0 interface{}
}

// M2mWiFiGainIdx as declared in include/m2m_types.h:661
type M2mWiFiGainIdx struct {
	U8GainTableIdx byte
	__PAD24__      [3]byte
	ref2363fa99    *C.tstrM2mWiFiGainIdx
	allocs2363fa99 interface{}
}

// M2mEnableLogs as declared in include/m2m_types.h:677
type M2mEnableLogs struct {
	U8Enable       byte
	__PAD24__      [3]byte
	ref8568bfe5    *C.tstrM2mEnableLogs
	allocs8568bfe5 interface{}
}

// M2mBatteryVoltage as declared in include/m2m_types.h:694
type M2mBatteryVoltage struct {
	U16BattVolt    uint16
	__PAD16__      [2]byte
	refdef96f15    *C.tstrM2mBatteryVoltage
	allocsdef96f15 interface{}
}

// M2mWiFiRoaming as declared in include/m2m_types.h:713
type M2mWiFiRoaming struct {
	U8EnableRoaming byte
	U8EnableDhcp    byte
	__PAD16__       [2]byte
	ref41e06f30     *C.tstrM2mWiFiRoaming
	allocs41e06f30  interface{}
}

// M2mWiFiXOSleepEnable as declared in include/m2m_types.h:729
type M2mWiFiXOSleepEnable struct {
	U8EnableXODuringSleep byte
	__PAD16__             [3]byte
	refe8bee5c5           *C.tstrM2mWiFiXOSleepEnable
	allocse8bee5c5        interface{}
}

// M2mWifiGainsParams as declared in include/m2m_types.h:1345
type M2mWifiGainsParams struct {
	U8PPAGFor11B   uint16
	U8PPAGFor11GN  uint16
	ref69e233ef    *C.tstrM2mWifiGainsParams
	allocs69e233ef interface{}
}

// M2mConnCredHdr as declared in include/m2m_types.h:1369
type M2mConnCredHdr struct {
	U16CredSize      uint16
	U8CredStoreFlags byte
	U8Channel        byte
	ref39e9ce5f      *C.tstrM2mConnCredHdr
	allocs39e9ce5f   interface{}
}

// M2mConnCredCmn as declared in include/m2m_types.h:1393
type M2mConnCredCmn struct {
	U8SsidLen      byte
	Au8Ssid        [32]byte
	U8Options      byte
	Au8Bssid       [6]byte
	U8AuthType     byte
	Au8Rsv         [3]byte
	reff07fc6b8    *C.tstrM2mConnCredCmn
	allocsf07fc6b8 interface{}
}

// M2mWifiWep as declared in include/m2m_types.h:1411
type M2mWifiWep struct {
	U8KeyIndex     byte
	U8KeyLen       byte
	Au8WepKey      [13]byte
	U8Rsv          byte
	ref3c80b060    *C.tstrM2mWifiWep
	allocs3c80b060 interface{}
}

// M2mWifiPsk as declared in include/m2m_types.h:1431
type M2mWifiPsk struct {
	U8PassphraseLen byte
	Au8Passphrase   [64]byte
	Au8Psk          [40]byte
	U8PskCalculated byte
	Au8Rsv          [2]byte
	refaf32dade     *C.tstrM2mWifiPsk
	allocsaf32dade  interface{}
}

// M2mWifi1xHdr as declared in include/m2m_types.h:1476
type M2mWifi1xHdr struct {
	U8Flags                    byte
	U8DomainLength             byte
	U8UserNameLength           byte
	U8HdrLength                byte
	U16PrivateKeyOffset        uint16
	U16PrivateKeyLength        uint16
	U16CertificateOffset       uint16
	U16CertificateLength       uint16
	Au8TlsSpecificRootNameSha1 [20]byte
	U32Rsv1                    uint32
	U32TlsHsFlags              uint32
	U32Rsv2                    uint32
	refcea6f160                *C.tstrM2mWifi1xHdr
	allocscea6f160             interface{}
}

// M2mWifiAuthInfoHdr as declared in include/m2m_types.h:1499
type M2mWifiAuthInfoHdr struct {
	U8Type         byte
	Au8Rsv         [3]byte
	U16InfoPos     uint16
	U16InfoLen     uint16
	ref879b4dbc    *C.tstrM2mWifiAuthInfoHdr
	allocs879b4dbc interface{}
}

// M2mWifiConnHdr as declared in include/m2m_types.h:1520
type M2mWifiConnHdr struct {
	StrConnCredHdr M2mConnCredHdr
	StrConnCredCmn M2mConnCredCmn
	ref3ce06260    *C.tstrM2mWifiConnHdr
	allocs3ce06260 interface{}
}

// M2mWifiApId as declared in include/m2m_types.h:1538
type M2mWifiApId struct {
	Au8SSID        [33]byte
	__PAD__        [3]byte
	refb4484f46    *C.tstrM2mWifiApId
	allocsb4484f46 interface{}
}

// M2MGenericResp as declared in include/m2m_types.h:1555
type M2MGenericResp struct {
	S8ErrorCode    int8
	__PAD24__      [3]byte
	ref3d54d4d0    *C.tstrM2MGenericResp
	allocs3d54d4d0 interface{}
}

// M2MWPSConnect as declared in include/m2m_types.h:1577
type M2MWPSConnect struct {
	U8TriggerType  byte
	AcPinNumber    [8]byte
	__PAD24__      [3]byte
	ref649eedf4    *C.tstrM2MWPSConnect
	allocs649eedf4 interface{}
}

// M2MWPSInfo as declared in include/m2m_types.h:1605
type M2MWPSInfo struct {
	U8AuthType     byte
	U8Ch           byte
	Au8SSID        [33]byte
	Au8PSK         [65]byte
	refa2ef9cb2    *C.tstrM2MWPSInfo
	allocsa2ef9cb2 interface{}
}

// M2MDefaultConnResp as declared in include/m2m_types.h:1627
type M2MDefaultConnResp struct {
	S8ErrorCode    int8
	__PAD24__      [3]byte
	ref4b129f8a    *C.tstrM2MDefaultConnResp
	allocs4b129f8a interface{}
}

// M2MScanOption as declared in include/m2m_types.h:1668
type M2MScanOption struct {
	U8NumOfSlot     byte
	U8SlotTime      byte
	U8ProbesPerSlot byte
	S8RssiThresh    int8
	refac5e69f9     *C.tstrM2MScanOption
	allocsac5e69f9  interface{}
}

// M2MStopScanOption as declared in include/m2m_types.h:1690
type M2MStopScanOption struct {
	U8StopOnFirstResult byte
	Au8Rsv              [3]byte
	refd0dabb43         *C.tstrM2MStopScanOption
	allocsd0dabb43      interface{}
}

// M2MScanRegion as declared in include/m2m_types.h:1708
type M2MScanRegion struct {
	U16ScanRegion  uint16
	__PAD16__      [2]byte
	reff9ba983f    *C.tstrM2MScanRegion
	allocsf9ba983f interface{}
}

// M2MScan as declared in include/m2m_types.h:1734
type M2MScan struct {
	U8ChNum            byte
	__RSVD8__          [1]byte
	U16PassiveScanTime uint16
	reff2f04759        *C.tstrM2MScan
	allocsf2f04759     interface{}
}

// CyptoResp as declared in include/m2m_types.h:1749
type CyptoResp struct {
	S8Resp         int8
	__PAD24__      [3]byte
	refd21c28b2    *C.tstrCyptoResp
	allocsd21c28b2 interface{}
}

// M2mScanDone as declared in include/m2m_types.h:1768
type M2mScanDone struct {
	U8NumofCh      byte
	S8ScanState    int8
	__PAD16__      [2]byte
	refb63164d6    *C.tstrM2mScanDone
	allocsb63164d6 interface{}
}

// M2mReqScanResult as declared in include/m2m_types.h:1785
type M2mReqScanResult struct {
	U8Index        byte
	__PAD24__      [3]byte
	refec922c38    *C.tstrM2mReqScanResult
	allocsec922c38 interface{}
}

// M2mWifiscanResult as declared in include/m2m_types.h:1817
type M2mWifiscanResult struct {
	U8index        byte
	S8rssi         int8
	U8AuthType     byte
	U8ch           byte
	Au8BSSID       [6]byte
	Au8SSID        [33]byte
	_PAD8_         byte
	reffc3bcece    *C.tstrM2mWifiscanResult
	allocsfc3bcece interface{}
}

// M2mWifiStateChanged as declared in include/m2m_types.h:1839
type M2mWifiStateChanged struct {
	U8CurrState    byte
	U8ErrCode      byte
	__PAD16__      [2]byte
	reff7464335    *C.tstrM2mWifiStateChanged
	allocsf7464335 interface{}
}

// M2mPsType as declared in include/m2m_types.h:1861
type M2mPsType struct {
	U8PsType       byte
	U8BcastEn      byte
	__PAD16__      [2]byte
	ref27aa5fd6    *C.tstrM2mPsType
	allocs27aa5fd6 interface{}
}

// M2mSlpReqTime as declared in include/m2m_types.h:1876
type M2mSlpReqTime struct {
	U32SleepTime   uint32
	refb428ccb8    *C.tstrM2mSlpReqTime
	allocsb428ccb8 interface{}
}

// M2mLsnInt as declared in include/m2m_types.h:1894
type M2mLsnInt struct {
	U16LsnInt      uint16
	__PAD16__      [2]byte
	ref56b27f4e    *C.tstrM2mLsnInt
	allocs56b27f4e interface{}
}

// M2MWifiMonitorModeCtrl as declared in include/m2m_types.h:1931
type M2MWifiMonitorModeCtrl struct {
	U8ChannelID      byte
	U8FrameType      byte
	U8FrameSubtype   byte
	Au8SrcMacAddress [6]byte
	Au8DstMacAddress [6]byte
	Au8BSSID         [6]byte
	U8EnRecvHdr      byte
	__PAD16__        [2]byte
	refc77b0df       *C.tstrM2MWifiMonitorModeCtrl
	allocsc77b0df    interface{}
}

// M2MWifiRxPacketInfo as declared in include/m2m_types.h:1985
type M2MWifiRxPacketInfo struct {
	U8FrameType      byte
	U8FrameSubtype   byte
	U8ServiceClass   byte
	U8Priority       byte
	U8HeaderLength   byte
	U8CipherType     byte
	Au8SrcMacAddress [6]byte
	Au8DstMacAddress [6]byte
	Au8BSSID         [6]byte
	U16DataLength    uint16
	U16FrameLength   uint16
	U32DataRateKbps  uint32
	S8RSSI           int8
	__PAD24__        [3]byte
	ref8954124c      *C.tstrM2MWifiRxPacketInfo
	allocs8954124c   interface{}
}

// M2MWifiTxPacketInfo as declared in include/m2m_types.h:2003
type M2MWifiTxPacketInfo struct {
	U16PacketSize   uint16
	U16HeaderLength uint16
	ref6669b7ac     *C.tstrM2MWifiTxPacketInfo
	allocs6669b7ac  interface{}
}

// M2MAPConfig as declared in include/m2m_types.h:2048
type M2MAPConfig struct {
	Au8SSID         [33]byte
	U8ListenChannel byte
	U8KeyIndx       byte
	U8KeySz         byte
	Au8WepKey       [27]byte
	U8SecType       byte
	U8SsidHide      byte
	Au8DHCPServerIP [4]byte
	Au8Key          [65]byte
	__PAD24__       [2]byte
	ref85c0615f     *C.tstrM2MAPConfig
	allocs85c0615f  interface{}
}

// M2MAPConfigExt as declared in include/m2m_types.h:2072
type M2MAPConfigExt struct {
	Au8DefRouterIP [4]byte
	Au8DNSServerIP [4]byte
	Au8SubnetMask  [4]byte
	refb9fc21bc    *C.tstrM2MAPConfigExt
	allocsb9fc21bc interface{}
}

// M2MAPModeConfig as declared in include/m2m_types.h:2093
type M2MAPModeConfig struct {
	StrApConfig    M2MAPConfig
	StrApConfigExt M2MAPConfigExt
	ref75fb89de    *C.tstrM2MAPModeConfig
	allocs75fb89de interface{}
}

// M2mServerInit as declared in include/m2m_types.h:2109
type M2mServerInit struct {
	U8Channel      byte
	__PAD24__      [3]byte
	ref8f1c66f0    *C.tstrM2mServerInit
	allocs8f1c66f0 interface{}
}

// M2mClientState as declared in include/m2m_types.h:2125
type M2mClientState struct {
	U8State        byte
	__PAD24__      [3]byte
	ref13134fc2    *C.tstrM2mClientState
	allocs13134fc2 interface{}
}

// M2Mservercmd as declared in include/m2m_types.h:2141
type M2Mservercmd struct {
	U8cmd          byte
	__PAD24__      [3]byte
	refdb7812b5    *C.tstrM2Mservercmd
	allocsdb7812b5 interface{}
}

// M2mSetMacAddress as declared in include/m2m_types.h:2161
type M2mSetMacAddress struct {
	Au8Mac         [6]byte
	__PAD16__      [2]byte
	ref34f51c02    *C.tstrM2mSetMacAddress
	allocs34f51c02 interface{}
}

// M2MDeviceNameConfig as declared in include/m2m_types.h:2176
type M2MDeviceNameConfig struct {
	Au8DeviceName  [48]byte
	ref8fb03441    *C.tstrM2MDeviceNameConfig
	allocs8fb03441 interface{}
}

// M2MIPConfig as declared in include/m2m_types.h:2208
type M2MIPConfig struct {
	U32StaticIP      uint32
	U32Gateway       uint32
	U32DNS           uint32
	U32AlternateDNS  uint32
	U32SubnetMask    uint32
	U32DhcpLeaseTime uint32
	ref56257aea      *C.tstrM2MIPConfig
	allocs56257aea   interface{}
}

// M2mIpRsvdPkt as declared in include/m2m_types.h:2221
type M2mIpRsvdPkt struct {
	U16PktSz       uint16
	U16PktOffset   uint16
	refbbe3c573    *C.tstrM2mIpRsvdPkt
	allocsbbe3c573 interface{}
}

// M2MProvisionModeConfig as declared in include/m2m_types.h:2253
type M2MProvisionModeConfig struct {
	StrApConfig            M2MAPConfig
	AcHttpServerDomainName [64]byte
	U8EnableRedirect       byte
	StrApConfigExt         M2MAPConfigExt
	__PAD24__              [3]byte
	ref7bc66fe8            *C.tstrM2MProvisionModeConfig
	allocs7bc66fe8         interface{}
}

// M2MProvisionInfo as declared in include/m2m_types.h:2281
type M2MProvisionInfo struct {
	Au8SSID        [33]byte
	Au8Password    [65]byte
	U8SecType      byte
	U8Status       byte
	ref32921097    *C.tstrM2MProvisionInfo
	allocs32921097 interface{}
}

// M2MConnInfo as declared in include/m2m_types.h:2306
type M2MConnInfo struct {
	AcSSID         [33]byte
	U8SecType      byte
	Au8IPAddr      [4]byte
	Au8MACAddress  [6]byte
	S8RSSI         int8
	U8CurrChannel  byte
	__PAD16__      [2]byte
	refa529de03    *C.tstrM2MConnInfo
	allocsa529de03 interface{}
}

// M2MP2PConnect as declared in include/m2m_types.h:2349
type M2MP2PConnect struct {
	U8ListenChannel byte
	__PAD24__       [3]byte
	ref7a58584c     *C.tstrM2MP2PConnect
	allocs7a58584c  interface{}
}

// M2MSNTPConfig as declared in include/m2m_types.h:2375
type M2MSNTPConfig struct {
	AcNTPServer    [33]byte
	__PAD8__       [2]byte
	ref649da0d7    *C.tstrM2MSNTPConfig
	allocs649da0d7 interface{}
}

// SystemTime as declared in include/m2m_types.h:2392
type SystemTime struct {
	U16Year        uint16
	U8Month        byte
	U8Day          byte
	U8Hour         byte
	U8Minute       byte
	U8Second       byte
	__PAD8__       byte
	ref47e6eb96    *C.tstrSystemTime
	allocs47e6eb96 interface{}
}

// M2MMulticastMac as declared in include/m2m_types.h:2413
type M2MMulticastMac struct {
	Au8macaddress  [6]byte
	U8AddRemove    byte
	__PAD8__       byte
	refda99fd36    *C.tstrM2MMulticastMac
	allocsda99fd36 interface{}
}

// Prng as declared in include/m2m_types.h:2435
type Prng struct {
	Pu8RngBuff     []byte
	U16PrngSize    uint16
	__PAD16__      [2]byte
	ref37c97783    *C.tstrPrng
	allocs37c97783 interface{}
}

// ConfAutoRate as declared in include/m2m_types.h:2555
type ConfAutoRate struct {
	U16ArMaxRecoveryFailThreshold uint16
	U16ArMinRecoveryFailThreshold uint16
	U8ArEnoughTxThreshold         byte
	U8ArSuccessTXThreshold        byte
	U8ArFailTxThreshold           byte
	__PAD24__                     [3]byte
	ref330ce733                   *C.tstrConfAutoRate
	allocs330ce733                interface{}
}

// TlsCrlEntry as declared in include/m2m_types.h:2599
type TlsCrlEntry struct {
	U8DataLen      byte
	Au8Data        [64]byte
	__PAD24__      [3]byte
	ref3a4218c3    *C.tstrTlsCrlEntry
	allocs3a4218c3 interface{}
}

// TlsCrlInfo as declared in include/m2m_types.h:2619
type TlsCrlInfo struct {
	U8CrlType      byte
	U8Rsv1         byte
	U8Rsv2         byte
	U8Rsv3         byte
	AstrTlsCrl     [10]TlsCrlEntry
	refa907de8c    *C.tstrTlsCrlInfo
	allocsa907de8c interface{}
}

// TlsSrvSecFileEntry as declared in include/m2m_types.h:2663
type TlsSrvSecFileEntry struct {
	AcFileName     [48]byte
	U32FileSize    uint32
	U32FileAddr    uint32
	refca7e9eec    *C.tstrTlsSrvSecFileEntry
	allocsca7e9eec interface{}
}

// TlsSrvSecHdr as declared in include/m2m_types.h:2683
type TlsSrvSecHdr struct {
	Au8SecStartPattern [8]byte
	U32nEntries        uint32
	U32NextWriteAddr   uint32
	AstrEntries        [8]TlsSrvSecFileEntry
	U32CRC             uint32
	refb8311c95        *C.tstrTlsSrvSecHdr
	allocsb8311c95     interface{}
}

// TlsSrvChunkHdr as declared in include/m2m_types.h:2703
type TlsSrvChunkHdr struct {
	U16Sig         uint16
	U16TotalSize32 uint16
	U16Offset32    uint16
	U16Size32      uint16
	ref11892db4    *C.tstrTlsSrvChunkHdr
	allocs11892db4 interface{}
}

// SslSetActiveCsList as declared in include/m2m_types.h:2707
type SslSetActiveCsList struct {
	U32CsBMP       uint32
	refd739fd93    *C.tstrSslSetActiveCsList
	allocsd739fd93 interface{}
}

// OtaInitHdr as declared in include/m2m_types.h:2848
type OtaInitHdr struct {
	U32OtaMagicValue  uint32
	U32OtaPayloadSize uint32
	ref28fa7407       *C.tstrOtaInitHdr
	allocs28fa7407    interface{}
}

// OtaControlSec as declared in include/m2m_types.h:2915
type OtaControlSec struct {
	U32OtaMagicValue                    uint32
	U32OtaFormatVersion                 uint32
	U32OtaSequenceNumber                uint32
	U32OtaLastCheckTime                 uint32
	U32OtaCurrentWorkingImagOffset      uint32
	U32OtaCurrentworkingImagFirmwareVer uint32
	U32OtaRollbackImageOffset           uint32
	U32OtaRollbackImageValidStatus      uint32
	U32OtaRollbackImagFirmwareVer       uint32
	U32OtaCortusAppWorkingOffset        uint32
	U32OtaCortusAppWorkingValidSts      uint32
	U32OtaCortusAppWorkingVer           uint32
	U32OtaCortusAppRollbackOffset       uint32
	U32OtaCortusAppRollbackValidSts     uint32
	U32OtaCortusAppRollbackVer          uint32
	U32OtaControlSecCrc                 uint32
	ref2695ef43                         *C.tstrOtaControlSec
	allocs2695ef43                      interface{}
}

// OtaUpdateStatusResp as declared in include/m2m_types.h:2933
type OtaUpdateStatusResp struct {
	U8OtaUpdateStatusType byte
	U8OtaUpdateStatus     byte
	_PAD16_               [2]byte
	ref30cc92d6           *C.tstrOtaUpdateStatusResp
	allocs30cc92d6        interface{}
}

// OtaUpdateInfo as declared in include/m2m_types.h:2960
type OtaUpdateInfo struct {
	U8NcfUpgradeVersion  uint32
	U8NcfCurrentVersion  uint32
	U8NcdUpgradeVersion  uint32
	U8NcdRequiredUpgrade byte
	U8DownloadUrlOffset  byte
	U8DownloadUrlSize    byte
	__PAD8__             byte
	refde57f176          *C.tstrOtaUpdateInfo
	allocsde57f176       interface{}
}

// OtaHostFileGetStatusResp as declared in include/m2m_types.h:2987
type OtaHostFileGetStatusResp struct {
	U32OtaFileSize     uint32
	U8OtaFileGetStatus byte
	U8CFHandler        byte
	__PAD16__          [2]byte
	ref7a592c76        *C.tstrOtaHostFileGetStatusResp
	allocs7a592c76     interface{}
}

// OtaHostFileReadStatusResp as declared in include/m2m_types.h:3016
type OtaHostFileReadStatusResp struct {
	U16FileBlockSz      uint16
	U8OtaFileReadStatus byte
	__PAD8__            byte
	PFileBuf            [128]byte
	refe198073a         *C.tstrOtaHostFileReadStatusResp
	allocse198073a      interface{}
}

// OtaHostFileEraseStatusResp as declared in include/m2m_types.h:3033
type OtaHostFileEraseStatusResp struct {
	U8OtaFileEraseStatus byte
	__PAD24__            [3]byte
	ref8dd1de0d          *C.tstrOtaHostFileEraseStatusResp
	allocs8dd1de0d       interface{}
}

// OtaStart as declared in include/m2m_types.h:3046
type OtaStart struct {
	U32TotalLen    uint32
	AcUrl          [256]byte
	AcSNI          [64]byte
	U8SSLFlags     byte
	__PAD24__      [3]byte
	refa906f3b7    *C.tstrOtaStart
	allocsa906f3b7 interface{}
}

// SockErr as declared in include/socket.h:686
type SockErr struct {
	U8ErrCode      byte
	ref7741dd47    *C.tstrSockErr
	allocs7741dd47 interface{}
}

// SocketBindMsg as declared in include/socket.h:776
type SocketBindMsg struct {
	Status         int8
	ref5e2a78d3    *C.tstrSocketBindMsg
	allocs5e2a78d3 interface{}
}

// SocketListenMsg as declared in include/socket.h:795
type SocketListenMsg struct {
	Status         int8
	refc9d3f17f    *C.tstrSocketListenMsg
	allocsc9d3f17f interface{}
}

// SocketAcceptMsg as declared in include/socket.h:816
type SocketAcceptMsg struct {
	Sock           int8
	ref91ace819    *C.tstrSocketAcceptMsg
	allocs91ace819 interface{}
}

// SocketConnectMsg as declared in include/socket.h:840
type SocketConnectMsg struct {
	Sock           int8
	S8Error        int8
	ref124e0d02    *C.tstrSocketConnectMsg
	allocs124e0d02 interface{}
}

// SocketRecvMsg as declared in include/socket.h:877
type SocketRecvMsg struct {
	Pu8Buffer        []byte
	S16BufferSize    int16
	U16RemainingSize uint16
	ref4cb67fc3      *C.tstrSocketRecvMsg
	allocs4cb67fc3   interface{}
}
