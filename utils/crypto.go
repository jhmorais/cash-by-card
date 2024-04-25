package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncryptPassword(pwd string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(pwd))
	pwd = hex.EncodeToString(hash.Sum(nil))
	return pwd
}
