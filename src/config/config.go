package config

import (
	"github.com/go-ini/ini"
	"log"
)

type Config struct {
	Mysql struct {
		Host     string `ini:"host"`
		User     string `ini:"user"`
		Password string `ini:"password"`
		DBName   string `ini:"dbname"`
	}
}

func LoadConfig(path string) Config {
	var cf Config
	load, err := ini.Load(path)
	if err != nil {
		log.Fatal("failed to load ini file")
	}
	err = load.Section("mysql").MapTo(&cf.Mysql)
	if err != nil {
		log.Fatal("failed to map ini file to struct")
	}
	return cf
}
