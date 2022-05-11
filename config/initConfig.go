package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path/filepath"
)

func InitConfig() {
	path, err := filepath.Abs(filepath.Join("config", "config.yml"))
	//注释掉的这行为测试用
	//path, err := filepath.Abs("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)

	}

	err = yaml.Unmarshal(yfile, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
