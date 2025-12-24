package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

//env-default:"production"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"` //means this fileds value is in yaml file named env
	StoragePath string `yaml:"storage_path" env-required:"true"`  //env-required is just required same as we do in schema
	HTTPServer  `yaml:"http_server" `
}

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		//check if the path is passed in flags
		flags := flag.String("config", "", "path to the configuration file") //name of the flag-first argument, default value-second argument,description-third argument
		flag.Parse()
		configPath = *flags

		// if the path is not passed in flags too than throw error and stop execution
		if configPath == "" {
			log.Fatal("Config path not found")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist : %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file : %s", err.Error()) //err.Error() is same as err.message
	}

	return &cfg
}
