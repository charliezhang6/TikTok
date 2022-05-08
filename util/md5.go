package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(password string, salt string) string {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(salt))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}
