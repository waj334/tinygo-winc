package winc

import "C"
import (
	"encoding/binary"
	"unsafe"

	"github.com/waj334/tinygo-winc/protocol"
)

type EccProvider interface {
	DeriveClientSharedSecret(serverInfo *EcdhReqInfo) (clientInfo *EcdhReqInfo, err error)
	DeriveServerSharedSecret(clientInfo *EcdhReqInfo) (serverInfo *EcdhReqInfo, err error)
	GenerateKey(info *EcdhReqInfo) (result *EcdhReqInfo, err error)
	GenerateSignature(address uint32, info *EcdsaSignReqInfo) (signature []byte, err error)
	VerifySignature(address uint32, info *EcdsaVerifyReqInfo) (result *EcdsaVerifyReqInfo, err error)
}

var (
	sslReplyChan = make(chan any, 1)
)

const (
	sslReqCertificateVerify protocol.OpcodeId = iota
	sslRequestEcc
	sslResponseEcc
	sslIndicateCertificateRevocationList
	sslRequestWriteOwnCertificates
	sslRequestSetCipherSuiteList
	sslResponseSetCipherSuiteList
	sslRespondWriteOwnCertificates
)

const (
	eccReqNone = iota
	eccReqClientEcdh
	eccReqServerEcdh
	eccReqGenKey
	eccReqSignGen
	eccReqSignVerify
)

type EcNamedCurve uint16

const (
	EcSecp192r1 EcNamedCurve = 19
	EcSecp256r1 EcNamedCurve = 23
	EcSecp384r1 EcNamedCurve = 24
	EcSecp521r1 EcNamedCurve = 25
	EcUnknown   EcNamedCurve = 255
)

const (
	SslCipherRsaWithAes128CbcSha           = _NBIT0
	SslCipherRsaWithAes128CbcSha256        = _NBIT1
	SslCipherDheRsaWithAes128CbcSha        = _NBIT2
	SslCipherDheRsaWithAes128CbcSha256     = _NBIT3
	SslCipherRsaWithAes128GcmSha256        = _NBIT4
	SslCipherDheRsaWithAes128GcmSha256     = _NBIT5
	SslCipherRsaWithAes256CbcSha           = _NBIT6
	SslCipherRsaWithAes256CbcSha256        = _NBIT7
	SslCipherDheRsaWithAes256CbcSha        = _NBIT8
	SslCipherDheRsaWithAes256CbcSha256     = _NBIT9
	SslCipherEcdheRsaWithAes128CbcSha      = _NBIT10
	SslCipherEcdheRsaWithAes256CbcSha      = _NBIT11
	SslCipherEcdheRsaWithAes128CbcSha256   = _NBIT12
	SslCipherEcdheEcdsaWithAes128CbcSha256 = _NBIT13
	SslCipherEcdheRsaWithAes128GcmSha256   = _NBIT14
	SslCipherEcdheEcdsaWithAes128GcmSha256 = _NBIT15

	SslEccOnlyCiphers = SslCipherEcdheEcdsaWithAes128CbcSha256 |
		SslCipherEcdheEcdsaWithAes128GcmSha256

	SslEccAllCiphers = SslCipherEcdheRsaWithAes128CbcSha |
		SslCipherEcdheRsaWithAes128CbcSha256 |
		SslCipherEcdheRsaWithAes128GcmSha256 |
		SslCipherEcdheEcdsaWithAes128CbcSha256 |
		SslCipherEcdheEcdsaWithAes128GcmSha256

	SslNonEccCiphersAes128 = SslCipherRsaWithAes128CbcSha |
		SslCipherRsaWithAes128CbcSha256 |
		SslCipherDheRsaWithAes128CbcSha |
		SslCipherDheRsaWithAes128CbcSha256 |
		SslCipherRsaWithAes128GcmSha256 |
		SslCipherDheRsaWithAes128GcmSha256

	SslEccCiphersAes256 = SslCipherEcdheRsaWithAes256CbcSha

	SslNonEccCiphersAes256 = SslCipherRsaWithAes256CbcSha |
		SslCipherRsaWithAes256CbcSha256 |
		SslCipherDheRsaWithAes256CbcSha |
		SslCipherDheRsaWithAes256CbcSha256

	SslCipherAll = SslCipherRsaWithAes128CbcSha |
		SslCipherRsaWithAes128CbcSha256 |
		SslCipherDheRsaWithAes128CbcSha |
		SslCipherDheRsaWithAes128CbcSha256 |
		SslCipherRsaWithAes128GcmSha256 |
		SslCipherDheRsaWithAes128GcmSha256 |
		SslCipherRsaWithAes256CbcSha |
		SslCipherRsaWithAes256CbcSha256 |
		SslCipherDheRsaWithAes256CbcSha |
		SslCipherDheRsaWithAes256CbcSha256 |
		SslCipherEcdheRsaWithAes128CbcSha |
		SslCipherEcdheRsaWithAes128CbcSha256 |
		SslCipherEcdheRsaWithAes128GcmSha256 |
		SslCipherEcdheEcdsaWithAes128CbcSha256 |
		SslCipherEcdheEcdsaWithAes128GcmSha256 |
		SslCipherEcdheRsaWithAes256CbcSha
)

func (w *WINC) handshakeResponse(strEccReqInfo []byte, data []byte) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err = w.hif.Send(GroupSSL, sslResponseEcc|protocol.OpcodeReqDataPkt, strEccReqInfo, data, uint16(len(strEccReqInfo))); err != nil {
		return
	}

	return
}

func (w *WINC) RetrieveHash(address uint32, buf []byte) (err error) {
	if address == 0 {
		return ErrInvalidParameter
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err = w.hif.Receive(address, buf, false); err != nil {
		return err
	}

	return
}

func (w *WINC) RetrieveCert(address uint32) (curve EcNamedCurve, hash, signature []byte, key *ECPoint, err error) {
	key = &ECPoint{}
	var offset uint32
	var keySz uint16
	var hashSz uint16
	var sigSz uint16

	// Receive the curve type
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&keySz)), 2), false); err != nil {
		return
	}
	curve = EcNamedCurve(Htons(keySz))
	offset += 2

	// Receive length of individual EC point (32)
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&key.Size)), 2), false); err != nil {
		return
	}
	key.Size = Htons(key.Size)
	offset += 2

	// Receive length of hash
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&hashSz)), 2), false); err != nil {
		return
	}
	hashSz = Htons(hashSz)
	offset += 2

	// Receive length of signature
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&sigSz)), 2), false); err != nil {
		return
	}
	sigSz = Htons(sigSz)
	offset += 2

	// Receive the EC Points
	if err = w.hif.Receive(address+offset, key.XY[:], false); err != nil {
		return
	}
	offset += uint32(keySz) * 2

	// Receive the hash
	hash = make([]byte, hashSz)
	if err = w.hif.Receive(address+offset, hash, false); err != nil {
		return
	}

	// Receive the signature
	signature = make([]byte, sigSz)
	if err = w.hif.Receive(address+offset, signature, false); err != nil {
		return
	}

	return
}

func (w *WINC) SetActiveCipherSuite(bits uint32) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	strCsList := sslSetActiveCsList{
		CipherSuiteBitmap: bits,
	}

	if err = w.hif.Send(GroupSSL, sslRequestSetCipherSuiteList, strCsList.bytes(), nil, 0); err != nil {
		return
	}

	return
}

//go:noinline
func (w *WINC) sslCallback(id protocol.OpcodeId, sz uint16, address uint32) (data any, err error) {
	if w.EccProvider == nil {
		return
	}

	switch id {
	case sslRequestEcc:
		payload := make([]byte, 112)
		if err = w.hif.Receive(address, payload, false); err != nil {
			return
		}

		// Advance address so that hashes can be received
		address += uint32(len(payload))

		// The first 2 bytes represent the ECC operation to be performed
		eecOp := binary.LittleEndian.Uint16(payload[:])

		var resp []byte
		var signature []byte

		// Call the respective ECC method
		switch eecOp {
		case eccReqClientEcdh:
			serverInfo := EcdhReqInfo{}
			serverInfo.read(payload)

			var clientInfo *EcdhReqInfo
			if clientInfo, err = w.EccProvider.DeriveClientSharedSecret(&serverInfo); err != nil || clientInfo == nil {
				clientInfo = &EcdhReqInfo{}
				clientInfo.Status = 12
			} else {
				clientInfo.Status = 0
			}

			clientInfo.REQ = serverInfo.REQ
			clientInfo.SeqNo = serverInfo.SeqNo
			clientInfo.UserData = serverInfo.UserData

			resp = clientInfo.bytes()

		case eccReqServerEcdh:
			clientInfo := EcdhReqInfo{}
			clientInfo.read(payload)

			var serverInfo *EcdhReqInfo
			if serverInfo, err = w.EccProvider.DeriveServerSharedSecret(&clientInfo); err != nil || serverInfo == nil {
				serverInfo = &EcdhReqInfo{}
				serverInfo.Status = 12
			} else {
				serverInfo.Status = 0
			}

			serverInfo.REQ = clientInfo.REQ
			serverInfo.SeqNo = clientInfo.SeqNo
			serverInfo.UserData = clientInfo.UserData

			resp = clientInfo.bytes()

		case eccReqGenKey:
			info := EcdhReqInfo{}
			info.read(payload)

			var keyInfo *EcdhReqInfo
			if keyInfo, err = w.EccProvider.GenerateKey(&info); err != nil || keyInfo == nil {
				keyInfo = &EcdhReqInfo{}
				keyInfo.Status = 12
			} else {
				keyInfo.Status = 0
			}

			keyInfo.REQ = info.REQ
			keyInfo.SeqNo = info.SeqNo
			keyInfo.UserData = info.UserData

			resp = keyInfo.bytes()

		case eccReqSignGen:
			info := EcdsaSignReqInfo{}
			info.read(payload)

			if signature, err = w.EccProvider.GenerateSignature(address, &info); err != nil || signature == nil {
				info.Status = 12
			} else {
				info.Status = 0
			}

			resp = info.bytes()

		case eccReqSignVerify:
			info := EcdsaVerifyReqInfo{}
			info.read(payload)

			var result *EcdsaVerifyReqInfo
			if result, err = w.EccProvider.VerifySignature(address, &info); err != nil || result == nil {
				result = &EcdsaVerifyReqInfo{}
				result.Status = 12
			} else {
				result.Status = 0
			}

			result.REQ = info.REQ
			result.SeqNo = info.SeqNo
			result.UserData = info.UserData

			resp = result.bytes()
		}

		// Respond to the handshake
		err = w.handshakeResponse(resp, signature)
	case sslResponseSetCipherSuiteList:
		//buf := make([]byte, 4)
		//if err = w.hif.Receive(address, buf, false); err != nil {
		//	return
		//}
		//
		//strCsList := sslSetActiveCsList{}
		//strCsList.read(buf)
		//
		//sslReplyChan <- &strCsList
	case sslRespondWriteOwnCertificates:
	}
	return
}
