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

const (
	// M2M_HIF_MAX_PACKET_SIZE as defined in include/m2m_hif.h:50
	M2M_HIF_MAX_PACKET_SIZE = (1600 - 4)
	// M2M_HIF_HDR_OFFSET as defined in include/m2m_hif.h:65
	M2M_HIF_HDR_OFFSET = 0x5f3860
	// SSL_MAX_OPT_LEN as defined in include/m2m_socket_host_if.h:55
	SSL_MAX_OPT_LEN = 0x5f3860
	// M2M_MAJOR_SHIFT as defined in include/m2m_types.h:54
	M2M_MAJOR_SHIFT = (8)
	// M2M_MINOR_SHIFT as defined in include/m2m_types.h:55
	M2M_MINOR_SHIFT = (4)
	// M2M_PATCH_SHIFT as defined in include/m2m_types.h:56
	M2M_PATCH_SHIFT = (0)
	// M2M_DRV_VERSION_SHIFT as defined in include/m2m_types.h:58
	M2M_DRV_VERSION_SHIFT = (16)
	// M2M_FW_VERSION_SHIFT as defined in include/m2m_types.h:59
	M2M_FW_VERSION_SHIFT = (0)
	// M2M_RELEASE_VERSION_MAJOR_NO as defined in include/m2m_types.h:115
	M2M_RELEASE_VERSION_MAJOR_NO = (19)
	// M2M_RELEASE_VERSION_MINOR_NO as defined in include/m2m_types.h:119
	M2M_RELEASE_VERSION_MINOR_NO = (7)
	// M2M_RELEASE_VERSION_PATCH_NO as defined in include/m2m_types.h:123
	M2M_RELEASE_VERSION_PATCH_NO = (7)
	// M2M_MIN_REQ_DRV_VERSION_MAJOR_NO as defined in include/m2m_types.h:131
	M2M_MIN_REQ_DRV_VERSION_MAJOR_NO = (19)
	// M2M_MIN_REQ_DRV_VERSION_MINOR_NO as defined in include/m2m_types.h:136
	M2M_MIN_REQ_DRV_VERSION_MINOR_NO = (3)
	// M2M_MIN_REQ_DRV_VERSION_PATCH_NO as defined in include/m2m_types.h:140
	M2M_MIN_REQ_DRV_VERSION_PATCH_NO = (0)
	// M2M_MIN_REQ_DRV_SVN_VERSION as defined in include/m2m_types.h:144
	M2M_MIN_REQ_DRV_SVN_VERSION = (0)
	// M2M_BUFFER_MAX_SIZE as defined in include/m2m_types.h:158
	M2M_BUFFER_MAX_SIZE = (uint64(1600) - 4)
	// M2M_MAC_ADDRES_LEN as defined in include/m2m_types.h:162
	M2M_MAC_ADDRES_LEN = 6
	// M2M_ETHERNET_HDR_OFFSET as defined in include/m2m_types.h:166
	M2M_ETHERNET_HDR_OFFSET = 34
	// M2M_ETHERNET_HDR_LEN as defined in include/m2m_types.h:171
	M2M_ETHERNET_HDR_LEN = 14
	// M2M_MAX_SSID_LEN as defined in include/m2m_types.h:176
	M2M_MAX_SSID_LEN = 33
	// M2M_MAX_PSK_LEN as defined in include/m2m_types.h:182
	M2M_MAX_PSK_LEN = 65
	// M2M_MIN_PSK_LEN as defined in include/m2m_types.h:188
	M2M_MIN_PSK_LEN = 9
	// M2M_DEVICE_NAME_MAX as defined in include/m2m_types.h:193
	M2M_DEVICE_NAME_MAX = 48
	// M2M_NTP_MAX_SERVER_NAME_LENGTH as defined in include/m2m_types.h:197
	M2M_NTP_MAX_SERVER_NAME_LENGTH = 32
	// M2M_LISTEN_INTERVAL as defined in include/m2m_types.h:201
	M2M_LISTEN_INTERVAL = 1
	// M2M_CUST_IE_LEN_MAX as defined in include/m2m_types.h:212
	M2M_CUST_IE_LEN_MAX = 252
	// M2M_CRED_STORE_FLAG as defined in include/m2m_types.h:216
	M2M_CRED_STORE_FLAG = 0x01
	// M2M_CRED_ENCRYPT_FLAG as defined in include/m2m_types.h:220
	M2M_CRED_ENCRYPT_FLAG = 0x02
	// M2M_CRED_IS_STORED_FLAG as defined in include/m2m_types.h:224
	M2M_CRED_IS_STORED_FLAG = 0x10
	// M2M_CRED_IS_ENCRYPTED_FLAG as defined in include/m2m_types.h:228
	M2M_CRED_IS_ENCRYPTED_FLAG = 0x20
	// M2M_WIFI_CONN_BSSID_FLAG as defined in include/m2m_types.h:233
	M2M_WIFI_CONN_BSSID_FLAG = 0x01
	// M2M_AUTH_1X_USER_LEN_MAX as defined in include/m2m_types.h:238
	M2M_AUTH_1X_USER_LEN_MAX = 100
	// M2M_AUTH_1X_PASSWORD_LEN_MAX as defined in include/m2m_types.h:242
	M2M_AUTH_1X_PASSWORD_LEN_MAX = 256
	// M2M_AUTH_1X_PRIVATEKEY_LEN_MAX as defined in include/m2m_types.h:245
	M2M_AUTH_1X_PRIVATEKEY_LEN_MAX = 256
	// M2M_AUTH_1X_CERT_LEN_MAX as defined in include/m2m_types.h:249
	M2M_AUTH_1X_CERT_LEN_MAX = 1584
	// M2M_802_1X_UNENCRYPTED_USERNAME_FLAG as defined in include/m2m_types.h:252
	M2M_802_1X_UNENCRYPTED_USERNAME_FLAG = 0x80
	// M2M_802_1X_PREPEND_DOMAIN_FLAG as defined in include/m2m_types.h:256
	M2M_802_1X_PREPEND_DOMAIN_FLAG = 0x40
	// M2M_802_1X_MSCHAP2_FLAG as defined in include/m2m_types.h:261
	M2M_802_1X_MSCHAP2_FLAG = 0x01
	// M2M_802_1X_TLS_FLAG as defined in include/m2m_types.h:264
	M2M_802_1X_TLS_FLAG = 0x02
	// M2M_802_1X_TLS_CLIENT_CERTIFICATE as defined in include/m2m_types.h:267
	M2M_802_1X_TLS_CLIENT_CERTIFICATE = 1
	// M2M_CONFIG_CMD_BASE as defined in include/m2m_types.h:284
	M2M_CONFIG_CMD_BASE = 1
	// M2M_STA_CMD_BASE as defined in include/m2m_types.h:287
	M2M_STA_CMD_BASE = 40
	// M2M_AP_CMD_BASE as defined in include/m2m_types.h:290
	M2M_AP_CMD_BASE = 70
	// M2M_P2P_CMD_BASE as defined in include/m2m_types.h:296
	M2M_P2P_CMD_BASE = 90
	// M2M_SERVER_CMD_BASE as defined in include/m2m_types.h:301
	M2M_SERVER_CMD_BASE = 100
	// M2M_GEN_CMD_BASE as defined in include/m2m_types.h:304
	M2M_GEN_CMD_BASE = 105
	// M2M_OTA_CMD_BASE as defined in include/m2m_types.h:311
	M2M_OTA_CMD_BASE = 100
	// M2M_CRYPTO_CMD_BASE as defined in include/m2m_types.h:319
	M2M_CRYPTO_CMD_BASE = 1
	// M2M_MAX_GRP_NUM_REQ as defined in include/m2m_types.h:324
	M2M_MAX_GRP_NUM_REQ = (127)
	// M2M_SHA256_CONTEXT_BUFF_LEN as defined in include/m2m_types.h:346
	M2M_SHA256_CONTEXT_BUFF_LEN = (128)
	// M2M_SCAN_DEFAULT_NUM_SLOTS as defined in include/m2m_types.h:349
	M2M_SCAN_DEFAULT_NUM_SLOTS = (2)
	// M2M_SCAN_DEFAULT_SLOT_TIME as defined in include/m2m_types.h:352
	M2M_SCAN_DEFAULT_SLOT_TIME = (30)
	// M2M_SCAN_DEFAULT_PASSIVE_SLOT_TIME as defined in include/m2m_types.h:355
	M2M_SCAN_DEFAULT_PASSIVE_SLOT_TIME = (300)
	// M2M_SCAN_DEFAULT_NUM_PROBE as defined in include/m2m_types.h:358
	M2M_SCAN_DEFAULT_NUM_PROBE = (2)
	// M2M_FASTCONNECT_DEFAULT_RSSI_THRESH as defined in include/m2m_types.h:361
	M2M_FASTCONNECT_DEFAULT_RSSI_THRESH = (-45)
	// M2M_MAGIC_APP as defined in include/m2m_types.h:460
	M2M_MAGIC_APP = (uint64(0xef522f61))
	// M2M_SUCCESS as defined in include/nm_common.h:58
	M2M_SUCCESS = 0
	// M2M_ERR_SEND as defined in include/nm_common.h:59
	M2M_ERR_SEND = -1
	// M2M_ERR_RCV as defined in include/nm_common.h:60
	M2M_ERR_RCV = -2
	// M2M_ERR_MEM_ALLOC as defined in include/nm_common.h:61
	M2M_ERR_MEM_ALLOC = -3
	// M2M_ERR_TIME_OUT as defined in include/nm_common.h:62
	M2M_ERR_TIME_OUT = -4
	// M2M_ERR_INIT as defined in include/nm_common.h:63
	M2M_ERR_INIT = -5
	// M2M_ERR_BUS_FAIL as defined in include/nm_common.h:64
	M2M_ERR_BUS_FAIL = -6
	// M2M_NOT_YET as defined in include/nm_common.h:65
	M2M_NOT_YET = -7
	// M2M_ERR_FIRMWARE as defined in include/nm_common.h:66
	M2M_ERR_FIRMWARE = -8
	// M2M_SPI_FAIL as defined in include/nm_common.h:67
	M2M_SPI_FAIL = -9
	// M2M_ERR_FIRMWARE_bURN as defined in include/nm_common.h:68
	M2M_ERR_FIRMWARE_bURN = -10
	// M2M_ACK as defined in include/nm_common.h:69
	M2M_ACK = -11
	// M2M_ERR_FAIL as defined in include/nm_common.h:70
	M2M_ERR_FAIL = -12
	// M2M_ERR_FW_VER_MISMATCH as defined in include/nm_common.h:71
	M2M_ERR_FW_VER_MISMATCH = -13
	// M2M_ERR_SCAN_IN_PROGRESS as defined in include/nm_common.h:72
	M2M_ERR_SCAN_IN_PROGRESS = -14
	// M2M_ERR_INVALID_ARG as defined in include/nm_common.h:75
	M2M_ERR_INVALID_ARG = -15
	// M2M_ERR_INVALID as defined in include/nm_common.h:76
	M2M_ERR_INVALID = -16
	// SOL_SOCKET as defined in include/socket.h:146
	SOL_SOCKET = 1
	// SOL_SSL_SOCKET as defined in include/socket.h:152
	SOL_SSL_SOCKET = 2
	// SO_SET_UDP_SEND_CALLBACK as defined in include/socket.h:158
	SO_SET_UDP_SEND_CALLBACK = 0x00
	// IP_ADD_MEMBERSHIP as defined in include/socket.h:173
	IP_ADD_MEMBERSHIP = 0x01
	// IP_DROP_MEMBERSHIP as defined in include/socket.h:179
	IP_DROP_MEMBERSHIP = 0x02
	// SO_TCP_KEEPALIVE as defined in include/socket.h:185
	SO_TCP_KEEPALIVE = 0x04
	// SO_TCP_KEEPIDLE as defined in include/socket.h:199
	SO_TCP_KEEPIDLE = 0x05
	// SO_TCP_KEEPINTVL as defined in include/socket.h:212
	SO_TCP_KEEPINTVL = 0x06
	// SO_TCP_KEEPCNT as defined in include/socket.h:225
	SO_TCP_KEEPCNT = 0x07
	// SO_SSL_BYPASS_X509_VERIF as defined in include/socket.h:263
	SO_SSL_BYPASS_X509_VERIF = 0x01
	// SO_SSL_SNI as defined in include/socket.h:278
	SO_SSL_SNI = 0x02
	// SO_SSL_ENABLE_SESSION_CACHING as defined in include/socket.h:290
	SO_SSL_ENABLE_SESSION_CACHING = 0x03
	// SO_SSL_ENABLE_SNI_VALIDATION as defined in include/socket.h:305
	SO_SSL_ENABLE_SNI_VALIDATION = 0x04
	// SO_SSL_ALPN as defined in include/socket.h:319
	SO_SSL_ALPN = 0x05
	// SSL_ENABLE_RSA_SHA_SUITES as defined in include/socket.h:336
	SSL_ENABLE_RSA_SHA_SUITES = 0x01
	// SSL_ENABLE_RSA_SHA256_SUITES as defined in include/socket.h:342
	SSL_ENABLE_RSA_SHA256_SUITES = 0x02
	// SSL_ENABLE_DHE_SHA_SUITES as defined in include/socket.h:348
	SSL_ENABLE_DHE_SHA_SUITES = 0x04
	// SSL_ENABLE_DHE_SHA256_SUITES as defined in include/socket.h:354
	SSL_ENABLE_DHE_SHA256_SUITES = 0x08
	// SSL_ENABLE_RSA_GCM_SUITES as defined in include/socket.h:360
	SSL_ENABLE_RSA_GCM_SUITES = 0x10
	// SSL_ENABLE_DHE_GCM_SUITES as defined in include/socket.h:366
	SSL_ENABLE_DHE_GCM_SUITES = 0x20
	// SSL_ENABLE_ALL_SUITES as defined in include/socket.h:372
	SSL_ENABLE_ALL_SUITES = 0x0000003F
	// SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA as defined in include/socket.h:386
	SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA = 0x5f3860
	// SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA256 as defined in include/socket.h:387
	SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA as defined in include/socket.h:388
	SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA = 0x5f3860
	// SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA256 as defined in include/socket.h:389
	SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_RSA_WITH_AES_128_GCM_SHA256 as defined in include/socket.h:390
	SSL_CIPHER_RSA_WITH_AES_128_GCM_SHA256 = 0x5f3860
	// SSL_CIPHER_DHE_RSA_WITH_AES_128_GCM_SHA256 as defined in include/socket.h:391
	SSL_CIPHER_DHE_RSA_WITH_AES_128_GCM_SHA256 = 0x5f3860
	// SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA as defined in include/socket.h:392
	SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA = 0x5f3860
	// SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA256 as defined in include/socket.h:393
	SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA as defined in include/socket.h:394
	SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA = 0x5f3860
	// SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA256 as defined in include/socket.h:395
	SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA as defined in include/socket.h:396
	SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA = 0x5f3860
	// SSL_CIPHER_ECDHE_RSA_WITH_AES_256_CBC_SHA as defined in include/socket.h:397
	SSL_CIPHER_ECDHE_RSA_WITH_AES_256_CBC_SHA = 0x5f3860
	// SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA256 as defined in include/socket.h:398
	SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 as defined in include/socket.h:399
	SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 = 0x5f3860
	// SSL_CIPHER_ECDHE_RSA_WITH_AES_128_GCM_SHA256 as defined in include/socket.h:400
	SSL_CIPHER_ECDHE_RSA_WITH_AES_128_GCM_SHA256 = 0x5f3860
	// SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 as defined in include/socket.h:401
	SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 = 0x5f3860
	// SSL_ECC_ONLY_CIPHERS as defined in include/socket.h:403
	SSL_ECC_ONLY_CIPHERS = SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	// SSL_ECC_ALL_CIPHERS as defined in include/socket.h:410
	SSL_ECC_ALL_CIPHERS = SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_ECDHE_RSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	// SSL_NON_ECC_CIPHERS_AES_128 as defined in include/socket.h:421
	SSL_NON_ECC_CIPHERS_AES_128 = SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_RSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_128_GCM_SHA256
	// SSL_ECC_CIPHERS_AES_256 as defined in include/socket.h:431
	SSL_ECC_CIPHERS_AES_256 = SSL_CIPHER_ECDHE_RSA_WITH_AES_256_CBC_SHA
	// SSL_NON_ECC_CIPHERS_AES_256 as defined in include/socket.h:436
	SSL_NON_ECC_CIPHERS_AES_256 = SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA | SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA | SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA256
	// SSL_CIPHER_ALL as defined in include/socket.h:447
	SSL_CIPHER_ALL = SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_DHE_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_RSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA | SSL_CIPHER_RSA_WITH_AES_256_CBC_SHA256 | SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA | SSL_CIPHER_DHE_RSA_WITH_AES_256_CBC_SHA256 | SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA | SSL_CIPHER_ECDHE_RSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_ECDHE_RSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256 | SSL_CIPHER_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 | SSL_CIPHER_ECDHE_RSA_WITH_AES_256_CBC_SHA
	// IP_PACKET_OFFSET as defined in include/socket.h:475
	IP_PACKET_OFFSET = 0x5f3860
	// SSL_TX_PACKET_OFFSET as defined in include/socket.h:479
	SSL_TX_PACKET_OFFSET = 0x5f3860
	// SSL_FLAGS_ACTIVE as defined in include/socket.h:481
	SSL_FLAGS_ACTIVE = 0x5f3860
	// SSL_FLAGS_BYPASS_X509 as defined in include/socket.h:482
	SSL_FLAGS_BYPASS_X509 = 0x5f3860
	// SSL_FLAGS_2_RESERVD as defined in include/socket.h:483
	SSL_FLAGS_2_RESERVD = 0x5f3860
	// SSL_FLAGS_3_RESERVD as defined in include/socket.h:484
	SSL_FLAGS_3_RESERVD = 0x5f3860
	// SSL_FLAGS_CACHE_SESSION as defined in include/socket.h:485
	SSL_FLAGS_CACHE_SESSION = 0x5f3860
	// SSL_FLAGS_NO_TX_COPY as defined in include/socket.h:486
	SSL_FLAGS_NO_TX_COPY = 0x5f3860
	// SSL_FLAGS_CHECK_SNI as defined in include/socket.h:487
	SSL_FLAGS_CHECK_SNI = 0x5f3860
	// SSL_FLAGS_DELAY as defined in include/socket.h:488
	SSL_FLAGS_DELAY = 0x5f3860
)
