package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LoadConfig() {
	// 从本地读取环境变量
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalln("cannot get env files..")
	}
}
