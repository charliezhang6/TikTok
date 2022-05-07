package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(password string, salt int64) string {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(string(salt)))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}
