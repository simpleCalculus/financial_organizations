package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type DB struct {
	Driver   string `json:"driver"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var defaultDB = DB{
	Driver:   "postgres",
	Username: "postgres",
	Password: "201640",
	DBName:   "alif",
}

func New() *DB {
	data, err := ioutil.ReadFile("config/config.json")

	if err != nil {
		log.Printf("an error occured while reading from file %v", err)
		return &defaultDB
	}

	var conf DB
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Printf("an error occured while unmarshallig to conf structure %v", err)
		return &defaultDB
	}

	log.Printf("config == > %v", conf)
	return &conf
}
