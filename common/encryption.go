package common

import "encoding/base64"

const base64Table = "helloLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw8912"
var coder = base64.NewEncoding(base64Table)

func Base64Encode(src string) string{
	return coder.EncodeToString([]byte(src))
}


func Base64Decode(src string)([]byte, error){
	return coder.DecodeString(src)
}
