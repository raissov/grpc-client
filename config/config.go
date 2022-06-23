package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

//Configs- structure for holding other structures that holds the env variables
type Configs struct {
	DB         *DBConf   `yaml:"db"`
	App        *App      `yaml:"server"`
	GRPCConfig *GRPCConf `yaml:"grpc-server"`
}

//DBConf - strucutre for holding env variables assosicated with database
type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	TimeOut  int    `yaml:"timeout"`
}

//App - structure for holding env variables associated with app
type App struct {
	AppPort         string `yaml:"port"`
	AppShutdownTime int    `yaml:"shutdown_time"`
}

type GRPCConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func New() (*Configs, error) {
	yamlBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// parse the YAML stored in the byte slice into the struct
	config := &Configs{}
	err = yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		log.Fatal(err)
	}
	return config, err
}
