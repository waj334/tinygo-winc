package winc

import "C"
import (
	"encoding/binary"
	"github.com/waj334/tinygo-winc/debug"
	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/types"
	"unsafe"
)

type EccProvider interface {
	DeriveClientSharedSecret(serverInfo *types.EcdhReqInfo) (clientInfo *types.EcdhReqInfo, err error)
	DeriveServerSharedSecret(clientInfo *types.EcdhReqInfo) (serverInfo *types.EcdhReqInfo, err error)
	GenerateKey(info *types.EcdhReqInfo) (result *types.EcdhReqInfo, err error)
	GenerateSignature(address uint32, info *types.EcdsaSignReqInfo) (signature []byte, err error)
	VerifySignature(address uint32, info *types.EcdsaVerifyReqInfo) (result *types.EcdsaVerifyReqInfo, err error)
}

var (
	sslReplyChan = make(chan any, 1)
)

func (w *WINC) handshakeResponse(strEccReqInfo []byte, data []byte) (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err = w.hif.Send(GroupSSL, protocol.OpcodeId(types.M2M_SSL_RESP_ECC)|protocol.OpcodeReqDataPkt, strEccReqInfo, data, uint16(len(strEccReqInfo))); err != nil {
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

func (w *WINC) RetrieveCert(address uint32) (curve types.EnumEcNamedCurve, hash, signature []byte, key *types.ECPoint, err error) {
	key = new(types.ECPoint)
	var offset uint32
	var keySz uint16
	var hashSz uint16
	var sigSz uint16

	// Receive the curve type
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&keySz)), 2), false); err != nil {
		return
	}
	curve = types.EnumEcNamedCurve(Htons(keySz))
	offset += 2

	// Receive length of individual EC point (32)
	if err = w.hif.Receive(address+offset, unsafe.Slice((*uint8)(unsafe.Pointer(&key.U16Size)), 2), false); err != nil {
		return
	}
	key.U16Size = Htons(key.U16Size)
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

	strCsList := types.SslSetActiveCsList{
		U32CsBMP: bits,
	}

	if err = w.hif.Send(GroupSSL, protocol.OpcodeId(types.M2M_SSL_REQ_SET_CS_LIST), strCsList.Bytes(), nil, 0); err != nil {
		return
	}

	return
}

//go:noinline
func (w *WINC) sslCallback(id protocol.OpcodeId, sz uint16, address uint32) (data any, err error) {
	if w.EccProvider == nil {
		return
	}

	switch types.EnumM2mSslCmd(id) {
	case types.M2M_SSL_REQ_ECC:
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
		switch types.EnumEccREQ(eecOp) {
		case types.ECC_REQ_CLIENT_ECDH:
			serverInfo := types.NewEcdhReqInfoRef(C.CBytes(payload))

			serverInfo.Deref()
			serverInfo.Free()
			serverInfo.StrPubKey.Deref()
			serverInfo.StrPubKey.Free()

			var clientInfo *types.EcdhReqInfo
			if clientInfo, err = w.EccProvider.DeriveClientSharedSecret(serverInfo); err != nil || clientInfo == nil {
				clientInfo = &types.EcdhReqInfo{}
				clientInfo.U16Status = 12
			}

			clientInfo.U16REQ = serverInfo.U16REQ
			clientInfo.U16Status = 0
			clientInfo.U32SeqNo = serverInfo.U32SeqNo
			clientInfo.U32UserData = serverInfo.U32UserData

			resp = clientInfo.Bytes()
			debug.DEBUG("length of response: %v", len(resp))
			debug.DEBUG("response: % #x", resp)

		case types.ECC_REQ_SERVER_ECDH:
			clientInfo := types.NewEcdhReqInfoRef(C.CBytes(payload))

			clientInfo.Deref()
			clientInfo.Free()
			clientInfo.StrPubKey.Deref()
			clientInfo.StrPubKey.Free()

			var serverInfo *types.EcdhReqInfo
			if serverInfo, err = w.EccProvider.DeriveServerSharedSecret(clientInfo); err != nil || serverInfo == nil {
				serverInfo.U16Status = 12
			}

			serverInfo.U16REQ = clientInfo.U16REQ
			serverInfo.U16Status = 0
			serverInfo.U32SeqNo = clientInfo.U32SeqNo
			serverInfo.U32UserData = clientInfo.U32UserData

			resp = clientInfo.Bytes()
			debug.DEBUG("length of response: %v", len(resp))
			debug.DEBUG("response: % #x", resp)

		case types.ECC_REQ_GEN_KEY:
			info := types.NewEcdhReqInfoRef(C.CBytes(payload))

			info.Deref()
			info.Free()
			info.StrPubKey.Deref()
			info.StrPubKey.Free()

			var keyInfo *types.EcdhReqInfo
			if info, err = w.EccProvider.GenerateKey(keyInfo); err != nil || info == nil {
				keyInfo = &types.EcdhReqInfo{}
				keyInfo.U16Status = 12
			}

			keyInfo.U16REQ = info.U16REQ
			keyInfo.U16Status = 0
			keyInfo.U32SeqNo = info.U32SeqNo
			keyInfo.U32UserData = info.U32UserData

			resp = keyInfo.Bytes()
			debug.DEBUG("length of response: %v", len(resp))
			debug.DEBUG("response: % #x", resp)

		case types.ECC_REQ_SIGN_GEN:
			info := types.NewEcdsaSignReqInfoRef(C.CBytes(payload))

			info.Deref()
			info.Free()

			if signature, err = w.EccProvider.GenerateSignature(address, info); err != nil || info == nil {
				info.U16Status = 12
			}

			info.U16Status = 0
			resp = info.Bytes()
			debug.DEBUG("length of response: %v", len(resp))
			debug.DEBUG("response: % #x", resp)

		case types.ECC_REQ_SIGN_VERIFY:
			info := types.NewEcdsaVerifyReqInfoRef(C.CBytes(payload))

			info.Deref()
			info.Free()

			var result *types.EcdsaVerifyReqInfo
			if info, err = w.EccProvider.VerifySignature(address, info); err != nil || info == nil {
				result = &types.EcdsaVerifyReqInfo{}
				info.U16Status = 12
			}

			result.U16REQ = info.U16REQ
			result.U16Status = 0
			result.U32SeqNo = info.U32SeqNo
			result.U32UserData = info.U32UserData

			resp = result.Bytes()
			debug.DEBUG("length of response: %v", len(resp))
			debug.DEBUG("response: % #x", resp)
		}

		// Respond to the handshake
		err = w.handshakeResponse(resp, signature)
	case types.M2M_SSL_RESP_SET_CS_LIST:
		//strCsList := &types.SslSetActiveCsList{}
		//if err = w.hif.Receive(address, strCsList.Bytes(), false); err != nil {
		//	return
		//}
		//
		//strCsList.Deref()
		//strCsList.Free()
		//
		//sslReplyChan <- strCsList
	case types.M2M_SSL_RESP_WRITE_OWN_CERTS:
	}
	return
}
