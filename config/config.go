package config

type Redis struct {
	RedisAddr     string `yaml:"Addr"`
	RedisPort     string `yaml:"Port"`
	RedisPassword string `yaml:"Password"`
	RedisDB       int    `yaml:"DB"`
}

type MySQL struct {
	Host     string `yaml:"Host"`
	User     string `yaml:"User"`
	Port     string `yaml:"Port"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
}

type Token struct {
	ExpireTime int64 `yaml:"ExpireTime"`
}

type configStruct struct {
	Token `yaml:"Token"`
	Redis `yaml:"Redis"`
	MySQL `yaml:"MySQL"`
}

var Config = configStruct{}

var UserKey = "user:"
