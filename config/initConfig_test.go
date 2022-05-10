package config

import (
	"gorm.io/gorm/utils/tests"
	"testing"
)

func TestInitConf(t *testing.T) {
	InitConfig()

	tests.AssertEqual(t, 30, Config.ExpireTime)
}
