/*******************************************************************************
  WINC Driver TLS/SSL Header File

  Company:
    Microchip Technology Inc.

  File Name:
    wdrv_winc_ssl.h

  Summary:
    WINC wireless driver TLS/SSL header file.

  Description:
    Provides an interface to configure TLS/SSL support.
 *******************************************************************************/

// DOM-IGNORE-BEGIN
/*******************************************************************************
* Copyright (C) 2019 Microchip Technology Inc. and its subsidiaries.
*
* Subject to your compliance with these terms, you may use Microchip software
* and any derivatives exclusively with Microchip products. It is your
* responsibility to comply with third party license terms applicable to your
* use of third party software (including open source software) that may
* accompany Microchip software.
*
* THIS SOFTWARE IS SUPPLIED BY MICROCHIP "AS IS". NO WARRANTIES, WHETHER
* EXPRESS, IMPLIED OR STATUTORY, APPLY TO THIS SOFTWARE, INCLUDING ANY IMPLIED
* WARRANTIES OF NON-INFRINGEMENT, MERCHANTABILITY, AND FITNESS FOR A
* PARTICULAR PURPOSE.
*
* IN NO EVENT WILL MICROCHIP BE LIABLE FOR ANY INDIRECT, SPECIAL, PUNITIVE,
* INCIDENTAL OR CONSEQUENTIAL LOSS, DAMAGE, COST OR EXPENSE OF ANY KIND
* WHATSOEVER RELATED TO THE SOFTWARE, HOWEVER CAUSED, EVEN IF MICROCHIP HAS
* BEEN ADVISED OF THE POSSIBILITY OR THE DAMAGES ARE FORESEEABLE. TO THE
* FULLEST EXTENT ALLOWED BY LAW, MICROCHIP'S TOTAL LIABILITY ON ALL CLAIMS IN
* ANY WAY RELATED TO THIS SOFTWARE WILL NOT EXCEED THE AMOUNT OF FEES, IF ANY,
* THAT YOU HAVE PAID DIRECTLY TO MICROCHIP FOR THIS SOFTWARE.
*******************************************************************************/
// DOM-IGNORE-END

#ifndef _WDRV_WINC_SSL_H
#define _WDRV_WINC_SSL_H

// *****************************************************************************
// *****************************************************************************
// Section: File includes
// *****************************************************************************
// *****************************************************************************

#include <stdint.h>

//#include "wdrv_winc_common.h"
//#ifndef WDRV_WINC_DEVICE_LITE_DRIVER
//#include "m2m_ssl.h"
//#endif

// *****************************************************************************
// *****************************************************************************
// Section: WINC Driver TLS/SSL Data Types
// *****************************************************************************
// *****************************************************************************

// *****************************************************************************
/* List of IANA cipher suite IDs

  Summary:
    List of IANA cipher suite IDs.

  Description:
    These defines list the IANA cipher suite IDs.

  Remarks:
    None.

*/

#define WDRV_WINC_TLS_NULL_WITH_NULL_NULL                       0x0000
#define WDRV_WINC_TLS_RSA_WITH_AES_128_CBC_SHA                  0x002f
#define WDRV_WINC_TLS_RSA_WITH_AES_128_CBC_SHA256               0x003c
#define WDRV_WINC_TLS_DHE_RSA_WITH_AES_128_CBC_SHA              0x0033
#define WDRV_WINC_TLS_DHE_RSA_WITH_AES_128_CBC_SHA256           0x0067
#define WDRV_WINC_TLS_RSA_WITH_AES_128_GCM_SHA256               0x009c
#define WDRV_WINC_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256           0x009e
#define WDRV_WINC_TLS_RSA_WITH_AES_256_CBC_SHA                  0x0035
#define WDRV_WINC_TLS_RSA_WITH_AES_256_CBC_SHA256               0x003d
#define WDRV_WINC_TLS_DHE_RSA_WITH_AES_256_CBC_SHA              0x0039
#define WDRV_WINC_TLS_DHE_RSA_WITH_AES_256_CBC_SHA256           0x006b
#define WDRV_WINC_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA            0xc013
#define WDRV_WINC_TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA            0xc014
#define WDRV_WINC_TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256         0xc027
#define WDRV_WINC_TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256       0xc023
#define WDRV_WINC_TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256         0xc02f
#define WDRV_WINC_TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256       0xc02b

// *****************************************************************************
/* The size of the the largest supported EC

  Summary:
    The size of the the largest supported EC. For now, assuming
    the 256-bit EC is the largest supported curve type.

  Description:
    These defines the size of the the largest supported EC.

  Remarks:
    None.

*/

#define WDRV_WINC_ECC_LARGEST_CURVE_SIZE     (32)

// *****************************************************************************
/* Maximum size of one coordinate of an EC point

  Summary:
    Maximum size of one coordinate of an EC point.

  Description:
    These defines the maximum size of one coordinate of an EC point.

  Remarks:
    None.

*/

#define WDRV_WINC_ECC_POINT_MAX_SIZE         WDRV_WINC_ECC_LARGEST_CURVE_SIZE

// *****************************************************************************
/* ECC Request Type

  Summary:
    Enumeration of ECC request types.

  Description:
    Types used for ECC requests from the WINC.

  Remarks:
    None.

*/

typedef enum
{
    WDRV_WINC_ECC_REQ_CLIENT_ECDH = ECC_REQ_CLIENT_ECDH,
    WDRV_WINC_ECC_REQ_SERVER_ECDH = ECC_REQ_SERVER_ECDH,
    WDRV_WINC_ECC_REQ_GEN_KEY = ECC_REQ_GEN_KEY,
    WDRV_WINC_ECC_REQ_SIGN_GEN = ECC_REQ_SIGN_GEN,
    WDRV_WINC_ECC_REQ_SIGN_VERIFY = ECC_REQ_SIGN_VERIFY
} WINC_WDRV_ECC_REQ_TYPE;

// *****************************************************************************
/* ECC Status Type

  Summary:
    Enumeration of ECC status types.

  Description:
    Types used for ECC responses to the WINC.

  Remarks:
    None.

*/

typedef enum
{
    WINC_WDRV_ECC_STATUS_SUCCESS,
    WINC_WDRV_ECC_STATUS_FAILURE,
} WINC_WDRV_ECC_STATUS;

// *****************************************************************************
/*  Elliptic Curve Point Representation

  Summary:
    Elliptic Curve point representation structure.

  Description:
    This structure contains information about Elliptic Curve point representation

  Remarks:
    None.

*/

typedef struct
{
    /* The X-coordinate of the ec point. */
    uint8_t     x[WDRV_WINC_ECC_POINT_MAX_SIZE];

    /* The Y-coordinate of the ec point. */
    uint8_t     y[WDRV_WINC_ECC_POINT_MAX_SIZE];

    /* Point size in bytes (for each of the coordinates). */
    uint16_t    size;

    /* ID for the corresponding private key. */
    uint16_t    privKeyID;
} WDRV_WINC_EC_POINT_REP;

// *****************************************************************************
/*  ECDSA Verify Request Information

  Summary:
    ECDSA Verify Request Information structure.

  Description:
    This structure contains information about ECDSA verify request.

  Remarks:
    None.

*/

typedef struct
{
    uint32_t    nSig;
} WDRV_WINC_ECDSA_VERIFY_REQ_INFO;

// *****************************************************************************
/*  ECDSA Sign Request Information

  Summary:
    ECDSA Sign Request Information structure.

  Description:
    This structure contains information about ECDSA sign request.

  Remarks:
    None.

*/

typedef struct
{
    uint16_t    curveType;
    uint16_t    hashSz;
} WDRV_WINC_ECDSA_SIGN_REQ_INFO;

// *****************************************************************************
/*  ECDH Request Information

  Summary:
    ECDH Request Information structure.

  Description:
    This structure contains information about ECDH request from WINC.

  Remarks:
    None.

*/

typedef struct
{
    WDRV_WINC_EC_POINT_REP  pubKey;
    uint8_t                 key[WDRV_WINC_ECC_POINT_MAX_SIZE];
} WDRV_WINC_ECDH_INFO;

// *****************************************************************************
/*  ECC Information Union

  Summary:
    Union combining possible structures for ECC Request Callback

  Description:
    This union contains possible structure returned to the ECC Request Callback.

  Remarks:
    None.
*/

typedef union
{
    WDRV_WINC_ECDH_INFO             ecdhInfo;
    WDRV_WINC_ECDSA_SIGN_REQ_INFO   ecdsaSignReqInfo;
    WDRV_WINC_ECDSA_VERIFY_REQ_INFO ecdsaVerifyReqInfo;
} WDRV_WINC_ECC_REQ_EX_INFO;

// *****************************************************************************
/*  ECC Handshake Information

  Summary:
    ECC handshake Information structure.

  Description:
    This structure contains information about ECC handshakes with the WINC.

  Remarks:
    None.

*/

typedef struct
{
    uint32_t data[2];
} WDRV_WINC_ECC_HANDSHAKE_INFO;

// *****************************************************************************
/*  Cipher Suite Context

  Summary:
    Cipher suite context structure.

  Description:
    This structure contains information about cipher suites.

  Remarks:
    None.

*/

typedef struct
{
    /* Bit mask of cipher suites. */
    uint32_t ciperSuites;
} WDRV_WINC_CIPHER_SUITE_CONTEXT;

#endif /* _WDRV_WINC_SSL_H */
