/*
 * Copyright (c) 2019. dvnlabs.ml and Buddylang
 * Davin Alfarizky Putra Basudewa <dbasudewa@gmail.com>
 * This project is for Buddylang Api for core functional.Dont expose to public
 * Please kept this source private.
 */

package config

import "encoding/base64"

var TokenSecret = []byte("AjkL1_ow2po_1sa-21")

func TokenSecretEncoded() []byte {
	sEnc := base64.StdEncoding.EncodeToString(TokenSecret)
	return []byte(sEnc)
}
