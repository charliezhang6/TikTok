package config

type Redis struct {
	RedisAddr     string `yaml:"Addr"`
	RedisPort     string `yaml:"Port"`
	RedisPassword string `yaml:"Password"`
	RedisDB       int    `yaml:"DB"`
}

type Token struct {
	ExpireTime int64 `yaml:"ExpireTime"`
}

type configStruct struct {
	Token `yaml:"Token"`
	Redis `yaml:"Redis"`
}

var Config = configStruct{}

var LoginTokenKey = "login_token:"
