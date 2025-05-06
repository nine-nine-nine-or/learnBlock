package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encode md5加密（小写）
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tem := h.Sum(nil)
	return hex.EncodeToString(tem)
}

// Md5EncodeBig md5加密（大写）
func Md5EncodeBig(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// MakePassword 加密
func MakePassword(plainPassword string, salt string) string {
	return Md5Encode(plainPassword + salt)
}

// DoPassword 判断加密数据是否正确
func DoPassword(plainPassword string, salt string, password string) bool {
	return Md5Encode(plainPassword+salt) == password
}
